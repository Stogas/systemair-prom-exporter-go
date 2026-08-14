package main

import (
	"strings"

	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type SystemairTempCollector struct {
	// ModbusClient which we will target for systemair-prom-exporter-go/systemairmodbus functions
	hvac *modbus.ModbusClient

	temp_mode_enabled          *prometheus.GaugeVec
	temp_degrees               *prometheus.GaugeVec
	temp_target_degrees        *prometheus.GaugeVec
	temp_controller_percentage prometheus.Gauge
}

func NewSystemairTempCollector(hvac *modbus.ModbusClient, namespace string) *SystemairTempCollector {
	subsystem := "temp"
	return &SystemairTempCollector{
		hvac: hvac,
		temp_mode_enabled: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "mode_enabled",
			Help:      "Unit temperature control mode. The currently enabled mode has a value of 1",
		}, []string{"mode"}),
		temp_degrees: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "degrees",
			Help:      "Air Temperature sensor values for: OAT (Outside Air Temp), EAT (Extract Air Temp), SAT (Supply Air Temp), OHT (Over Heat Temp). Values are in celsius, Min -40 C, Max 80 C",
		}, []string{"sensor"}),
		temp_target_degrees: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "target_degrees",
			Help:      "Target Celsius temperatures for the \"room\" (selected temperature control mode) and the \"supply\" (to achieve the \"room\" temperature). Min 12 C, Max 30 C",
		}, []string{"type"}),
		temp_controller_percentage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "controller_percentage",
			Help:      "\"Output of the SATC\" in percentage. Min 0, Max 100",
		}),
	}
}

func (e *SystemairTempCollector) Describe(ch chan<- *prometheus.Desc) {
	e.temp_mode_enabled.Describe(ch)
	e.temp_degrees.Describe(ch)
	e.temp_target_degrees.Describe(ch)
	e.temp_controller_percentage.Describe(ch)
}

func (e *SystemairTempCollector) Collect(ch chan<- prometheus.Metric) {
	for _, mode := range []string{"supply", "room", "extract"} {
		e.temp_mode_enabled.WithLabelValues(mode).Set(0)
	}
	tempMode, err := systemairmodbus.GetTempMode(e.hvac)
	if err != nil {
		logModbusError("collector temp mode", err)
	} else {
		e.temp_mode_enabled.WithLabelValues(strings.ToLower(tempMode)).Set(1)
	}
	e.temp_mode_enabled.Collect(ch)

	for _, sensor := range []string{"OAT", "SAT", "EAT", "OHT"} {
		temp, err := systemairmodbus.GetTemp(e.hvac, sensor)
		if err != nil {
			logModbusError("collector temp degrees "+sensor, err)
		} else {
			e.temp_degrees.WithLabelValues(sensor).Set(temp)
		}
	}
	e.temp_degrees.Collect(ch)

	for _, target := range []string{"room", "supply"} {
		temp, err := systemairmodbus.GetTempTarget(e.hvac, target)
		if err != nil {
			logModbusError("collector temp target "+target, err)
		} else {
			e.temp_target_degrees.WithLabelValues(target).Set(temp)
		}
	}
	e.temp_target_degrees.Collect(ch)

	demandPct, err := systemairmodbus.GetTempDemandPercentage(e.hvac)
	if err != nil {
		logModbusError("collector temp controller percentage", err)
	} else {
		e.temp_controller_percentage.Set(float64(demandPct))
	}
	e.temp_controller_percentage.Collect(ch)
}
