package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func setBoolGauge(gauge prometheus.Gauge, value bool) {
	if value {
		gauge.Set(1)
		return
	}
	gauge.Set(0)
}

func collectBoolGauge(ch chan<- prometheus.Metric, gauge prometheus.Gauge, value bool, err error, context string) {
	if err != nil {
		logModbusError(context, err)
	} else {
		setBoolGauge(gauge, value)
	}
	gauge.Collect(ch)
}

func collectFloatGauge(ch chan<- prometheus.Metric, gauge prometheus.Gauge, value float64, err error, context string) {
	if err != nil {
		logModbusError(context, err)
	} else {
		gauge.Set(value)
	}
	gauge.Collect(ch)
}

func collectLabelledBoolGauge(gaugeVec *prometheus.GaugeVec, label string, value bool, err error, context string) {
	if err != nil {
		logModbusError(context, err)
		return
	}
	if value {
		gaugeVec.WithLabelValues(label).Set(1)
		return
	}
	gaugeVec.WithLabelValues(label).Set(0)
}
