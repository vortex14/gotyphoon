package prometheus

import (
	"context"
	"sync"
	"testing"
	"time"

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

		md := MetricData{Name: nameCounter}

		newCountMetric := Metric{
			Type: TypeCounter, Description: descriptionCounter,
			MetricData: md,
		}

		metrics.AddNewMetric(newCountMetric)

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

func TestGauge(t *testing.T) {
	Convey("test basic gauge", t, func() {

		var metrics = Metrics{
			Config: MetricsConfig{ProjectName: "my-project", ComponentName: "component-1"},
		}

		const (
			nameGauge        = "test_gauge"
			descriptionGauge = "test gauge description"
		)

		md := MetricData{Name: nameGauge}

		newGaugeMetric := Metric{
			Type: TypeGauge, Description: descriptionGauge,
			MetricData: md,
		}

		metrics.AddNewMetric(newGaugeMetric)

		metrics.Add(md)
		metrics.Add(md)
		metrics.Add(md)

		_dto := metrics.GetDTO(md)

		So(*_dto.Gauge.Value, ShouldEqual, 3)

		metrics.Dec(md)
		metrics.Dec(md)
		metrics.Dec(md)

		_dto = metrics.GetDTO(md)

		So(*_dto.Gauge.Value, ShouldEqual, 0)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func(md MetricData) {
			metrics.Add(md)

			time.Sleep(5 * time.Second)

			metrics.Dec(md)

			wg.Done()
		}(md)

		go func(md MetricData) {
			metrics.Add(md)

			time.Sleep(2 * time.Second)

			metrics.Dec(md)
			wg.Done()
		}(md)

		time.Sleep(1 * time.Second)

		_dto = metrics.GetDTO(md)

		So(*_dto.Gauge.Value, ShouldEqual, 2)

		wg.Wait()

		_dto = metrics.GetDTO(md)

		So(*_dto.Gauge.Value, ShouldEqual, 0)

	})
}

func TestMetricsCounterWithExceptionsAndLabels(t *testing.T) {

	Convey("metrics counter with exceptions and labels", t, func() {

		var metrics = Metrics{
			Config: MetricsConfig{ProjectName: "my-project", ComponentName: "component-1"},
		}

		const (
			nameCounter        = "test_counter"
			descriptionCounter = "test description"
		)

		newCountMetric := Metric{
			Type: TypeCounterVec, Description: descriptionCounter,
			MetricData: MetricData{Name: nameCounter},
			LabelsKeys: []string{"label1", "label2"},
		}

		metrics.AddNewMetric(newCountMetric)

		md := MetricData{Name: nameCounter, Labels: map[string]string{"label1": "1", "label2": "2"}}

		metrics.Add(md)
		metrics.Add(md)
		metrics.Add(md)

		metrics.SetException(md)
		metrics.SetException(md)

		_dto := metrics.GetDTO(md)

		_dtoE := metrics.GetDTO(MetricData{
			Name:        nameCounter,
			IsException: true,
			Labels:      map[string]string{"label1": "1", "label2": "2"},
		})

		So(*_dto.Counter.Value, ShouldEqual, 3)

		So(*_dtoE.Counter.Value, ShouldEqual, 2)

	})
}

func TestGaugeWithLabels(t *testing.T) {
	Convey("test gauge with labels", t, func() {

		var metrics = Metrics{
			Config: MetricsConfig{ProjectName: "my-project", ComponentName: "component-1"},
		}

		const (
			nameGauge        = "test_gauge_labels"
			descriptionGauge = "test description"
		)

		newGaugeMetric := Metric{
			Type: TypeGaugeVec, Description: descriptionGauge,
			MetricData: MetricData{Name: nameGauge},
			LabelsKeys: []string{"label1", "label2"},
		}

		metrics.AddNewMetric(newGaugeMetric)

		md := MetricData{Name: nameGauge, Labels: map[string]string{"label1": "1", "label2": "2"}}

		metrics.Add(md)
		metrics.Add(md)
		metrics.Add(md)

		_dto := metrics.GetDTO(md)
		So(*_dto.Gauge.Value, ShouldEqual, 3)

		metrics.Dec(md)
		metrics.Dec(md)
		metrics.Dec(md)

		_dto = metrics.GetDTO(md)
		So(*_dto.Gauge.Value, ShouldEqual, 0)

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
