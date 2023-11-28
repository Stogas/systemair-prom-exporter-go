package main

import (
	"fmt"
	"os"
	"systemair-prom-exporter-go/systemairmodbus"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	fmt.Println("Starting")
	var err error

	conf := &modbus.ClientConfiguration{
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    19200,
		DataBits: 8,
		Parity:   modbus.PARITY_NONE,
		StopBits: 1,
		Timeout:  500 * time.Millisecond,
	}

	client := CreateModbusClient(conf)
	fmt.Println("Client created")

	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Modbus client Open() failed with error: %v\n", err)
		os.Exit(2)
		// note: multiple Open() attempts can be made on the same client until
		// the connection succeeds (i.e. err == nil), calling the constructor again
		// is unnecessary.
	}
	fmt.Println("Client opened")

	var eat_temp float32

	eat_temp = systemairmodbus.GetTempEAT(client) / 10
	fmt.Println("EAT temp: ", eat_temp)

	// close the TCP connection/serial port
	client.Close()
}

func CreateModbusClient(conf *modbus.ClientConfiguration) *modbus.ModbusClient {
	var client  *modbus.ModbusClient
	var err      error

	client, err = modbus.NewClient(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Modbus client creation failed with error: %v\n", err)
		os.Exit(2)
	}

	return client
}
