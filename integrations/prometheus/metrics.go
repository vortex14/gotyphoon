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
	*Metric

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
	Runtime                    RunTime `yaml:"runtime" json:"runtime"`
	ProjectName, ComponentName string
}

type Metrics struct {
	singleton.Singleton
	Config   MetricsConfig
	metrics  map[string]*TyphoonMetric
	measurer Measurer

	//List map[string]TyphoonMetric `yaml:"list" json:"list"`
}

func (m *Metrics) init() {
	m.Construct(func() {
		m.metrics = make(map[string]*TyphoonMetric)
		m.measurer = NewMeasurer(m.Config)
	})

	for metricName, metricInfo := range m.metrics {

		if metricInfo.Active {
			continue
		}

		switch metricInfo.Type {
		case TypeSummaryVec:
			m.measurer.AddSummaryVec(metricName, metricInfo.Description, metricInfo.Labels...)
		case TypeSummary:
			m.measurer.AddSummary(metricName, metricInfo.Description)
		case TypeCounterVec:
			m.measurer.AddCounterVec(metricName, metricInfo.Description, metricInfo.Labels...)
		case TypeCounter:
			m.measurer.AddCounter(metricName, metricInfo.Description)
		case TypeGaugeVec:
			m.measurer.AddGaugeVec(metricName, metricInfo.Description, metricInfo.Labels...)
		case TypeGauge:
			m.measurer.AddGauge(metricName, metricInfo.Description)
		}
	}

}

func (m *Metrics) AddNewMetric(metric *Metric) {
	m.init()

	if _, ok := m.metrics[metric.Name]; !ok {

		m.metrics[metric.Name] = &TyphoonMetric{
			Metric:        metric,
			ComponentName: m.Config.ComponentName,
			ProjectName:   m.Config.ProjectName,
		}

	}

}

func (m *Metrics) AddMetric(metric ...*Metric) {
	//for _, _metric := range metric {
	//
	//}
}

func (m *Metrics) SetException(data *MetricData) {

}

func (m *Metrics) Update(data *MetricData) {

}

func (m *Metrics) Add(data *MetricData) {

}

func (m *Metrics) Dec(data *MetricData) {

}
