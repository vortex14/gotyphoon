package proxy

import (
	Errors "errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	u_ "github.com/ahl5esoft/golang-underscore"
	Gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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
	BlockedTime int
	CheckTime   int
	RedisHost   string
}

type ProxyCollection struct {
	singleton.Singleton
	mu              sync.Mutex
	Settings        *Settings
	PrefixNamespace string

	redisService *redis.Service

	stats map[string]int

	list    []string
	locked  []string
	allowed []string
	banned  []string
}

func (c *ProxyCollection) GetFullKeyPath(key string) string {
	return fmt.Sprintf("%s:%s", c.PrefixNamespace, key)
}

func (c *ProxyCollection) GetFullBanPathByKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", BanKey, c.PrefixNamespace, key)
}

func (c *ProxyCollection) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.RemoveBanHistory()
	err = c.RemoveHistory()
	c.init()
	return err
}

func (c *ProxyCollection) blockAvailableProxyByIndex(index int) error {
	logrus.Debug("block proxy by index ", index)
	proxy := c.allowed[index]
	c.allowed = append(c.allowed[:index], c.allowed[index+1:]...)
	c.locked = append(c.locked, proxy)
	return c.redisService.SetExp(c.GetFullKeyPath(proxy), "-", c.Settings.BlockedTime)
}

func (c *ProxyCollection) IsLocked(proxy string) bool {
	return len(c.redisService.Get(c.GetFullKeyPath(proxy))) > 0
}

func (c *ProxyCollection) unblockProxyByValue(proxyAddress string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	index := u_.Chain(c.locked).FindIndex(func(r string, _ int) bool {
		return proxyAddress == r
	})

	if index == -1 {
		return
	}

	proxy := c.locked[index]
	logrus.Debug(fmt.Sprintf("unblock proxy by index: %d , proxy: %s", index, proxy))
	c.locked = append(c.locked[:index], c.locked[index+1:]...)
	c.allowed = append(c.allowed, proxy)
}

func (c *ProxyCollection) RemoveHistory() error {
	for _, value := range c.list {
		err := c.redisService.Remove(c.GetFullKeyPath(value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProxyCollection) RemoveBanHistory() error {
	for _, value := range c.list {
		err := c.redisService.Remove(c.GetFullBanPathByKey(value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProxyCollection) CountByLocked() int {
	return c.redisService.Count(fmt.Sprintf("%s:*", c.PrefixNamespace))
}

func (c *ProxyCollection) CountByBans() int {
	return c.redisService.Count(fmt.Sprintf("%s:*", BanKey))
}

func (c *ProxyCollection) unblockingProxy() {
	for {
		time.Sleep(time.Duration(c.Settings.CheckTime) * time.Second)
		logrus.Debug("checking blocked timeout proxy")
		for _, value := range c.locked {
			proxyData := c.redisService.Get(c.GetFullKeyPath(value))
			if len(proxyData) == 0 {
				logrus.Debug("Unblocked proxy by timeout: ", value)
				c.unblockProxyByValue(value)
			}

		}

		//c.PrintStats()
	}
}

func (c *ProxyCollection) Block(proxy string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	indexBanned := u_.Chain(c.banned).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})

	if indexBanned > -1 {
		return nil
	}

	c.banned = append(c.banned, proxy)

	err := c.redisService.Set(c.GetFullBanPathByKey(proxy), "-")
	if err != nil {
		return err
	}

	indexLocked := u_.Chain(c.locked).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})

	if indexLocked != -1 {
		c.locked = append(c.locked[:indexLocked], c.locked[indexLocked+1:]...)
	}

	indexAllowed := u_.Chain(c.allowed).FindIndex(func(r string, _ int) bool {
		return r == proxy
	})

	if indexAllowed != -1 {
		c.allowed = append(c.allowed[:indexAllowed], c.allowed[indexAllowed+1:]...)
	}

	return nil

}

func (c *ProxyCollection) GetProxy() (error, string) {
	c.Init()
	c.mu.Lock()
	defer c.mu.Unlock()

	saltTime := rand.NewSource(time.Now().UnixNano())
	randomState := rand.New(saltTime)

	count := len(c.allowed)

	if len(c.allowed) == 0 {
		return Errors.New("not found proxy"), ""
	}

	randomIndex := randomState.Intn(count)
	selected := c.allowed[randomIndex]

	c.stats[selected] += 1

	errBlock := c.blockAvailableProxyByIndex(randomIndex)
	if errBlock != nil {
		return errBlock, ""
	}

	return nil, selected
}

func (c *ProxyCollection) PrintStats() {
	println("--------")
	fmt.Printf("Count available: %d , blocked: %d ", len(c.allowed), len(c.locked))
	println("")
	for proxy, stat := range c.stats {
		println(fmt.Sprintf(">>> proxy: %s, stats: %d", proxy, stat))
	}
}

func (c *ProxyCollection) init() {
	proxyEnvList := os.Getenv("PROXY_LIST")
	if len(proxyEnvList) == 0 {
		logrus.Fatal("env PROXY_LIST not found")
	}

	proxyList := strings.Split(proxyEnvList, "\n")
	c.list = proxyList
	c.allowed, c.locked, c.banned = c.getLists(proxyList)
}

func (c *ProxyCollection) Init() *ProxyCollection {
	c.Construct(func() {

		c.stats = make(map[string]int)

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

		c.init()

		l := fmt.Sprintf("Init %d proxies; locked: %d", len(c.allowed), len(c.locked))

		logrus.Debug(l)

		go c.unblockingProxy()

	})

	return c
}

func (c *ProxyCollection) getLists(proxyList []string) ([]string, []string, []string) {
	var allowed []string
	var locked []string
	var banned []string

	for _, value := range proxyList {
		if len(value) == 0 {
			continue
		}

		banKey := fmt.Sprintf("%s:%s", "bans", c.GetFullKeyPath(value))

		switch {
		case len(c.redisService.Get(banKey)) == 0 && len(c.redisService.Get(c.GetFullKeyPath(value))) == 0:
			allowed = append(allowed, value)
		case len(c.redisService.Get(c.GetFullKeyPath(value))) > 0:
			locked = append(locked, value)
		case len(c.redisService.Get(banKey)) > 0:
			banned = append(banned, value)
		}

	}

	return allowed, locked, banned
}

func Constructor(redisHost string) interfaces.ServerInterface {

	proxyCollection := ProxyCollection{
		PrefixNamespace: "domain_",
		Settings: &Settings{
			BlockedTime: 45, //second
			CheckTime:   3,
			RedisHost:   redisHost,
		},
	}

	return (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			Port:  PORT,
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

						err, proxyValue := proxyCollection.GetProxy()
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

						err := proxyCollection.Block(blockProxy[0])
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
						ctx.JSON(200, proxyCollection.stats)
					},
				},
			},
		},
	)
}
