package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Measurer interface {
	AddSummaryVec(name, description string, labels ...string)
	SummaryVec(name string) *prometheus.SummaryVec
	AddSummary(name, description string)
	Summary(name string) prometheus.Summary

	AddCounterVec(name, description string, labels ...string)
	CounterVec(name string) *prometheus.CounterVec
	AddCounter(name, description string)
	Counter(name string) prometheus.Counter

	AddGaugeVec(name, description string, labels ...string)
	GaugeVec(name string) *prometheus.GaugeVec
	AddGauge(name, description string)
	Gauge(name string) prometheus.Gauge
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

func (m *measurer) AddSummaryVec(name, description string, labels ...string) {
	m.summaryVec[name] = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: name,
		Help: description,
	}, labels)
}

func (m *measurer) SummaryVec(name string) *prometheus.SummaryVec {
	return m.summaryVec[name]
}

func (m *measurer) AddSummary(name, description string) {
	m.summary[name] = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: name,
		Help: description,
	})
}

func (m *measurer) Summary(name string) prometheus.Summary {
	return m.summary[name]
}

func (m *measurer) AddCounterVec(name, description string, labels ...string) {
	m.counterVec[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: description,
	}, labels)
}

func (m *measurer) CounterVec(name string) *prometheus.CounterVec {
	return m.counterVec[name]
}

func (m *measurer) AddCounter(name, description string) {
	m.counter[name] = prometheus.NewCounter(prometheus.CounterOpts{
		Help: description,
		Name: name,
	})
}

func (m *measurer) Counter(name string) prometheus.Counter {
	return m.counter[name]
}

func (m *measurer) AddGaugeVec(name, description string, labels ...string) {
	m.gaugeVec[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	}, labels)
}

func (m *measurer) GaugeVec(name string) *prometheus.GaugeVec {
	return m.gaugeVec[name]
}

func (m *measurer) AddGauge(name, description string) {
	m.gauge[name] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	})
}

func (m *measurer) Gauge(name string) prometheus.Gauge {
	return m.gauge[name]
}

func NewMeasurer(config MetricsConfig) Measurer {
	m := &measurer{
		summaryVec: make(map[string]*prometheus.SummaryVec),
		summary:    make(map[string]prometheus.Summary),
		counterVec: make(map[string]*prometheus.CounterVec),
		counter:    make(map[string]prometheus.Counter),
		gaugeVec:   make(map[string]*prometheus.GaugeVec),
		gauge:      make(map[string]prometheus.Gauge),
	}

	// Run collect runtime metrics
	_runtimeMetricsCollector := &runtimeMetricsCollector{
		config:   config,
		measurer: m,
	}
	go _runtimeMetricsCollector.Run()

	return m
}
