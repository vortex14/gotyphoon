package chromedp

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

type ChromeHeadless struct {
	singleton.Singleton

	LOG interfaces.LoggerInterface

	Proxy  string
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *ChromeHeadless) init() {
	c.Construct(func() {
		var cx context.Context
		c.LOG = log.New(log.D{"service": "chromeless"})
		if len(c.Proxy) > 0 {
			o := append(chromedp.DefaultExecAllocatorOptions[:],
				//... any options here
				chromedp.ProxyServer(c.Proxy),
			)

			cx, _ = chromedp.NewExecAllocator(context.Background(), o...)
		} else {
			cx, _ = chromedp.NewContext(context.Background())
		}

		ctx, cancel := chromedp.NewContext(cx)
		c.ctx = ctx
		c.cancel = cancel
	})
}

func (c *ChromeHeadless) Close() {
	c.Destruct(func() {
		c.cancel()
	})
}

// setheaders returns a task list that sets the passed headers.
func (c *ChromeHeadless) setHeaders(host string, headers map[string]interface{}, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(headers),
		chromedp.Navigate(host),
		chromedp.Text(`//*`, res),
	}
}

func (c *ChromeHeadless) Request(url string) string {
	if c.ctx == nil {
		c.init()
	}
	// run task list
	var res string
	err := chromedp.Run(c.ctx, c.setHeaders(url, map[string]interface{}{},
		&res,
	))
	if err != nil {
		c.LOG.Error(err.Error())
	}

	return res
}
