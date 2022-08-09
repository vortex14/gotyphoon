package proxy

import (
	Context "context"
	"fmt"
	underscore "github.com/ahl5esoft/golang-underscore"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/elazarl/goproxy"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/vortex14/gotyphoon/log"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Task "github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/elements/models/timer"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"github.com/vortex14/gotyphoon/interfaces"
)

var validProxyAddress = "http://localhost:11316"

var proxyList = []string{
	"http://E2Wr4v:f3perf@78.76.190.53:9965",
	"http://E2Wr4v:f3perf@78.76.202.52:9034",
	"http://E2Wr4v:f3perf@78.76.200.55:9349",
	"http://E2Wr4v:f3perf@78.76.233.60:9878",
	"http://E2Wr4v:f3perf@78.76.222.60:9878",
	validProxyAddress,
}

var BlockedTime = 6
var CheckTime = 3

var settings = &Settings{
	PrefixNamespace: "domain",
	CheckHosts:      []string{"https://2ip.ru"},
	BlockedTime:     BlockedTime,
	CheckTime:       CheckTime,
	RedisHost:       "localhost",
	ConcurrentCheck: 3,
	Port:            11222,
}

func init() {
	log.InitD()
	_ = os.Setenv("PROXY_LIST", strings.Join(proxyList, "\n"))
}

func TestReadPart(t *testing.T) {

	lastIndex := 0
	maxCycleIter := 10
	cycleIter := 1

	for {
		wg := sync.WaitGroup{}
		println("iter ", lastIndex, len(proxyList), lastIndex+settings.ConcurrentCheck, lastIndex == len(proxyList), cycleIter)

		if maxCycleIter == cycleIter {
			break
		}

		if lastIndex == len(proxyList) {
			lastIndex = 0
			println("reset")
		}

		step := lastIndex + settings.ConcurrentCheck
		residue := len(proxyList) - lastIndex

		if step > residue {
			step = residue
		}

		for i, proxy := range proxyList[lastIndex : lastIndex+step] {
			println("check", proxy, i)

			if (i + 1) >= settings.ConcurrentCheck {
				println("busy ...")
				break
			}

			wg.Add(1)

			go func() {
				time.Sleep(time.Duration(settings.CheckTime) * time.Second)
				wg.Done()
			}()
			lastIndex += 1
		}

		time.Sleep(1 * time.Second)
		cycleIter += 1
		wg.Wait()

	}

	println("Done", cycleIter)
}

func TestRunProxyServer(t *testing.T) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	Convey("run proxy server", t, func() {

		go func() {
			_ = http.ListenAndServe(":11313", proxy)
		}()
		time.Sleep(2 * time.Second)

	})

	Convey("create request from local proxy server", t, func() {
		var proxyAddress = "http://localhost:11313"

		taskTest := fake.CreateDefaultTask()

		taskTest.SetFetcherUrl("https://2ip.ru/")
		taskTest.SetProxyAddress(proxyAddress)
		taskTest.SetProxyServerUrl("")

		ctxGroup := Task.NewTaskCtx(taskTest)

		err := (&forms.PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name:     "Http strategy",
				Required: true,
			},
			Stages: []interfaces.BasePipelineInterface{
				net_http.CreateRequestPipeline(),
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

						resp := doc.Find(".ip").Text()

						logger.Info(resp)

						Convey("check response 2ip.ru", func(c C) {
							c.So(len(resp) > 0, ShouldBeTrue)
						})

						return nil, context
					},
					Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
						logger.Error("pipeline error")
					},
				},
			},
		}).Run(ctxGroup)

		So(err, ShouldBeNil)

	})
}

func MakeRequestThroughProxy() error {
	taskTest := fake.CreateDefaultTask()

	taskTest.SetFetcherUrl("https://2ip.ru/")
	taskTest.SetProxyAddress(validProxyAddress)
	taskTest.SetProxyServerUrl(fmt.Sprintf("http://localhost:%d", settings.Port))

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

					resp := doc.Find(".ip").Text()

					logger.Info(resp)

					Convey("check response 2ip.ru", func(c C) {
						c.So(len(resp) > 0, ShouldBeTrue)
					})

					return nil, context
				},
				Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
					logger.Error("pipeline error")
				},
			},
		},
	}).Run(ctxGroup)
}

func TestAvailableProxy(t *testing.T) {

	logger := log.New(log.D{"test": "test"})

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	coll := (&Collection{
		Settings: settings,
	}).Init()
	_ = coll.RemoveBanHistory()

	proxyService := Constructor(settings)

	go func() {
		_ = proxyService.Run()
	}()

	time.Sleep(4 * time.Second)

	Convey("create a new request to unavailable proxy", t, func() {

		So(MakeRequestThroughProxy(), ShouldBeError)

	})

	Convey("checking blocked proxy", t, func() {

		coll2 := (&Collection{
			Settings: settings,
		}).Init()

		println(coll2.banned)

		index := underscore.Chain(coll2.banned).FindIndex(func(r string, _ int) bool {
			return validProxyAddress == r
		})

		So(index > -1, ShouldBeTrue)
	})

	Convey("Run proxy server", t, func() {

		coll2 := (&Collection{
			Settings: settings,
		}).Init()

		logger.Debug(fmt.Sprintf("%+v", coll2.banned))

		go func() {
			_ = http.ListenAndServe(":11316", proxy)
		}()

	})

	Convey("check unbanning proxy", t, func() {

		So(coll.IsBannedProxy(validProxyAddress), ShouldBeTrue)

		count := 0
		for {
			if count > 5 {
				So(coll.IsBannedProxy(validProxyAddress), ShouldBeFalse)
				So(coll.IsAllowed(validProxyAddress), ShouldBeTrue)
				break
			}

			time.Sleep(2 * time.Second)

			count += 1
		}

	})
}

func TestBlockProxyPermanently(t *testing.T) {
	Convey("block proxy permanently", t, func() {

		println("TestBlockProxyPermanently!!!")

		coll := (&Collection{
			Settings: settings,
		}).Init()

		err := coll.Clear()
		So(err, ShouldBeNil)

		e := coll.RemoveBanHistory()
		So(e, ShouldBeNil)
		count := coll.CountByBans()
		So(count, ShouldEqual, 0)

		err = coll.Block(proxyList[0])
		So(err, ShouldBeNil)

		count = coll.CountByBans()
		So(count, ShouldEqual, 1)
		e = coll.RemoveBanHistory()
		So(e, ShouldBeNil)

		count = coll.CountByBans()
		So(count, ShouldEqual, 0)

	})
}

func TestBlockAllProxies(t *testing.T) {

	coll := (&Collection{
		Settings: settings,
	}).Init()

	_ = coll.RemoveBanHistory()

	Convey("Block all proxy", t, func() {
		for _, proxy := range proxyList {
			err := coll.Block(proxy)
			So(err, ShouldBeNil)
		}

		count := coll.CountByBans()

		So(count, ShouldEqual, len(proxyList))

	})

	Convey("check a new collection state", t, func() {
		coll2 := (&Collection{
			Settings: settings,
		}).Init()

		So(coll2.CountByBans(), ShouldEqual, len(proxyList))

		err := coll.RemoveBanHistory()
		So(err, ShouldBeNil)

		count := coll.CountByBans()

		So(count, ShouldEqual, 0)
	})
}

func TestBlockProxyByExpiration(t *testing.T) {
	Convey("blocking proxy by time", t, func(c C) {
		coll := (&Collection{
			Settings: settings,
		}).Init()

		e := coll.Clear()
		So(e, ShouldBeNil)

		println(coll.CountByLocked(), coll.CountByBans())

		err, proxy := coll.GetProxy()
		So(err, ShouldBeNil)

		st0 := timer.SetTimeout(func(args ...interface{}) {
			println(proxy, " ", coll.IsLocked(proxy))
			c.So(coll.IsLocked(proxy), ShouldBeTrue)
		}, BlockedTime+1)

		st0.Await()

		st := timer.SetTimeout(func(args ...interface{}) {
			println(proxy, " ", coll.IsLocked(proxy))
			c.So(coll.IsLocked(proxy), ShouldBeFalse)
		}, (BlockedTime+1)*1000)

		st.Await()
	})
}

func TestRemoveRedisCollection(t *testing.T) {

	Convey("remove locked proxies", t, func() {
		coll := (&Collection{
			Settings: settings,
		}).Init()

		err := coll.RemoveHistory()

		So(err, ShouldBeNil)

		So(coll.CountByLocked(), ShouldEqual, 0)

	})

}

func TestGetConcurrentProxyCollection(t *testing.T) {
	Convey("create a new proxy collection", t, func() {
		coll := (&Collection{
			Settings: settings,
		}).Init()

		_ = coll.RemoveBanHistory()
		_ = coll.RemoveHistory()

		So(len(coll.allowed), ShouldEqual, len(proxyList))
		g := sync.WaitGroup{}
		g.Add(4)
		resultProxyEnv := make(map[string]string)

		Convey("locking all proxy", func(c C) {

			go func() {
				err, proxy := coll.GetProxy()
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			go func() {
				err, proxy := coll.GetProxy()
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			go func() {
				err, proxy := coll.GetProxy()
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			go func() {
				err, proxy := coll.GetProxy()
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			g.Wait()
			existCount := 0
			for _, v := range proxyList {
				_, ok := resultProxyEnv[v]
				if ok {
					existCount += 1
				}
			}

			So(existCount, ShouldEqual, len(proxyList)-2)

			coll.PrintStats()
			err := coll.RemoveHistory()

			So(err, ShouldBeNil)

			So(coll.CountByLocked(), ShouldEqual, 0)

		})
	})
}
