package prometheus

import (
	"github.com/vortex14/gotyphoon/elements/models/singleton"
)

type RunTime struct {
	Duration string `yaml:"duration" json:"duration"`
	CPU      bool   `yaml:"cpu" json:"cpu"`
	Mem      bool   `yaml:"mem" json:"mem"`
	GC       bool   `yaml:"gc" json:"gc"`
}

type MetricData struct {
	Name   string   `yaml:"name" json:"name"`
	Labels []string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Value  float32  `yaml:"value" json:"value"`

	AutoIncrement bool
	AutoDecrement bool
}

type Metric struct {
	MetricData
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
}

type TyphoonMetric struct {
	Metric

	Active bool

	ProjectName    string
	ComponentName  string
	PrometheusPath string
}

type MetricsInterface interface {
	SetException(data *MetricData)
	AddNewMetric(metric *Metric)
	AddMetric(metric ...*Metric)
	Update(data *MetricData)
	Add(data *MetricData)
	Dec(data *MetricData)
}

type MetricsConfig struct {
	Runtime RunTime `yaml:"runtime" json:"runtime"`
}

type Metrics struct {
	singleton.Singleton
	Config   MetricsConfig
	measurer *measurer

	//List map[string]TyphoonMetric `yaml:"list" json:"list"`
}

func (tm *Metrics) init() {
	tm.Construct(func() {

	})
}

func (tm *Metrics) AddNewMetric(metric *Metric) {

}

func (tm *Metrics) AddMetric(metric ...*Metric) {
	//for _, _metric := range metric {
	//
	//}
}

func (tm *Metrics) SetException(data *MetricData) {

}

func (tm *Metrics) Update(data *MetricData) {

}

func (tm *Metrics) Add(data *MetricData) {

}

func (tm *Metrics) Dec(data *MetricData) {

}
