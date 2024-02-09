package main

import (
	"fmt"
	"net/http"
	"sync"
	"systemair-prom-exporter-go/systemairmodbus"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simonvetter/modbus"
)

type SystemairCollector struct {
	// Ensure only a single Collect() process can take place
	// Otherwise, we might run into issues with the serial modbus interface
	mutex sync.Mutex

	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	// Metrics for each library function
	temp_controller_percentage prometheus.Gauge
}

func NewSystemairCollector(hvac *modbus.ModbusClient) *SystemairCollector {
	return &SystemairCollector{
		hvac: hvac,
		temp_controller_percentage: prometheus.NewGauge(prometheus.GaugeOpts{
					Name: "hvac_temp_controller_percentage",
					Help: "\"Output of the SATC\" in percentage. Min 0, Max 100",
			}),
	}
}

func (e *SystemairCollector) Describe(ch chan<- *prometheus.Desc) {
	// Register all metrics with Prometheus
	e.temp_controller_percentage.Describe(ch)
}

func (e *SystemairCollector) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // Ensure a single collection at a time
	defer e.mutex.Unlock()

	// Call your library's functions sequentially and set the metrics
	valueA := systemairmodbus.GetTempDemandPercentage(e.hvac) // Assume this returns a uint16
	e.temp_controller_percentage.Set(float64(valueA))

	// Collect metrics to Prometheus
	e.temp_controller_percentage.Collect(ch)
}

// StartExporter, given an address and path, creates a Prometheus Collector,
//   registers it to the default Prometheus Registerer,
//   registers a Prometheus HTTP handler,
//   and starts an HTTP server.
func StartExporter(addr string, path string, hvac *modbus.ModbusClient) {
	collector := NewSystemairCollector(hvac)
	prometheus.MustRegister(collector)

	http.Handle(path, promhttp.Handler())

	// Create a channel to communicate with the goroutine.
	errChan := make(chan error)

	// Start the HTTP server in a new goroutine.
	go func() {
			// ListenAndServe always returns a non-nil error.
			err := http.ListenAndServe(addr, nil)
			if err != nil {
					// Send any errors back through the channel.
					errChan <- err
			}
	}()

	// Give the server a moment to start.
	time.Sleep(100 * time.Millisecond)

	// Check if there was an error starting the server.
	select {
	case err := <-errChan:
			// Handle the error, e.g., log it or exit.
			fmt.Println("Failed to start HTTP listener:", err)
	default:
			// If no error, server started successfully.
			fmt.Printf("HTTP listener is successfully serving on: %v\n", addr + path)
	}

	// Keep the main goroutine alive indefinitely.
	select {}
}