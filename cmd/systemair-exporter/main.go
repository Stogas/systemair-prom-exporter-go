package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"

	"github.com/simonvetter/modbus"
)

type Config struct {
  HTTPPort              int

	MonitoringEnabled     bool
	PercentageIncrease    float64
	RefreshDuration       uint
	AveragePeriod         int

	RestApiEnabled        bool
}

func parseFlags() Config {
	var cfg Config
	flag.BoolVar(&cfg.RestApiEnabled, "restApiEnabled", false, "Enable/disable REST API")
	flag.BoolVar(&cfg.MonitoringEnabled, "humidityMonitoring", false, "Enable/disable humidity monitoring")
	flag.Float64Var(&cfg.PercentageIncrease, "percentageIncreaseThreshold", 5.0, "Percentage increase over running average needed to trigger a spike alert")
	flag.UintVar(&cfg.RefreshDuration, "refreshDuration", 15, "Refresh mode duration for clearing out the humidity spike. Must be between 1 and 240")
	flag.IntVar(&cfg.AveragePeriod, "averagePeriod", 10, "Period over which to calculate the running average (in minutes)")
	flag.IntVar(&cfg.HTTPPort, "port", 9999, "HTTP port to listen on (for both Prometheus exporter & Rest API)")

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

	cfg := parseFlags()

	if cfg.MonitoringEnabled {
		fmt.Println("Humidity monitoring enabled. Upon a configured humidity spike, will trigger a refresh for the configured time.")
		go monitorHumidity(cfg, client)
	} else {
		fmt.Println("Humidity monitoring disabled.")
	}

	if cfg.RestApiEnabled {
		fmt.Println("REST API enabled.")
		go StartAPI(":9998", cfg, client)
	} else {
		fmt.Println("REST API disabled.")
	}

	fmt.Println("Application is running. Press Ctrl+C to stop.")

	StartExporter(":" + strconv.Itoa(cfg.HTTPPort), "/metrics", client)
}