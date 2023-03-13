package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	PR "github.com/vortex14/gotyphoon/integrations/prometheus"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var measurer = PR.NewMeasurer(PR.Metrics{})

const (
	nameCounter        = "test_counter"
	descriptionCounter = "test description"
)

func ping(w http.ResponseWriter, req *http.Request) {
	c := measurer.Counter(nameCounter)
	c.Inc()

	fmt.Fprintf(w, "pong")
}

func main() {

	measurer.AddCounter(nameCounter, descriptionCounter)
	c := measurer.Counter(nameCounter)

	fmt.Printf("%+v", c)

	prometheus.MustRegister(c)

	http.HandleFunc("/ping", ping)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8090", nil)
}
