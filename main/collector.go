package main

import (
	"strings"
	"sync"
	"systemair-prom-exporter-go/systemairmodbus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type SystemairCollector struct {
	// Ensure only a single Collect() process can take place
	// Otherwise, we might run into issues with the serial modbus interface
	mutex sync.Mutex

	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	temp_mode_enabled *prometheus.GaugeVec
	temp_controller_percentage prometheus.Gauge
}

func NewSystemairCollector(hvac *modbus.ModbusClient, namespace string) *SystemairCollector {
	return &SystemairCollector{
		hvac: hvac,
		temp_mode_enabled: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: namespace,
					Subsystem: "temp",
					Name: "mode_enabled",
					Help: "Unit temperature control mode. The currently enabled mode has a value of 1",
			}, []string{"mode"}),
		temp_controller_percentage: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: namespace,
					Subsystem: "temp",
					Name: "controller_percentage",
					Help: "\"Output of the SATC\" in percentage. Min 0, Max 100",
			}),
	}
}

func (e *SystemairCollector) Describe(ch chan<- *prometheus.Desc) {
	// Register all metrics with Prometheus
	e.temp_mode_enabled.Describe(ch)
	e.temp_controller_percentage.Describe(ch)
}

func (e *SystemairCollector) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // Ensure a single collection at a time
	defer e.mutex.Unlock()

	for _, labelValue := range []string{"supply", "room", "extract"} {
			e.temp_mode_enabled.WithLabelValues(labelValue).Set(0)
	}
	e.temp_mode_enabled.WithLabelValues(strings.ToLower(systemairmodbus.GetTempMode(e.hvac))).Set(1)


	e.temp_controller_percentage.Set(float64(systemairmodbus.GetTempDemandPercentage(e.hvac)))

	// Collect metrics to Prometheus
	e.temp_mode_enabled.Collect(ch)
	e.temp_controller_percentage.Collect(ch)
}