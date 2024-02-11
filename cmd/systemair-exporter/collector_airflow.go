package main

import (
	"systemair-prom-exporter-go/pkg/systemairmodbus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type SystemairAirflowCollector struct {
	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	fan_speed_rpm *prometheus.GaugeVec
	fan_speed_percentage *prometheus.GaugeVec
}

func NewSystemairAirflowCollector(hvac *modbus.ModbusClient, namespace string) *SystemairAirflowCollector {
	subsystem := "fan"
	return &SystemairAirflowCollector{
		hvac: hvac,
		fan_speed_rpm: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name: "speed_rpm",
				Help: "Supply/Extract Air Fan RPM indication from TACHO. Min 0 RPM, Max 5000 RPM",
		}, []string{"fan"}),
		fan_speed_percentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name: "speed_percentage",
				Help: "SAF/EAF fan speed in percentage. Min 0 %, Max 100 %",
		}, []string{"fan"}),
	}
}

func (e *SystemairAirflowCollector) Describe(ch chan<- *prometheus.Desc) {
	e.fan_speed_rpm.Describe(ch)
	e.fan_speed_percentage.Describe(ch)
}

func (e *SystemairAirflowCollector) Collect(ch chan<- prometheus.Metric) {
	for _, fan := range []string{"SAF", "EAF"} {
		e.fan_speed_rpm.WithLabelValues(fan).Set(float64(systemairmodbus.GetFanRPM(e.hvac, fan)))
		e.fan_speed_percentage.WithLabelValues(fan).Set(float64(systemairmodbus.GetFanPercentage(e.hvac, fan)))
	}
	e.fan_speed_rpm.Collect(ch)
	e.fan_speed_percentage.Collect(ch)
}