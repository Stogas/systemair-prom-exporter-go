package main

import (
	"systemair-prom-exporter-go/pkg/systemairmodbus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type SystemairHeaterCollector struct {
	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	triac_active prometheus.Gauge
	triac_voltage prometheus.Gauge
	heater_active prometheus.Gauge
	heater_voltage prometheus.Gauge
	heatexchanger_active prometheus.Gauge
	heatexchanger_voltage prometheus.Gauge

	eco *prometheus.GaugeVec
	freecooling *prometheus.GaugeVec
}

func NewSystemairHeaterCollector(hvac *modbus.ModbusClient, namespace string) *SystemairHeaterCollector {
	return &SystemairHeaterCollector{
		hvac: hvac,
		triac_active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "triac_active",
			Help: "Boolean gauge on whether the TRIAC Electric Heater is active. Min 0 (not active), Max 1 (active)",
		}),
		triac_voltage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "triac_voltage",
			Help: "Voltage applied to the TRIAC Electric Heater. Min 0 V, Max 10 V",
		}),
		heater_active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "heater_active",
			Help: "Boolean gauge on whether the Electric Heater is active. Min 0 (not active), Max 1 (active)",
		}),
		heater_voltage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "heater_voltage",
			Help: "Voltage applied to the Electric Heater. Min 0 V, Max 10 V",
		}),
		heatexchanger_active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "heatexchanger_active",
			Help: "Boolean gauge on whether the Heat Exchanger is active. Min 0 (not active), Max 1 (active)",
		}),
		heatexchanger_voltage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "heatexchanger_voltage",
			Help: "Voltage applied to the Heat Exchanger. Min 0 V, Max 10 V",
		}),
		eco: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "eco",
			Help: "Whether ECO mode is enabled and active. Min 0 (not enabled/not active), Max 1 (enabled/active)",
		}, []string{"state"}),
		freecooling: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name: "freecooling",
			Help: "Whether Freecooling mode is enabled and active. Min 0 (not enabled/not active), Max 1 (enabled/active)",
		}, []string{"state"}),
	}
}

func (e *SystemairHeaterCollector) Describe(ch chan<- *prometheus.Desc) {
	e.triac_active.Describe(ch)
	e.triac_voltage.Describe(ch)
	e.heater_active.Describe(ch)
	e.heater_voltage.Describe(ch)
	e.heatexchanger_active.Describe(ch)
	e.heatexchanger_voltage.Describe(ch)

	e.eco.Describe(ch)
	e.freecooling.Describe(ch)
}

func (e *SystemairHeaterCollector) Collect(ch chan<- prometheus.Metric) {
	if systemairmodbus.GetTRIACActive(e.hvac) {
		e.triac_active.Set(1)
	} else {
		e.triac_active.Set(0)
	}
	e.triac_active.Collect(ch)

	e.triac_voltage.Set(systemairmodbus.GetTRIACVoltage(e.hvac))
	e.triac_voltage.Collect(ch)

	if systemairmodbus.GetHeaterActive(e.hvac) {
		e.heater_active.Set(1)
	} else {
		e.heater_active.Set(0)
	}
	e.heater_active.Collect(ch)

	e.heater_voltage.Set(systemairmodbus.GetHeaterVoltage(e.hvac))
	e.heater_voltage.Collect(ch)

	if systemairmodbus.GetHeatExchangerActive(e.hvac) {
		e.heatexchanger_active.Set(1)
	} else {
		e.heatexchanger_active.Set(0)
	}
	e.heatexchanger_active.Collect(ch)

	e.heatexchanger_voltage.Set(systemairmodbus.GetHeatExchangerVoltage(e.hvac))
	e.heatexchanger_voltage.Collect(ch)

	if systemairmodbus.GetEcoEnabled(e.hvac) {
		e.eco.WithLabelValues("enabled").Set(1)
	} else {
		e.eco.WithLabelValues("enabled").Set(0)
	}
	if systemairmodbus.GetEcoActive(e.hvac) {
		e.eco.WithLabelValues("active").Set(1)
	} else {
		e.eco.WithLabelValues("active").Set(0)
	}
	e.eco.Collect(ch)

	if systemairmodbus.GetFreecoolingEnabled(e.hvac) {
		e.freecooling.WithLabelValues("enabled").Set(1)
	} else {
		e.freecooling.WithLabelValues("enabled").Set(0)
	}
	if systemairmodbus.GetFreecoolingActive(e.hvac) {
		e.freecooling.WithLabelValues("active").Set(1)
	} else {
		e.freecooling.WithLabelValues("active").Set(0)
	}
	e.freecooling.Collect(ch)
}