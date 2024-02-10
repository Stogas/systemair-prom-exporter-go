package main

import (
	"strings"
	"systemair-prom-exporter-go/systemairmodbus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type SystemairMiscCollector struct {
	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	iaq_level *prometheus.GaugeVec
	airfilter_remaining_seconds prometheus.Gauge
	humidity_percentage *prometheus.GaugeVec
	usermode_enabled *prometheus.GaugeVec
	usermode_remaining_seconds prometheus.Gauge
}

func NewSystemairMiscCollector(hvac *modbus.ModbusClient, namespace string) *SystemairMiscCollector {
	return &SystemairMiscCollector{
		hvac: hvac,
		iaq_level: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name: "iaq_level",
				Help: "Actual IAQ level. Min 0 (not the current level), Max 1 (current level)",
		}, []string{"level"}),
		airfilter_remaining_seconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "airfilter_remaining_seconds",
			Help: "Remaining time for air filter expiry in seconds",
		}),
		humidity_percentage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "humidity_percentage",
			Help: "Current and wanted humidity values in percentage. Min 0, Max 100.",
		}, []string{"type"}),
		usermode_enabled: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "usermode_enabled",
			Help: "Currently active user mode. Min 0 (not active), Max 1 (active).",
		}, []string{"mode"}),
		usermode_remaining_seconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "usermode_remaining_seconds",
			Help: "Remaining time for the currently active user mode in seconds",
		}),
	}
}

func (e *SystemairMiscCollector) Describe(ch chan<- *prometheus.Desc) {
	e.iaq_level.Describe(ch)
	e.airfilter_remaining_seconds.Describe(ch)
	e.humidity_percentage.Describe(ch)
	e.usermode_enabled.Describe(ch)
	e.usermode_remaining_seconds.Describe(ch)
}

func (e *SystemairMiscCollector) Collect(ch chan<- prometheus.Metric) {
	for _, value := range []string{"Economic", "Good", "Improving"} {
		e.iaq_level.WithLabelValues(value).Set(0)
	}
	e.iaq_level.WithLabelValues(systemairmodbus.GetIAQ(e.hvac)).Set(1)
	e.iaq_level.Collect(ch)

	e.airfilter_remaining_seconds.Set(systemairmodbus.GetFilterRemaining(e.hvac).Seconds())
	e.airfilter_remaining_seconds.Collect(ch)

	for _, value := range []string{"sensor", "demand"} {
		e.humidity_percentage.WithLabelValues(value).Set(float64(systemairmodbus.GetHumidity(e.hvac, value)))
	}
	e.humidity_percentage.Collect(ch)

	for _, value := range []string{"auto", "manual", "crowded", "refresh", "fireplace", "away", "holiday", "cookerhood", "vacuumcleaner", "cdi1", "cdi2", "cdi3", "pressureguard"} {
		e.usermode_enabled.WithLabelValues(value).Set(0)
	}
	e.usermode_enabled.WithLabelValues(strings.ToLower(systemairmodbus.GetUsermode(e.hvac))).Set(1)
	e.usermode_enabled.Collect(ch)

	e.usermode_remaining_seconds.Set(systemairmodbus.GetUsermodeRemaining(e.hvac).Seconds())
	e.usermode_remaining_seconds.Collect(ch)
}