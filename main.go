package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

func main() {

	//Create a new instance of the foocollector and
	//register it with the prometheus client.
	foo := newFooCollector()
	prometheus.MustRegister(foo)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
