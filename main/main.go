package main

import (
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

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
	client := CreateAndOpenModbusClient(conf)
	defer client.Close()

	PrintModbusRegisters(client)

	StartExporter(":2112", "/metrics", client)
}