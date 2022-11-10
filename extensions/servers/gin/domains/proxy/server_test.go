package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/elazarl/goproxy"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/vortex14/gotyphoon/elements/models/timer"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	net_html "github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
)

var (
	validProxyAddress = "http://localhost:11316"

	proxyServiceURL = "http://localhost:11222"

	proxyList = []string{
		"http://E2Wr4v:f3perf@78.76.190.53:9965",
		"http://E2Wr4v:f3perf@78.76.202.52:9034",
		"http://E2Wr4v:f3perf@78.76.200.55:9349",
		"http://E2Wr4v:f3perf@78.76.233.60:9878",
		"http://E2Wr4v:f3perf@78.76.222.60:9878",
		validProxyAddress,
	}

	BlockedTime      = 6
	CheckTime        = 3
	CheckBlockedTime = 3

	validURL = "https://2ip.ru/"

	endpointURL, _ = url.Parse(validURL)

	LOG interfaces.LoggerInterface

	settings = &Settings{
		PrefixNamespace:  "domain",
		CheckHosts:       []string{"https://2ip.ru"},
		BlockedTime:      BlockedTime,
		CheckTime:        CheckTime,
		CheckBlockedTime: CheckBlockedTime,
		RedisHost:        "localhost",
		ConcurrentCheck:  3,
		Port:             11222,
	}
)

func init() {
	log.InitD()
	LOG = log.New(map[string]interface{}{"service": "proxy"})

	_ = os.Setenv("PROXY_LIST", strings.Join(proxyList, "\n"))

	availableMap := map[string]string{
		"2ip.ru": "https://2ip.ru",
	}

	_, s := utils.DumpPrettyJson(availableMap)

	_ = os.Setenv("AVAILABLE_LIST", s)

}

func TestUrl(t *testing.T) {
	Convey("test url", t, func() {
		So(endpointURL.Hostname(), ShouldEqual, "2ip.ru")
		So(endpointURL.String(), ShouldEqual, "https://2ip.ru/")
	})
}

func TestIp(t *testing.T) {
	Convey("get raw ip address without auth", t, func() {
		v, e := url.Parse(proxyList[0])
		So(e, ShouldBeNil)
		So(v.Hostname(), ShouldEqual, "78.76.190.53")
		So(v.Host, ShouldEqual, "78.76.190.53:9965")
	})
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

		err := net_html.MakeRequestThroughProxy(taskTest, func(logger interfaces.LoggerInterface,
			response *http.Response, doc *goquery.Document) bool {

			resp := doc.Find(".ip").Text()
			return len(resp) > 0

		})

		So(err, ShouldBeNil)

	})
}

func TestAvailableProxy(t *testing.T) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	go func() {

		proxyService := Constructor(settings)
		_ = proxyService.Run()
	}()

	time.Sleep(4 * time.Second)

	Convey("create a new request to unavailable proxy", t, func() {

		taskTest := fake.CreateDefaultTask()

		taskTest.SetFetcherUrl(validURL)
		taskTest.SetProxyAddress(validProxyAddress)
		taskTest.SetProxyServerUrl(fmt.Sprintf("http://localhost:%d", settings.Port))

		So(net_html.MakeRequestThroughProxy(taskTest, func(logger interfaces.LoggerInterface,
			response *http.Response, doc *goquery.Document) bool {

			resp := doc.Find(".ip").Text()

			logger.Info(resp)

			return len(resp) > 0

		}), ShouldBeError)

	})

	time.Sleep(4 * time.Second)

	Convey("checking blocked proxy", t, func() {
		checkURL := fmt.Sprintf("%s/is_blocked?url=%s&proxy=%s", proxyServiceURL, validURL, validProxyAddress)
		LOG.Debug("check proxy by url: ", checkURL)
		e, _, data := net_http.MakeBasicRequest(checkURL)
		So(e, ShouldBeNil)
		So(len(*data) > 1, ShouldBeTrue)
		LOG.Debug(*data)
		So(strings.Contains(*data, "true"), ShouldBeTrue)

	})

	Convey("Run proxy server", t, func() {

		go func() {
			_ = http.ListenAndServe(":11316", proxy)
		}()

	})

	Convey("check unbanning proxy", t, func() {

		count := 0
		for {
			if count > 5 {
				checkURL := fmt.Sprintf("%s/is_blocked?url=%s&proxy=%s", proxyServiceURL, validURL, validProxyAddress)
				err, _, data := net_http.MakeBasicRequest(checkURL)

				LOG.Debug("check proxy by url: ", checkURL)

				So(err, ShouldBeNil)
				So(len(*data) > 1, ShouldBeTrue)
				LOG.Debug(*data)
				So(strings.Contains(*data, "false"), ShouldBeTrue)

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

		e := coll.RemoveBanHistory(endpointURL.Hostname())
		So(e, ShouldBeNil)
		count := coll.CountByBans(endpointURL.Hostname())
		So(count, ShouldEqual, 0)

		err = coll.Block(proxyList[0], endpointURL)
		So(err, ShouldBeNil)

		count = coll.CountByBans(endpointURL.Hostname())
		So(count, ShouldEqual, 1)
		e = coll.RemoveBanHistory(endpointURL.Hostname())
		So(e, ShouldBeNil)

		count = coll.CountByBans(endpointURL.Hostname())
		So(count, ShouldEqual, 0)

	})
}

func TestBlockAllProxies(t *testing.T) {

	coll := (&Collection{
		Settings: settings,
	}).Init()

	_ = coll.RemoveBanHistory(endpointURL.Hostname())

	Convey("Block all proxy", t, func() {
		for _, proxy := range proxyList {
			err := coll.Block(proxy, endpointURL)
			So(err, ShouldBeNil)
		}

		count := coll.CountByBans(endpointURL.Hostname())

		So(count, ShouldEqual, len(proxyList))

	})

	Convey("check a new collection state", t, func() {
		coll2 := (&Collection{
			Settings: settings,
		}).Init()

		So(coll2.CountByBans(endpointURL.Hostname()), ShouldEqual, len(proxyList))

		err := coll.RemoveBanHistory(endpointURL.Hostname())
		So(err, ShouldBeNil)

		count := coll.CountByBans(endpointURL.Hostname())

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

		println(coll.CountByLocked(endpointURL.Hostname()), coll.CountByBans(endpointURL.Hostname()))

		err, proxy := coll.GetProxy(endpointURL)
		So(err, ShouldBeNil)

		st0 := timer.SetTimeout(func(args ...interface{}) {
			println(proxy, " ", coll.IsLocked(endpointURL.Hostname(), proxy))
			c.So(coll.IsLocked(endpointURL.Hostname(), proxy), ShouldBeTrue)
		}, BlockedTime+1)

		st0.Await()

		st := timer.SetTimeout(func(args ...interface{}) {
			println(proxy, " ", coll.IsLocked(endpointURL.Hostname(), proxy))
			c.So(coll.IsLocked(endpointURL.Hostname(), proxy), ShouldBeFalse)
		}, (BlockedTime+1)*1000)

		st.Await()
	})
}

func TestRemoveRedisCollection(t *testing.T) {

	Convey("remove locked proxies", t, func() {
		coll := (&Collection{
			Settings: settings,
		}).Init()

		err := coll.RemoveHistory(endpointURL.Hostname())

		So(err, ShouldBeNil)

		So(coll.CountByLocked(endpointURL.Hostname()), ShouldEqual, 0)

	})

}

func TestGetConcurrentProxyCollection(t *testing.T) {
	Convey("create a new proxy collection", t, func() {
		coll := (&Collection{
			Settings: settings,
		}).Init()
		hostname := endpointURL.Hostname()

		LOG.Debug(hostname)

		_ = coll.RemoveBanHistory(hostname)
		_ = coll.RemoveHistory(hostname)

		g := sync.WaitGroup{}
		g.Add(4)
		resultProxyEnv := make(map[string]string)

		Convey("locking all proxy", func(c C) {

			go func() {
				err, proxy := coll.GetProxy(endpointURL)
				LOG.Debug(err, proxy)
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			go func() {
				err, proxy := coll.GetProxy(endpointURL)
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			go func() {
				err, proxy := coll.GetProxy(endpointURL)
				resultProxyEnv[proxy] = "-"
				c.So(err, ShouldBeNil)
				g.Done()
			}()

			go func() {
				err, proxy := coll.GetProxy(endpointURL)
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
			err := coll.RemoveHistory(hostname)

			So(err, ShouldBeNil)

			So(coll.CountByLocked(endpointURL.Hostname()), ShouldEqual, 0)

		})
	})
}
