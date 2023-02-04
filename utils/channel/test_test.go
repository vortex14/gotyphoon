package channel

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/log"
	"sync"
	"testing"
	"time"
)

func init() {
	log.InitD()
}

func TestForCh(t *testing.T) {
	LOG := log.New(map[string]interface{}{"ch": "read"})

	Convey("test ch", t, func(c C) {
		ch := make(chan bool, 10)

		ch <- true
		ch <- false

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func(wg *sync.WaitGroup) {
			count := 0
			for b := range ch {
				LOG.Info("received: ", b)

				count += 1

			}

			c.So(count, ShouldEqual, 2)
			wg.Done()
		}(&wg)

		go func(wg *sync.WaitGroup) {
			select {
			case <-time.After(time.Second * 3):
				close(ch)
				wg.Done()
				break
			}

		}(&wg)

		wg.Wait()

	})

}
