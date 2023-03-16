package main

import (
	"fmt"
	PR "github.com/vortex14/gotyphoon/integrations/prometheus"
	"github.com/vortex14/gotyphoon/log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var metrics = PR.Metrics{
	Config: PR.MetricsConfig{ProjectName: "my-project", ComponentName: "component-1"},
}

const (
	nameCounter        = "test_counter"
	descriptionCounter = "test description"
)

func init() {
	log.InitD()
}

func ping(w http.ResponseWriter, req *http.Request) {
	//c := measurer.Counter(nameCounter)
	//c.Inc()

	metrics.Add(PR.MetricData{Name: nameCounter, Labels: map[string]string{"label1": "test1", "be": "asf"}})

	metrics.Add(PR.MetricData{Name: nameCounter, Labels: map[string]string{"label1": "l", "be": "1"}})
	metrics.Add(PR.MetricData{Name: nameCounter, Labels: map[string]string{"label1": "l", "be": "1"}})
	//metrics.Add(PR.MetricData{Name: nameCounter, Labels: map[string]string{"label2": "test1_test2"}})

	fmt.Fprintf(w, "pong")
}

func main() {

	newCountMetric := PR.Metric{
		Type: PR.TypeCounterVec, Description: descriptionCounter,
		MetricData: PR.MetricData{Name: nameCounter},
		LabelsKeys: []string{"label1", "be"},
	}

	metrics.AddNewMetric(newCountMetric)

	//fmt.Printf("%+v", c)

	//prometheus.MustRegister(c)

	http.HandleFunc("/ping", ping)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8090", nil)
}
