package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simonvetter/modbus"
)

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