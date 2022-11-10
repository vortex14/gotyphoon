package proxy

import (
	Context "context"
	Errors "errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	Task "github.com/vortex14/gotyphoon/elements/models/task"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	u_ "github.com/ahl5esoft/golang-underscore"
	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

const (
	DESCRIPTION  = "Proxy Service"
	ResourceName = "Proxy"
	NAME         = "Rotator"
	PATH         = "/"
	PORT         = 8987
	BanKey       = "bans"
)

type Settings struct {
	BlockedTime      int
	CheckTime        int
	CheckBlockedTime int
	RedisHost        string
	ConcurrentCheck  int
	Port             int
	PrefixNamespace  string
	CheckHosts       []string
}

type Collection struct {
	singleton.Singleton
	mu  sync.Mutex
	LOG interfaces.LoggerInterface

	Settings *Settings

	redisService *redis.Service

	stats map[string]map[string]int

	list    []string
	locked  map[string][]string
	allowed map[string][]string
	banned  map[string][]string
}

func (c *Collection) GetFullKeyPath(host string, key string) string {
	return fmt.Sprintf("%s:%s:%s", c.Settings.PrefixNamespace, host, key)
}

func (c *Collection) GetFullBanPathByKey(host string, key string) string {
	return fmt.Sprintf("%s:%s:%s:%s", BanKey, c.Settings.PrefixNamespace, host, key)
}

func (c *Collection) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	for domain, _ := range c.stats {
		err = c.RemoveBanHistory(domain)
		c.LOG.Debug("remove history for ", domain)
		err = c.RemoveHistory(domain)
		if err != nil {
			break
		}
	}

	c.init()
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
		err := c.redisService.Remove(c.GetFullBanPathByKey(host, value))
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

func (c *Collection) CountByLocked() int {
	return c.redisService.Count(fmt.Sprintf("%s:*", c.Settings.PrefixNamespace))
}

func (c *Collection) CountByBans() int {
	return c.redisService.Count(fmt.Sprintf("%s:*", BanKey))
}

func (c *Collection) unblockingProxy(host string) {
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

		//c.PrintStats()
	}
}

func (c *Collection) MakeRequestThroughProxy(availableUrl string, proxy string) error {
	taskTest := fake.CreateDefaultTask()

	taskTest.SetFetcherUrl(availableUrl)
	taskTest.SetProxyAddress(proxy)
	taskTest.SetSaveData("SKIP_CN", "skip")

	ctxGroup := Task.NewTaskCtx(taskTest)

	return (&forms.PipelineGroup{
		MetaInfo: &label.MetaInfo{
			Name:     "Http strategy",
			Required: true,
		},
		Stages: []interfaces.BasePipelineInterface{
			net_http.CreateProxyRequestPipeline(&forms.Options{Retry: forms.RetryOptions{MaxCount: 2}}),
			&html.ResponseHtmlPipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "Response pipeline",
					},
				},
				Fn: func(context Context.Context,
					task interfaces.TaskInterface, logger interfaces.LoggerInterface,
					request *http.Request, response *http.Response,
					data *string, doc *goquery.Document) (error, Context.Context) {

					if response.StatusCode > 400 {
						return Errors.New("not ready "), context
					}
					return nil, context
				},
				Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
					logger.Error("pipeline error")
				},
			},
		},
	}).Run(ctxGroup)
}

func (c *Collection) checkingBlocked(host string, availableUrl string) {
	lastIndex := 0
	for {
		c.LOG.Debug(fmt.Sprintf("checking blocked for %s ... every %ds Count: %d", host, c.Settings.CheckBlockedTime, len(c.banned)))
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
				err := c.MakeRequestThroughProxy(availableUrl, proxy)
				if err == nil {
					ansE := c.RemoveProxyBan(host, proxy)
					if ansE == nil {
						c.LOG.Debug(fmt.Sprintf("proxy %s is available for domain %s throught %s", proxy, host, availableUrl))
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

func (c *Collection) Block(proxy string, host string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

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

func (c *Collection) GetProxy(host string) (error, string) {
	c.Init()
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.stats[host]; !ok {
		c.stats[host] = make(map[string]int)
		c.allowed[host], c.locked[host], c.banned[host] = c.getLists(host)
		go c.unblockingProxy(host)
		go c.checkingBlocked(host, "https://2ip.ru")

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

func (c *Collection) init() {
	proxyEnvList := os.Getenv("PROXY_LIST")
	if len(proxyEnvList) == 0 {
		c.LOG.Error("env PROXY_LIST not found")
	}

	proxyList := strings.Split(proxyEnvList, "\n")
	c.list = proxyList

}

func (c *Collection) Init() *Collection {
	c.Construct(func() {
		c.init()
		c.LOG = log.New(log.D{"proxy": "rotator"})
		c.stats = make(map[string]map[string]int)

		redisService := &redis.Service{
			Config: &interfaces.ServiceRedis{
				Name: "Redis proxy data",
				Details: struct {
					Host     string      `yaml:"host"`
					Port     int         `yaml:"port"`
					Password interface{} `yaml:"password"`
				}(struct {
					Host     string
					Port     int
					Password interface{}
				}{Host: c.Settings.RedisHost, Port: 6379}),
			},
		}

		redisService.Init()

		if !redisService.Ping() {
			os.Exit(1)
		}

		c.redisService = redisService

		c.allowed = make(map[string][]string)
		c.locked = make(map[string][]string)
		c.banned = make(map[string][]string)

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
	return fmt.Sprintf("%s:%s", "bans", c.GetFullKeyPath(host, value))
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
	//c.LOG.Debug(">>>>>>>>>>>>>", host, "  ", allowed, locked, banned)
	return allowed, locked, banned
}

func Constructor(opts *Settings) interfaces.ServerInterface {

	proxyCollection := Collection{
		Settings: opts,
	}

	proxyCollection.Init()

	return (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			Port:  opts.Port,
			Level: interfaces.INFO,
			MetaInfo: &label.MetaInfo{
				Name:        NAME,
				Description: DESCRIPTION,
			},
		},
	}).Init().InitLogger().AddResource(
		&forms.Resource{
			MetaInfo: &label.MetaInfo{
				Path:        PATH,
				Name:        ResourceName,
				Description: DESCRIPTION,
			},
			Actions: map[string]interfaces.ActionInterface{
				ping.PATH:  ping.Controller,
				graph.PATH: graph.Controller,
				"proxy": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "proxy",
							Path:        "proxy",
							Description: "Get a new proxy",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get a new proxy.")
						proxy := fake.CreateFakeProxy()
						params := ctx.Request.URL.Query()
						endpointURL := params["url"]
						if len(endpointURL) == 0 {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: Errors.New("not found url for /proxy").Error()})
							return
						}

						u, err := url.Parse(endpointURL[0])
						if err != nil {
							ctx.JSON(500, "not valid url for /proxy")
							return
						}

						err, proxyValue := proxyCollection.GetProxy(u.Hostname())
						if err != nil {
							logger.Debug(err.Error())
							proxy.Success = false
							proxy.Proxy = err.Error()
						}

						proxy.Proxy = proxyValue
						ctx.JSON(200, proxy)
					},
				},
				"blocked": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "proxy",
							Path:        "blocked",
							Description: "blocked proxies",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get stats of blocked proxies")

						ctx.JSON(200, proxyCollection.banned)
					},
				},
				"is_blocked": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "is blocked?",
							Path:        "is_blocked",
							Description: "is blocked the proxy ?",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get stats of blocked proxies")

						params := ctx.Request.URL.Query()
						endpointURL := params["url"]
						proxy := params["proxy"]
						if len(endpointURL) == 0 || len(proxy) == 0 {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: Errors.New("query params hasn't url or proxy").Error()})
							return
						}

						u, err := url.Parse(endpointURL[0])
						if err != nil {
							ctx.JSON(500, "not valid url")
							return
						}

						index := u_.Chain(proxyCollection.banned[u.Hostname()]).FindIndex(func(r string, _ int) bool {
							return proxy[0] == r
						})
						if index > -1 {
							ctx.JSON(200, struct {
								Status  bool
								Message string
								Index   int
							}{
								Status: true,
								Index:  index,
							})
						} else {
							ctx.JSON(404, struct {
								Status  bool
								Message string
								Index   int
							}{
								Status:  false,
								Message: "not found",
								Index:   index,
							})
						}

					},
				},
				"locked": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "locked-proxy",
							Path:        "locked",
							Description: "locked proxies",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get stats of locked proxies")

						ctx.JSON(200, proxyCollection.locked)
					},
				},
				"active": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "proxy",
							Path:        "active",
							Description: "Active list",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Received request for get active proxies")

						params := ctx.Request.URL.Query()
						endpointURL := params["url"]
						if len(endpointURL) == 0 {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: Errors.New("query params hasn't url for show stats by active proxies list of the domain").Error()})
							return
						}

						u, err := url.Parse(endpointURL[0])
						if err != nil {
							ctx.JSON(500, "not valid url for /active")
							return
						}

						ctx.JSON(200, append(proxyCollection.locked[u.Hostname()], proxyCollection.allowed[u.Hostname()]...))
					},
				},
				"block": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "block",
							Path:        "block",
							Description: "block the proxy",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						params := ctx.Request.URL.Query()
						blockURL := params["url"]
						blockProxy := params["proxy"]
						statusCode := params["code"]
						if len(blockURL) == 0 || len(blockProxy) == 0 || len(statusCode) == 0 {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: Errors.New("invalid request for /block").Error()})
							return
						}

						logger.Debug(fmt.Sprintf("block proxy by url %s, proxy: %s, status: %s", blockURL[0], blockProxy[0], statusCode[0]))

						u, err := url.Parse(blockURL[0])
						if err != nil {
							ctx.JSON(500, "not valid url for /block")
							return
						}

						err = proxyCollection.Block(blockProxy[0], u.Hostname())
						if err != nil {
							ctx.JSON(500, struct {
								Success bool
								Message string
							}{Success: false, Message: err.Error()})
							return
						}

						ctx.JSON(200, struct {
							Success bool
						}{Success: true})
					},
				},
				"bad_request": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "bad",
							Path:        "bad_request",
							Description: "error code 500",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Some access errors.")
						ctx.JSON(500, nil)
					},
				},
				"stats": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "stats",
							Path:        "stats",
							Description: "Stats requests",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("Check stats.")

						params := ctx.Request.URL.Query()
						sURL := params["url"]

						if len(sURL) == 0 {
							ctx.JSON(500, "not valid url for /stats")
							return
						}

						u, err := url.Parse(sURL[0])
						if err != nil {
							ctx.JSON(500, "not valid url for /stats")
							return
						}

						ctx.JSON(200, struct {
							Stats   interface{}
							Active  interface{}
							Blocked interface{}
							Locked  interface{}
							List    interface{}
							Allowed interface{}
						}{
							Stats:   proxyCollection.stats,
							Active:  append(proxyCollection.locked[u.Hostname()], proxyCollection.allowed[u.Hostname()]...),
							Blocked: proxyCollection.banned,
							Allowed: proxyCollection.allowed[u.Hostname()],
							Locked:  proxyCollection.locked,
							List:    proxyCollection.list,
						})
					},
				},
				"clear": &gin.Action{
					Action: &forms.Action{
						MetaInfo: &label.MetaInfo{
							Name:        "clear",
							Path:        "clear",
							Description: "clear proxy history",
						},
						Methods: []string{interfaces.GET},
					},
					GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
						logger.Debug("clear proxy history.")
						e := proxyCollection.Clear()
						ctx.JSON(200, e)
					},
				},
			},
		},
	)
}
