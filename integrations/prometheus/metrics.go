package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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
	IsException   bool
}

type Metric struct {
	MetricData

	Type        string   `yaml:"type" json:"type"`
	LabelsKeys  []string `yaml:"labelsKeys" json:"labelsKeys"`
	Description string   `yaml:"description" json:"description"`
}

type TyphoonMetric struct {
	singleton.Singleton

	Metric

	Active bool

	ProjectName    string
	ComponentName  string
	prometheusPath string
}

func (tm *TyphoonMetric) GetPrometheusPath() string {
	tm.Construct(func() {
		projectName := strings.TrimSpace(strings.ReplaceAll(tm.ProjectName, "-", "_"))
		componentName := strings.TrimSpace(strings.ReplaceAll(tm.ComponentName, "-", "_"))

		tm.prometheusPath = strings.Join([]string{projectName, componentName, tm.Name}, "_")
	})
	return tm.prometheusPath
}

type MetricsInterface interface {
	SetException(data MetricData)
	AddNewMetric(metric Metric)
	AddMetric(metric ...Metric)
	Update(data MetricData)
	GetDTO(data MetricData) *dto.Metric
	Add(data MetricData)
	Dec(data MetricData)
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

		m.LOG.Debugf("init metric %s, prometheusPath: %s, labels: %+v", metricName, prometheusPath, metricInfo.LabelsKeys)

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

		metricInfo.Active = true
	}

}

func (m *Metrics) AddNewMetric(metric Metric) {
	m.init()

	if _, ok := m.metrics[metric.Name]; !ok {

		m.LOG.Debugf("add a new metric %s", metric.Name)
		m.metrics[metric.Name] = &TyphoonMetric{
			Metric:        metric,
			ComponentName: m.Config.ComponentName,
			ProjectName:   m.Config.ProjectName,
		}

		exceptionName := strings.Join([]string{metric.Name, "exceptions"}, "_")

		exceptionDescription := strings.Join([]string{metric.Description, " Only exceptions"}, ";")

		exceptionMetric := &TyphoonMetric{
			Metric: Metric{
				Type:        metric.Type,
				Description: exceptionDescription,
				MetricData: MetricData{
					Name: exceptionName,
				},
			},
			ComponentName: m.Config.ComponentName,
			ProjectName:   m.Config.ProjectName,
		}

		switch metric.Type {
		case TypeCounter:
			m.metrics[exceptionName] = exceptionMetric
		case TypeCounterVec:
			exceptionMetric.LabelsKeys = metric.LabelsKeys

			m.metrics[exceptionName] = exceptionMetric
		}
		m.init()

	}

}

func (m *Metrics) AddMetric(metric ...Metric) {
	//for _, _metric := range metric {
	//
	//}
}

func (m *Metrics) SetException(data MetricData) {
	if tm, ok := m.metrics[data.Name]; ok {
		name := strings.Join([]string{tm.GetPrometheusPath(), "exceptions"}, "_")
		switch tm.Type {
		case TypeCounter:
			metric := m.measurer.Counter(name)
			metric.Inc()
		case TypeCounterVec:
			metricVec := m.measurer.CounterVec(name)
			me, err := metricVec.GetMetricWith(data.Labels)

			if err != nil {
				m.LOG.Error(err.Error())
			} else {
				me.Add(1)
			}
		}
	}
}

func (m *Metrics) Update(data MetricData) {

}

func (m *Metrics) Add(data MetricData) {
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

			me, err := metricVec.GetMetricWith(data.Labels)
			if err != nil {
				m.LOG.Error(err.Error())
			} else {
				me.Add(value)
			}
		case TypeGauge:
			metric := m.measurer.Gauge(name)
			metric.Inc()
		case TypeGaugeVec:

			metricVec := m.measurer.GaugeVec(name)
			me, err := metricVec.GetMetricWith(data.Labels)
			if err != nil {
				m.LOG.Error(err.Error())
			} else {
				me.Inc()
			}

		}
	}
}

func (m *Metrics) Dec(data MetricData) {
	if tm, ok := m.metrics[data.Name]; ok {

		name := tm.GetPrometheusPath()

		switch tm.Type {
		case TypeGauge:
			metric := m.measurer.Gauge(name)
			metric.Dec()
		case TypeGaugeVec:
			metricVec := m.measurer.GaugeVec(name)
			metric, err := metricVec.GetMetricWith(data.Labels)

			if err != nil {
				m.LOG.Error(err.Error())
			} else {
				metric.Dec()
			}
		}
	}
}

func (m *Metrics) GetDTO(data MetricData) *dto.Metric {
	if tm, ok := m.metrics[data.Name]; ok {
		o := &dto.Metric{}
		path := tm.GetPrometheusPath()
		if data.IsException {
			path = strings.Join([]string{path, "exceptions"}, "_")
		}
		switch tm.Type {
		case TypeSummaryVec:
			//c, e := m.measurer.SummaryVec(path).GetMetricWith(data.Labels).(prometheus.Summary)
			//if e != false {
			//	m.LOG.Error("SummaryVec is not ")
			//}
			//_ = c.Write(o)
		case TypeSummary:
			_ = m.measurer.Summary(path).Write(o)
		case TypeCounterVec:
			c, e := m.measurer.CounterVec(path).GetMetricWith(data.Labels)

			if e != nil {
				m.LOG.Error(e.Error())
			} else {
				_ = c.Write(o)
			}

		case TypeCounter:
			_ = m.measurer.Counter(path).Write(o)
		case TypeGaugeVec:

			c, e := m.measurer.GaugeVec(path).GetMetricWith(data.Labels)

			if e != nil {
				m.LOG.Error(e.Error())
			} else {
				_ = c.Write(o)
			}
		case TypeGauge:
			_ = m.measurer.Gauge(path).Write(o)
		}

		return o

	}
	return nil
}
