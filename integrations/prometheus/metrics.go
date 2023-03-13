package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"strings"
)

type RunTime struct {
	Duration string `yaml:"duration" json:"duration"`
	CPU      bool   `yaml:"cpu" json:"cpu"`
	Mem      bool   `yaml:"mem" json:"mem"`
	GC       bool   `yaml:"gc" json:"gc"`
}

type MetricData struct {
	Name   string            `yaml:"name" json:"name"`
	Labels prometheus.Labels `yaml:"labels,omitempty" json:"labels,omitempty"`
	Value  float64           `yaml:"value" json:"value"`

	AutoIncrement bool
	AutoDecrement bool
}

type Metric struct {
	*MetricData

	Type        string   `yaml:"type" json:"type"`
	LabelsKeys  []string `yaml:"labelsKeys", json:"labelsKeys"`
	Description string   `yaml:"description" json:"description"`
}

type TyphoonMetric struct {
	singleton.Singleton

	*Metric

	Active bool

	ProjectName    string
	ComponentName  string
	prometheusPath string
}

func (tm *TyphoonMetric) GetPrometheusPath() string {
	tm.Construct(func() {

		projectName := strings.ReplaceAll(tm.ProjectName, "-", "_")
		componentName := strings.ReplaceAll(tm.ComponentName, "-", "_")

		tm.prometheusPath = strings.Join([]string{projectName, componentName, tm.Name}, "_")
	})
	return tm.prometheusPath
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
	LOG interfaces.LoggerInterface

	singleton.Singleton
	Config   MetricsConfig
	metrics  map[string]*TyphoonMetric
	measurer Measurer
}

func (m *Metrics) init() {
	m.Construct(func() {
		m.LOG = log.New(map[string]interface{}{"metrics": "prometheus"})
		m.metrics = make(map[string]*TyphoonMetric)
		m.measurer = NewMeasurer(m.Config)
	})

	for metricName, metricInfo := range m.metrics {

		if metricInfo.Active {
			continue
		}
		prometheusPath := metricInfo.GetPrometheusPath()

		m.LOG.Debugf("init metric %s, prometheusPath: %s", metricName, prometheusPath)

		var collector prometheus.Collector

		switch metricInfo.Type {
		case TypeSummaryVec:
			m.measurer.AddSummaryVec(prometheusPath, metricInfo.Description, metricInfo.LabelsKeys...)
			collector = m.measurer.SummaryVec(prometheusPath)
		case TypeSummary:
			m.measurer.AddSummary(prometheusPath, metricInfo.Description)
			collector = m.measurer.Summary(prometheusPath)
		case TypeCounterVec:
			m.measurer.AddCounterVec(prometheusPath, metricInfo.Description, metricInfo.LabelsKeys...)
			collector = m.measurer.CounterVec(prometheusPath)
		case TypeCounter:
			m.measurer.AddCounter(prometheusPath, metricInfo.Description)
			collector = m.measurer.Counter(prometheusPath)
		case TypeGaugeVec:
			m.measurer.AddGaugeVec(prometheusPath, metricInfo.Description, metricInfo.LabelsKeys...)
			collector = m.measurer.GaugeVec(prometheusPath)
		case TypeGauge:
			m.measurer.AddGauge(prometheusPath, metricInfo.Description)
			collector = m.measurer.Gauge(prometheusPath)
		}

		prometheus.MustRegister(collector)
	}

}

func (m *Metrics) AddNewMetric(metric *Metric) {
	m.init()

	if _, ok := m.metrics[metric.Name]; !ok {

		m.LOG.Debugf("add a new metric %s", metric.Name)
		m.metrics[metric.Name] = &TyphoonMetric{
			Metric:        metric,
			ComponentName: m.Config.ComponentName,
			ProjectName:   m.Config.ProjectName,
		}
		m.init()

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
	if tm, ok := m.metrics[data.Name]; ok {

		value := tm.Value

		if value == 0 {
			value = 1
		}

		name := tm.GetPrometheusPath()

		switch tm.Type {
		case TypeCounter:
			metric := m.measurer.Counter(name)
			metric.Add(value)
		case TypeCounterVec:
			metricVec := m.measurer.CounterVec(name)
			metric, err := metricVec.GetMetricWith(tm.Labels)

			if err != nil {
				m.LOG.Error(err.Error())
			} else {
				metric.Add(value)
			}
		}
	}
}

func (m *Metrics) Dec(data *MetricData) {

}
