package main

import (
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

func NewSystemairCollector(hvac *modbus.ModbusClient) *SystemairCollector {
	return &SystemairCollector{
		hvac: hvac,
		temp_mode_enabled: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "hvac_temp_mode_enabled",
					Help: "",
			}, []string{"supply", "room", "extract"}),
		temp_controller_percentage: prometheus.NewGauge(prometheus.GaugeOpts{
					Name: "hvac_temp_controller_percentage",
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

	e.temp_controller_percentage.Set(float64(systemairmodbus.GetTempDemandPercentage(e.hvac)))

	// Collect metrics to Prometheus
	e.temp_mode_enabled.Collect(ch)
	e.temp_controller_percentage.Collect(ch)
}