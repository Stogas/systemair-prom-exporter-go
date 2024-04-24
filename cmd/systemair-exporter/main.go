package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"

	"github.com/simonvetter/modbus"
)

type HumidityConfig struct {
	MonitoringEnabled     bool
	PercentageIncrease    float64
	RefreshDuration       uint
	AveragePeriod         int
}

func parseFlags() HumidityConfig {
	var cfg HumidityConfig
	flag.BoolVar(&cfg.MonitoringEnabled, "humidityMonitoring", false, "Enable/disable humidity monitoring")
	flag.Float64Var(&cfg.PercentageIncrease, "percentageIncreaseThreshold", 5.0, "Percentage increase over running average needed to trigger a spike alert")
	flag.UintVar(&cfg.RefreshDuration, "refreshDuration", 15, "Refresh mode duration for clearing out the humidity spike. Must be between 1 and 240")
	flag.IntVar(&cfg.AveragePeriod, "averagePeriod", 10, "Period over which to calculate the running average (in minutes)")
	flag.Parse()
	return cfg
}

func main() {
	fmt.Println("Starting Systemair-Prom-Exporter")

	// Create a modbus client on which data will be fetched
	conf := &modbus.ClientConfiguration{
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    19200,
		DataBits: 8,
		Parity:   modbus.PARITY_NONE,
		StopBits: 1,
		Timeout:  2000 * time.Millisecond,
	}
	client := systemairmodbus.CreateAndOpenModbusClient(conf)
	defer client.Close()

	systemairmodbus.PrintModbusRegisters(client)

	cfgHumidity := parseFlags()
	if cfgHumidity.MonitoringEnabled {
		fmt.Println("Humidity monitoring enabled. Upon a configured humidity spike, will trigger a refresh for the configured time.")
		go monitorHumidity(cfgHumidity, client)
	} else {
		fmt.Println("Humidity monitoring disabled.")
	}

	fmt.Println("Application is running. Press Ctrl+C to stop.")

	StartExporter(":9999", "/metrics", client)
}