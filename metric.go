package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricServer struct {
	mChan <-chan float64
}

func NewMetricServer(mChan <-chan float64) *MetricServer {
	return &MetricServer{mChan}
}


var (
	sshConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "awesome_ssh_tarpit_open_connections",
		Help: "The amount of clients in the ssh tarpit",
	})
)


func (m MetricServer) listenForMetric() {
	for m := range m.mChan {
		sshConnections.Add(m)
	}
}

func (m MetricServer) ListenAndServe(addr string) {
	go m.listenForMetric()
	h := promhttp.Handler()
	http.Handle("/", h)
	http.Handle("/metrics", h)
	fmt.Printf("Starting httpd, binding on %s \n", addr)
	http.ListenAndServe(addr, nil)
}
