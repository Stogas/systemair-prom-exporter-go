package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simonvetter/modbus"
)

// StartExporter, given an address and path, creates Prometheus Collectors,
//   registers the to the default Prometheus Registerer,
//   registers a Prometheus HTTP handler,
//   and starts an HTTP server.
func StartExporter(addr string, path string, m *modbus.ModbusClient) {
	collectorTemp := NewSystemairTempCollector(m, "hvac")
	prometheus.MustRegister(collectorTemp)
	collectorAirflow := NewSystemairAirflowCollector(m, "hvac")
	prometheus.MustRegister(collectorAirflow)
	collectorHeater := NewSystemairHeaterCollector(m, "hvac")
	prometheus.MustRegister(collectorHeater)
	collectorMisc := NewSystemairMiscCollector(m, "hvac")
	prometheus.MustRegister(collectorMisc)

	http.Handle(path, promhttp.Handler())

	fmt.Printf("Starting HTTP listener on: %v\n", addr + path)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("HTTP listener returned error: %v\n", err)
	}
}