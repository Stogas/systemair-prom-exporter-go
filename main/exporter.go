package main

import (
	"fmt"
	"os"
	"systemair-prom-exporter-go/systemairmodbus"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	fmt.Println("Starting Systemair-Prom-Exporter")
	var err error

	conf := &modbus.ClientConfiguration{
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    19200,
		DataBits: 8,
		Parity:   modbus.PARITY_NONE,
		StopBits: 1,
		Timeout:  2000 * time.Millisecond,
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

  // Airflow values
	fmt.Println()
	fmt.Printf("SAF: %d RPM\n", systemairmodbus.GetFanSAF_RPM(client)) // hvac_fan_speed_rpm{fan='SAF'}
	fmt.Printf("EAF: %d RPM\n", systemairmodbus.GetFanEAF_RPM(client))// hvac_fan_speed_rpm{fan='EAF'}
	fmt.Printf("SAF: %d %%\n", systemairmodbus.GetFanSAFPercentage(client)) // hvac_fan_speed_percentage{fan='SAF'}
	fmt.Printf("EAF: %d %%\n", systemairmodbus.GetFanEAFPercentage(client)) // hvac_fan_speed_percentage{fan='EAF'}

	// Misc values
	fmt.Println()
	fmt.Printf("Humidity: %d %%\n", systemairmodbus.GetHumidity(client)) // hvac_humidity_percentage{type='sensor'}
	fmt.Printf("Humidity demand: %d %%\n", systemairmodbus.GetHumidityDemand(client)) // hvac_humidity_percentage{type='demand'}
	fmt.Printf("IAQ: %s\n", systemairmodbus.GetIAQ(client)) // hvac_iaq_level{level='<result>'}=1
	fmt.Printf("Usermode: %s\n", systemairmodbus.GetUsermode(client)) // hvac_usermode_enabled{mode='<result>'}=1
	fmt.Printf("Usermode remaining: %v\n", systemairmodbus.GetUsermodeRemaining(client))
	fmt.Printf("Filter remaining: %v\n", systemairmodbus.GetFilterRemaining(client)) // hvac_airfilter_remaining_seconds

  // Heater values
	fmt.Println()
	fmt.Printf("Heat exchanger active: %t\n", systemairmodbus.GetHeatExchangerActive(client)) // hvac_heater_active
	fmt.Printf("Heat exchanger voltage: %.1f V\n", systemairmodbus.GetHeatExchangerVoltage(client)) // hvac_heater_voltage
	// TODO: heater & TRIAC correctness is unverified - these two might mean opposite things;
	// also, I'm not sure what TRIAC means and if it's actually an electric heater
	fmt.Printf("Electric heater active: %t\n", systemairmodbus.GetHeaterActive(client)) // hvac_heater_active
	fmt.Printf("Electric heater voltage: %.1f V\n", systemairmodbus.GetHeaterVoltage(client)) // hvac_heater_voltage
	fmt.Printf("TRIAC Electric heater active: %t\n", systemairmodbus.GetTRIACActive(client)) // hvac_triac_active
	fmt.Printf("TRIAC Electric heater voltage: %.1f V\n", systemairmodbus.GetTRIACVoltage(client)) // hvac_triac_voltage

	// Temperature values
	fmt.Println()
	fmt.Printf("Temp supply mode: %s\n", systemairmodbus.GetTempMode(client)) // hvac_temp_mode_enabled
	fmt.Printf("SATC controller output: %d %%\n", systemairmodbus.GetTempDemandPercentage(client)) // hvac_temp_controller_percentage
	fmt.Printf("Target room: %.1f C\n", systemairmodbus.GetTempTargetRoom(client)) // hvac_temp_target_degrees{type='room'}
	fmt.Printf("Target supply: %.1f C\n", systemairmodbus.GetTempTargetSupply(client)) // hvac_temp_target_degrees{type='supply'}
	fmt.Printf("OAT: %.1f C\n", systemairmodbus.GetTempOAT(client)) // hvac_temp_degrees{sensor='OAT'}
	fmt.Printf("SAT: %.1f C\n", systemairmodbus.GetTempSAT(client)) // hvac_temp_degrees{sensor='SAT'}
	fmt.Printf("EAT: %.1f C\n", systemairmodbus.GetTempEAT(client)) // hvac_temp_degrees{sensor='EAT'}
	fmt.Printf("OHT: %.1f C\n", systemairmodbus.GetTempOHT(client)) // hvac_temp_degrees{sensor='OHT'}
	
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

	client.SetEncoding(modbus.BIG_ENDIAN,modbus.LOW_WORD_FIRST)

	return client
}
