package prometheus

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	dto "github.com/prometheus/client_model/go"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func TestMetricCounter(t *testing.T) {

	Convey("basic counter", t, func() {
		m := NewMeasurer(MetricsConfig{})

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

func TestMetricsCounterWithExceptions(t *testing.T) {

	Convey("metrics counter with exceptions", t, func() {

		var metrics = Metrics{
			Config: MetricsConfig{ProjectName: "my-project", ComponentName: "component-1"},
		}

		const (
			nameCounter        = "test_counter"
			descriptionCounter = "test description"
		)

		newCountMetric := Metric{
			Type: TypeCounter, Description: descriptionCounter,
			MetricData: MetricData{Name: nameCounter},
		}

		metrics.AddNewMetric(newCountMetric)

		md := MetricData{Name: nameCounter}

		metrics.Add(md)
		metrics.Add(md)
		metrics.Add(md)

		metrics.SetException(md)
		metrics.SetException(md)

		_dto := metrics.GetDTO(md)

		_dtoE := metrics.GetDTO(MetricData{
			Name:        nameCounter,
			IsException: true,
		})

		So(*_dto.Counter.Value, ShouldEqual, 3)

		So(*_dtoE.Counter.Value, ShouldEqual, 2)

	})
}

func TestPipelineCounterMetric(t *testing.T) {
	Convey("pipeline counter", t, func() {

		L := log.New(map[string]interface{}{"logger": "pipeline"})
		ctx := log.NewCtx(context.Background(), L)

		p := &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name:        "new pipeline",
				Description: "some pipeline",
			},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				return nil, ctx
			},
		}

		var Err error

		for i := 0; i < 5; i++ {
			p.Run(ctx, func(pipeline interfaces.BasePipelineInterface, err error) {
				Err = err
			}, func(ctx context.Context) {

			})
		}

		So(Err, ShouldBeNil)

	})
}
