package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Measurer interface {
	AddSummaryVec(name string, labels ...string)
	SummaryVec(name string) *prometheus.SummaryVec
	AddSummary(name string)
	Summary(name string) prometheus.Summary

	AddCounterVec(name string, labels ...string)
	CounterVec(name string) *prometheus.CounterVec
	AddCounter(name string)
	Counter(name string) prometheus.Counter

	AddGaugeVec(name string, labels ...string)
	GaugeVec(name string) *prometheus.GaugeVec
	AddGauge(name string)
	Gauge(name string) prometheus.Gauge
}

type Metrics struct {
	Runtime struct {
		Duration string `yaml:"duration" json:"duration" mapstructure:"duration"`
		CPU      bool   `yaml:"cpu" json:"cpu" mapstructure:"cpu"`
		Mem      bool   `yaml:"mem" json:"mem" mapstructure:"mem"`
		GC       bool   `yaml:"gc" json:"gc" mapstructure:"gc"`
	} `yaml:"runtime" json:"runtime" mapstructure:"runtime"`
	List map[string]struct {
		Type   string   `yaml:"type" json:"type" mapstructure:"type"`
		Labels []string `yaml:"labels,omitempty" json:"labels,omitempty" mapstructure:"labels,omitempty"`
	} `yaml:"list" json:"list" mapstructure:"list"`
}

type measurer struct {
	summaryVec map[string]*prometheus.SummaryVec
	summary    map[string]prometheus.Summary
	counterVec map[string]*prometheus.CounterVec
	counter    map[string]prometheus.Counter
	gaugeVec   map[string]*prometheus.GaugeVec
	gauge      map[string]prometheus.Gauge

	runtimeMetricsCollector *runtimeMetricsCollector
}

func (m *measurer) AddSummaryVec(name string, labels ...string) {
	m.summaryVec[name] = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: name,
	}, labels)
}

func (m *measurer) SummaryVec(name string) *prometheus.SummaryVec {
	return m.summaryVec[name]
}

func (m *measurer) AddSummary(name string) {
	m.summary[name] = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: name,
	})
}

func (m *measurer) Summary(name string) prometheus.Summary {
	return m.summary[name]
}

func (m *measurer) AddCounterVec(name string, labels ...string) {
	m.counterVec[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
	}, labels)
}

func (m *measurer) CounterVec(name string) *prometheus.CounterVec {
	return m.counterVec[name]
}

func (m *measurer) AddCounter(name string) {
	m.counter[name] = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
	})
}

func (m *measurer) Counter(name string) prometheus.Counter {
	return m.counter[name]
}

func (m *measurer) AddGaugeVec(name string, labels ...string) {
	m.gaugeVec[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
	}, labels)
}

func (m *measurer) GaugeVec(name string) *prometheus.GaugeVec {
	return m.gaugeVec[name]
}

func (m *measurer) AddGauge(name string) {
	m.gauge[name] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
	})
}

func (m *measurer) Gauge(name string) prometheus.Gauge {
	return m.gauge[name]
}

func NewMeasurer(config Metrics) Measurer {
	m := &measurer{
		summaryVec: make(map[string]*prometheus.SummaryVec),
		summary:    make(map[string]prometheus.Summary),
		counterVec: make(map[string]*prometheus.CounterVec),
		counter:    make(map[string]prometheus.Counter),
		gaugeVec:   make(map[string]*prometheus.GaugeVec),
		gauge:      make(map[string]prometheus.Gauge),
	}

	for metricName, metricInfo := range config.List {
		switch metricInfo.Type {
		case TypeSummaryVec:
			m.AddSummaryVec(metricName, metricInfo.Labels...)
		case TypeSummary:
			m.AddSummary(metricName)
		case TypeCounterVec:
			m.AddCounterVec(metricName, metricInfo.Labels...)
		case TypeCounter:
			m.AddCounter(metricName)
		case TypeGaugeVec:
			m.AddGaugeVec(metricName, metricInfo.Labels...)
		case TypeGauge:
			m.AddGauge(metricName)
		}
	}

	// Run collect runtime metrics
	_runtimeMetricsCollector := &runtimeMetricsCollector{
		config:   config,
		measurer: m,
	}
	go _runtimeMetricsCollector.Run()

	return m
}
