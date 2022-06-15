package capturer

import (
	"bytes"
	"image/jpeg"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init()  {
	log.InitD()
}

type Capturer struct {
	singleton.Singleton
	awaitabler.Object

	Count int
	Active bool
	Quality int
	Delay float32
	IsFullScreen bool
	Output chan []byte

	LOG interfaces.LoggerInterface
}

func (c *Capturer) Stop() *Capturer {
	c.LOG.Debug("stop capturing")
	c.Active = false
	return c
}

func (c *Capturer) Capture() *Capturer {
	c.Construct(func() {
		c.Active = true
		c.Add()
		c.Output = make(chan []byte)
		c.LOG = log.New(log.D{"agent": "Capture img"})
		c.LOG.Debug("init")

		go func() {
			defer close(c.Output)
			c.LOG.Debug("start ")
			for {
				if !c.Active { return }
				n := screenshot.NumActiveDisplays()
				for i := 0; i < n; i++ {
					for {

						c.LOG.Debug("get img")
						bounds := screenshot.GetDisplayBounds(i)
						imgC, err := screenshot.CaptureRect(bounds)

						if err != nil {
							panic(err)
						}

						var buff bytes.Buffer
						jpeg.Encode(&buff, imgC, &jpeg.Options{Quality: c.Quality})
						imgBytes := buff.Bytes()
						c.LOG.Debug("send to channel:", c.Count)
						c.Output <- imgBytes

						time.Sleep(time.Duration(c.Delay) * time.Second)

						c.Count ++
					}

				}
			}

		}()


	})

	return c
}
