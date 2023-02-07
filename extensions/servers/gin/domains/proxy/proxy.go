package proxy

import (
	Errors "errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	u_ "github.com/ahl5esoft/golang-underscore"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

type Collection struct {
	singleton.Singleton
	mu  sync.RWMutex
	LOG interfaces.LoggerInterface

	Settings *Settings

	redisService *redis.Service

	stats map[string]map[string]int

	list    []string
	locked  map[string][]string
	allowed map[string][]string
	banned  map[string][]string

	observableHosts    map[string]bool
	availableEndpoints map[string]*UpdateCheckPayload // domain key : available endpoint
}

func (c *Collection) GetFullKeyPath(host string, key string) string {
	return fmt.Sprintf("%s:%s:%s", c.Settings.PrefixNamespace, host, key)
}

func (c *Collection) GetFullBanPathByKey(host string, key string) string {
	return fmt.Sprintf("%s:%s:%s:%s", c.Settings.PrefixNamespace, BanKey, host, key)
}

func (c *Collection) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	for domain := range c.stats {
		err = c.RemoveBanHistory(domain)
		c.LOG.Debug("remove history for ", domain)
		err = c.RemoveHistory(domain)
		if err != nil {
			break
		}
	}

	return err
}

func (c *Collection) blockAvailableProxyByIndex(host string, index int) error {
	c.LOG.Debug("block proxy by index ", index)
	proxy := c.allowed[host][index]
	c.allowed[host] = append(c.allowed[host][:index], c.allowed[host][index+1:]...)
	c.locked[host] = append(c.locked[host], proxy)
	return c.redisService.SetExp(c.GetFullKeyPath(host, proxy), "-", c.Settings.BlockedTime)
}

func (c *Collection) IsLocked(host string, proxy string) bool {
	return len(c.redisService.Get(c.GetFullKeyPath(host, proxy))) > 0
}

func (c *Collection) unblockProxyByValue(host string, proxyAddress string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	index := u_.Chain(c.locked[host]).FindIndex(func(r string, _ int) bool {
		return proxyAddress == r
	})

	if index == -1 {
		return
	}

	proxy := c.locked[host][index]
	c.LOG.Debug(fmt.Sprintf("unblock proxy for %s by index: %d , proxy: %s", host, index, proxy))
	c.locked[host] = append(c.locked[host][:index], c.locked[host][index+1:]...)
	c.allowed[host] = append(c.allowed[host], proxy)
}

func (c *Collection) RemoveHistory(host string) error {

	for _, value := range c.list {
		err := c.redisService.Remove(c.GetFullKeyPath(host, value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Collection) RemoveBanHistory(host string) error {
	for _, value := range c.list {
		err := c.RemoveProxyBan(host, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Collection) RemoveProxyBan(host string, proxy string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.redisService.Remove(c.GetFullBanPathByKey(host, proxy))
	indexBanned := u_.Chain(c.banned[host]).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})
	if err == nil && indexBanned > -1 {
		c.banned[host] = append(c.banned[host][:indexBanned], c.banned[host][indexBanned+1:]...)
		c.allowed[host] = append(c.allowed[host], proxy)
	}

	return err
}

func (c *Collection) CountByLocked(host string) int {
	return c.redisService.Count(fmt.Sprintf("%s:%s:*", c.Settings.PrefixNamespace, host))
}

func (c *Collection) CountByBans(host string) int {
	return c.redisService.Count(fmt.Sprintf("%s:%s:%s:*", c.Settings.PrefixNamespace, BanKey, host))
}

func (c *Collection) unblockingProxy(availableURL *url.URL) {
	host := availableURL.Hostname()
	for {
		time.Sleep(time.Duration(c.Settings.CheckTime) * time.Second)
		c.LOG.Debug("checking locked timeout proxy")
		for _, value := range c.locked[host] {
			proxyData := c.redisService.Get(c.GetFullKeyPath(host, value))
			if len(proxyData) == 0 {
				c.LOG.Debug("Unblocked proxy by timeout: ", value)
				c.unblockProxyByValue(host, value)
			}

		}

	}
}

func (c *Collection) getAvailablePayload(availableURL *url.URL) *UpdateCheckPayload {

	host := availableURL.Host
	if _, ok := c.availableEndpoints[host]; !ok {
		c.availableEndpoints[host] = &UpdateCheckPayload{Url: availableURL}
		_, _ = c.availableEndpoints[host].ParseUrl()
	}
	return c.availableEndpoints[host]
}

func (c *Collection) setAvailablePayload(payload *UpdateCheckPayload) {
	host := payload.Url.Hostname()
	_, _ = payload.ParseUrl()
	c.availableEndpoints[host] = payload
}

func (c *Collection) checkingBlocked(availableURL *url.URL) {
	lastIndex := 0

	host := availableURL.Hostname()

	for {
		c.LOG.Debug(fmt.Sprintf("checking blocked for host: %s ... every %ds Count: %d", host, c.Settings.CheckBlockedTime, len(c.banned[host])))
		wg := &sync.WaitGroup{}

		if lastIndex == len(c.banned[host]) || lastIndex > len(c.banned[host]) {
			lastIndex = 0
		}

		step := lastIndex + c.Settings.ConcurrentCheck
		residue := len(c.banned[host]) - lastIndex

		if step > len(c.banned[host]) {
			step = lastIndex + residue
		}

		if len(c.banned[host]) == 0 {
			c.LOG.Debug("not found proxies for checking")
			time.Sleep(time.Duration(c.Settings.CheckBlockedTime) * time.Second)
			continue
		}
		c.LOG.Debug(fmt.Sprintf("lastIndex: %d ; Step: %d; residue: %d", lastIndex, step, residue))
		for i, proxy := range c.banned[host][lastIndex:step] {
			c.LOG.Debug(fmt.Sprintf("check %s", proxy))
			if (i + 1) >= c.Settings.ConcurrentCheck {
				break
			}

			wg.Add(1)

			go func(wg *sync.WaitGroup, proxy string) {
				task := fake.CreateDefaultTask()

				payload := c.getAvailablePayload(availableURL)

				task.SetFetcherUrl(payload.UrlSource)
				task.SetProxyAddress(proxy)
				task.SetHeaders(payload.Headers)

				c.LOG.Debug("create a new request to ", payload.UrlSource, " through "+proxy+
					fmt.Sprintf(" Headers keys: %d", len(payload.Headers)))

				err := netHttp.MakeRequestThroughProxy(task, func(logger interfaces.LoggerInterface,
					response *http.Response, doc *goquery.Document) bool {

					return !(response.StatusCode > 400)
				})
				if err == nil {
					ansE := c.RemoveProxyBan(host, proxy)
					if ansE == nil {
						c.LOG.Debug(fmt.Sprintf("proxy %s   %s throught %s is available", proxy, host, payload.UrlSource))
					} else {
						c.LOG.Error(ansE)
					}

				} else {
					c.LOG.Error(err)
				}

				wg.Done()

			}(wg, proxy)
			lastIndex += 1
		}
		c.LOG.Debug("waiting for checked")
		wg.Wait()

		time.Sleep(time.Duration(c.Settings.CheckBlockedTime) * time.Second)
		c.LOG.Debug("next list step to check")

	}
}

func (c *Collection) Block(proxy string, u *url.URL) error {
	host := u.Hostname()
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.observableHosts[host]; !ok {
		c.observableHosts[host] = true
		go c.unblockingProxy(u)
		payload := c.getAvailablePayload(u)
		c.setAvailablePayload(payload)
		go c.checkingBlocked(u)
	}

	indexBanned := u_.Chain(c.banned[host]).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})

	if indexBanned > -1 {
		return nil
	}

	c.banned[host] = append(c.banned[host], proxy)

	err := c.redisService.Set(c.GetFullBanPathByKey(host, proxy), "-")
	if err != nil {
		return err
	}

	indexLocked := u_.Chain(c.locked[host]).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})

	if indexLocked != -1 {
		c.locked[host] = append(c.locked[host][:indexLocked], c.locked[host][indexLocked+1:]...)
	}

	indexAllowed := u_.Chain(c.allowed[host]).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})

	if indexAllowed != -1 {
		c.allowed[host] = append(c.allowed[host][:indexAllowed], c.allowed[host][indexAllowed+1:]...)
	}

	return nil

}

func (c *Collection) GetProxy(u *url.URL) (error, string) {
	c.Init()
	c.mu.Lock()
	defer c.mu.Unlock()
	host := u.Hostname()

	if _, ok := c.stats[host]; !ok {
		c.stats[host] = make(map[string]int)
		c.allowed[host], c.locked[host], c.banned[host] = c.getLists(host)
		c.observableHosts[host] = true
		go c.unblockingProxy(u)
		go c.checkingBlocked(u)
	}

	saltTime := rand.NewSource(time.Now().UnixNano())
	randomState := rand.New(saltTime)

	count := len(c.allowed[host])

	if len(c.allowed[host]) == 0 {
		return Errors.New("not found proxy"), ""
	}

	randomIndex := randomState.Intn(count)
	selected := c.allowed[host][randomIndex]

	c.stats[host][selected] += 1

	errBlock := c.blockAvailableProxyByIndex(host, randomIndex)
	if errBlock != nil {
		return errBlock, ""
	}

	return nil, selected
}

func (c *Collection) PrintStats() {
	println("--------")
	fmt.Printf("Count available: %d , blocked: %d ", len(c.allowed), len(c.locked))
	println("")
	for domain, stat := range c.stats {
		println(fmt.Sprintf(">>> domain: %s", domain))

		for proxy, value := range stat {
			println(fmt.Sprintf(">>> proxy: %s; stat: %d", proxy, value))
		}
	}
}

func (c *Collection) Init() *Collection {
	c.Construct(func() {
		c.LOG = log.New(log.D{"proxy": "rotator"})

		c.locked = make(map[string][]string)
		c.banned = make(map[string][]string)
		c.allowed = make(map[string][]string)
		c.stats = make(map[string]map[string]int)
		c.observableHosts = make(map[string]bool)
		c.availableEndpoints = make(map[string]*UpdateCheckPayload)

		proxyEnvList := os.Getenv("PROXY_LIST")
		if len(proxyEnvList) == 0 {
			c.LOG.Error("env PROXY_LIST not found")
			os.Exit(1)
		}
		availableEnvList := os.Getenv("AVAILABLE_LIST")

		if len(availableEnvList) == 0 {
			c.LOG.Warning("env AVAILABLE_LIST not found. will be set first url for get proxy by domain")
		}

		proxyList := strings.Split(proxyEnvList, "\n")
		c.list = proxyList

		redisService := &redis.Service{
			Config: &interfaces.ServiceRedis{
				Name: "Redis proxy data",
				Details: interfaces.RedisDetails{
					Host:     c.Settings.RedisDetails.Host,
					Port:     c.Settings.RedisDetails.Port,
					Password: c.Settings.Password,
				},
			},
		}

		redisService.Init()

		if !redisService.Ping() {
			os.Exit(1)
		}

		c.redisService = redisService

	})

	return c
}

func (c *Collection) IsBannedProxy(host string, proxy string) bool {
	return len(c.redisService.Get(c.getBanKey(host, proxy))) > 0
}

func (c *Collection) IsAllowed(proxy string) bool {
	return u_.Chain(c.allowed).FindIndex(func(r string, _ int) bool {
		return r == proxy
	}) > -1
}

func (c *Collection) getBanKey(host string, value string) string {
	return fmt.Sprintf("%s:%s:%s", c.Settings.PrefixNamespace, BanKey, c.GetFullKeyPath(host, value))
}

func (c *Collection) getLists(host string) ([]string, []string, []string) {
	allowed := make([]string, 0)
	locked := make([]string, 0)
	banned := make([]string, 0)

	for _, value := range c.list {
		if len(value) == 0 {
			continue
		}

		banKey := c.getBanKey(host, value)

		switch {
		case len(c.redisService.Get(banKey)) == 0 && len(c.redisService.Get(c.GetFullKeyPath(host, value))) == 0:
			allowed = append(allowed, value)
		case len(c.redisService.Get(c.GetFullKeyPath(host, value))) > 0:
			locked = append(locked, value)
		case len(c.redisService.Get(banKey)) > 0:
			banned = append(banned, value)
		}

	}

	return allowed, locked, banned
}
