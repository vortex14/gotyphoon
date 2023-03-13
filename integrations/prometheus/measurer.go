package prometheus

import "github.com/prometheus/client_golang/prometheus"

type Measurer interface {
	Summary(name string) *prometheus.SummaryVec
	AddSummary(name, description string, labels ...string)

	AddCounter(name, description string, labels ...string)
	Counter(name string) *prometheus.CounterVec

	AddGauge(name, description string, labels ...string)
	Gauge(name string) *prometheus.GaugeVec
}

type measurer struct {
	summary map[string]*prometheus.SummaryVec
	counter map[string]*prometheus.CounterVec
	gauge   map[string]*prometheus.GaugeVec

	runtimeMetricsCollector *runtimeMetricsCollector
}

func (m *measurer) AddSummary(name, description string, labels ...string) {
	m.summary[name] = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: name,
	}, labels)
}

func (m *measurer) Summary(name string) *prometheus.SummaryVec {
	return m.summary[name]
}

func (m *measurer) AddCounter(name, description string, labels ...string) {
	m.counter[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
	}, labels)
}

func (m *measurer) Counter(name string) *prometheus.CounterVec {
	return m.counter[name]
}

func (m *measurer) AddGauge(name, description string, labels ...string) {
	m.gauge[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
	}, labels)
}

func (m *measurer) Gauge(name string) *prometheus.GaugeVec {
	return m.gauge[name]
}

func NewMeasurer(config MetricsConfig) Measurer {
	m := &measurer{
		summary: make(map[string]*prometheus.SummaryVec),
		counter: make(map[string]*prometheus.CounterVec),
		gauge:   make(map[string]*prometheus.GaugeVec),
	}

	// Run collect runtime metrics
	_runtimeMetricsCollector := &runtimeMetricsCollector{
		config:   config,
		measurer: m,
	}
	go _runtimeMetricsCollector.Run()

	return m
}
