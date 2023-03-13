package prometheus

import (
	dto "github.com/prometheus/client_model/go"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMetricCounter(t *testing.T) {

	Convey("counter", t, func() {
		m := NewMeasurer(Metrics{})

		m.AddCounter("my_counter", "my_counter description")
		
		c := m.Counter("my_counter")

		c.Inc()
		c.Inc()
		c.Inc()

		o := &dto.Metric{}

		e := c.Write(o)

		So(e, ShouldBeNil)
		So(*o.Counter.Value, ShouldEqual, 3)
		So(o.Label, ShouldBeEmpty)

	})
}
