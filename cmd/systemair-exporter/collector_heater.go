package main

import (
	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type SystemairHeaterCollector struct {
	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	triac_active          prometheus.Gauge
	triac_voltage         prometheus.Gauge
	heater_active         prometheus.Gauge
	heater_voltage        prometheus.Gauge
	heatexchanger_active  prometheus.Gauge
	heatexchanger_voltage prometheus.Gauge

	eco         *prometheus.GaugeVec
	freecooling *prometheus.GaugeVec
}

func NewSystemairHeaterCollector(hvac *modbus.ModbusClient, namespace string) *SystemairHeaterCollector {
	return &SystemairHeaterCollector{
		hvac: hvac,
		triac_active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "triac_active",
			Help:      "Boolean gauge on whether the TRIAC Electric Heater is active. Min 0 (not active), Max 1 (active)",
		}),
		triac_voltage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "triac_voltage",
			Help:      "Voltage applied to the TRIAC Electric Heater. Min 0 V, Max 10 V",
		}),
		heater_active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heater_active",
			Help:      "Boolean gauge on whether the Electric Heater is active. Min 0 (not active), Max 1 (active)",
		}),
		heater_voltage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heater_voltage",
			Help:      "Voltage applied to the Electric Heater. Min 0 V, Max 10 V",
		}),
		heatexchanger_active: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heatexchanger_active",
			Help:      "Boolean gauge on whether the Heat Exchanger is active. Min 0 (not active), Max 1 (active)",
		}),
		heatexchanger_voltage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heatexchanger_voltage",
			Help:      "Voltage applied to the Heat Exchanger. Min 0 V, Max 10 V",
		}),
		eco: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "eco",
			Help:      "Whether ECO mode is enabled and active. Min 0 (not enabled/not active), Max 1 (enabled/active)",
		}, []string{"state"}),
		freecooling: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "freecooling",
			Help:      "Whether Freecooling mode is enabled and active. Min 0 (not enabled/not active), Max 1 (enabled/active)",
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
	triacActive, err := systemairmodbus.GetTRIACActive(e.hvac)
	collectBoolGauge(ch, e.triac_active, triacActive, err, "collector triac active")

	triacVoltage, err := systemairmodbus.GetTRIACVoltage(e.hvac)
	collectFloatGauge(ch, e.triac_voltage, triacVoltage, err, "collector triac voltage")

	heaterActive, err := systemairmodbus.GetHeaterActive(e.hvac)
	collectBoolGauge(ch, e.heater_active, heaterActive, err, "collector heater active")

	heaterVoltage, err := systemairmodbus.GetHeaterVoltage(e.hvac)
	collectFloatGauge(ch, e.heater_voltage, heaterVoltage, err, "collector heater voltage")

	heatexchangerActive, err := systemairmodbus.GetHeatExchangerActive(e.hvac)
	collectBoolGauge(ch, e.heatexchanger_active, heatexchangerActive, err, "collector heatexchanger active")

	heatexchangerVoltage, err := systemairmodbus.GetHeatExchangerVoltage(e.hvac)
	collectFloatGauge(ch, e.heatexchanger_voltage, heatexchangerVoltage, err, "collector heatexchanger voltage")

	ecoEnabled, err := systemairmodbus.GetEcoEnabled(e.hvac)
	collectLabelledBoolGauge(e.eco, "enabled", ecoEnabled, err, "collector eco enabled")

	ecoActive, err := systemairmodbus.GetEcoActive(e.hvac)
	collectLabelledBoolGauge(e.eco, "active", ecoActive, err, "collector eco active")
	e.eco.Collect(ch)

	freecoolingEnabled, err := systemairmodbus.GetFreecoolingEnabled(e.hvac)
	collectLabelledBoolGauge(e.freecooling, "enabled", freecoolingEnabled, err, "collector freecooling enabled")

	freecoolingActive, err := systemairmodbus.GetFreecoolingActive(e.hvac)
	collectLabelledBoolGauge(e.freecooling, "active", freecoolingActive, err, "collector freecooling active")
	e.freecooling.Collect(ch)
}
