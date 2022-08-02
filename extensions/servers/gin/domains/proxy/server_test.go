package proxy

import (
	Context "context"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"github.com/vortex14/gotyphoon/interfaces"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elazarl/goproxy"
	. "github.com/smartystreets/goconvey/convey"
	Task "github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/elements/models/timer"
)

var proxyList = []string{
	"http://E2Wr4v:f3perf@78.76.190.53:9965",
	"http://E2Wr4v:f3perf@78.76.202.52:9034",
	"http://E2Wr4v:f3perf@78.76.200.55:9349",
	"http://E2Wr4v:f3perf@78.76.233.60:9878",
	"http://E2Wr4v:f3perf@78.76.222.60:9878",
}

var BlockedTime = 6
var CheckTime = 3

func init() {
	_ = os.Setenv("PROXY_LIST", strings.Join(proxyList, "\n"))
}

func TestRunProxyServer(t *testing.T) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	Convey("run proxy server", t, func() {

		go func() {
			log.Fatal(http.ListenAndServe(":11313", proxy))
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
				net_http.CreatePrepareRequestPipeline(),
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

func TestBlockProxyPermanently(t *testing.T) {
	Convey("block proxy permanently", t, func() {

		println("TestBlockProxyPermanently!!!")

		coll := (&ProxyCollection{
			PrefixNamespace: "domain",
			Settings: &Settings{
				BlockedTime: BlockedTime, //second
				CheckTime:   CheckTime,
			},
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

func TestRecoveryState(t *testing.T) {
	coll := (&ProxyCollection{
		PrefixNamespace: "domain",
		Settings: &Settings{
			BlockedTime: BlockedTime, //second
			CheckTime:   CheckTime,
		},
	}).Init()

	_ = coll.Clear()

	Convey("Block all proxy", t, func() {
		for _, proxy := range proxyList {
			err := coll.Block(proxy)
			So(err, ShouldBeNil)
		}

		count := coll.CountByBans()

		So(count, ShouldEqual, len(proxyList))
	})

	Convey("check a new collection state", t, func() {
		coll2 := (&ProxyCollection{
			PrefixNamespace: "domain",
			Settings: &Settings{
				BlockedTime: BlockedTime, //second
				CheckTime:   CheckTime,
			},
		}).Init()

		So(coll2.CountByBans(), ShouldEqual, len(proxyList))
	})

	Convey("remove ban history", t, func() {
		err := coll.RemoveBanHistory()

		So(err, ShouldBeNil)

		count := coll.CountByBans()

		So(count, ShouldEqual, 0)
	})
}

func TestBlockProxyByExpiration(t *testing.T) {
	Convey("blocking proxy by time", t, func(c C) {
		coll := (&ProxyCollection{
			PrefixNamespace: "domain",
			Settings: &Settings{
				BlockedTime: BlockedTime, //second
				CheckTime:   CheckTime,
			},
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
		coll := (&ProxyCollection{
			PrefixNamespace: "domain",
			Settings: &Settings{
				BlockedTime: BlockedTime, //second
				CheckTime:   CheckTime,
			},
		}).Init()

		err := coll.RemoveHistory()

		So(err, ShouldBeNil)

		So(coll.CountByLocked(), ShouldEqual, 0)

	})

}

func TestGetConcurrentProxyCollection(t *testing.T) {
	Convey("create a new proxy collection", t, func() {
		coll := (&ProxyCollection{
			PrefixNamespace: "domain",
			Settings: &Settings{
				BlockedTime: BlockedTime, //second
				CheckTime:   CheckTime,
			},
		}).Init()

		So(len(coll.allowed), ShouldEqual, 5)
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

			So(existCount, ShouldEqual, len(proxyList)-1)

			coll.PrintStats()
			err := coll.RemoveHistory()

			So(err, ShouldBeNil)

			So(coll.CountByLocked(), ShouldEqual, 0)

		})
	})
}