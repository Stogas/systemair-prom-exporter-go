package systemairmodbus

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

// PrintModbusRegisters prints all of the supported modbus registers values to standard output in a human-readable format.
func PrintModbusRegisters(client *modbus.ModbusClient) {
	// Airflow values
	fmt.Println()
	printFanRPM(client, "SAF")
	printFanRPM(client, "EAF")
	printFanPercentage(client, "SAF")
	printFanPercentage(client, "EAF")

	// Misc values
	fmt.Println()
	printHumidity(client, "sensor")
	printHumidity(client, "demand")
	printIAQ(client)
	printUsermode(client)
	printUsermodeRemaining(client)
	printFilterRemaining(client)

	// Heater values
	fmt.Println()
	printHeatExchangerActive(client)
	printHeatExchangerVoltage(client)
	// Note: heater & TRIAC correctness is unverified - these two might mean opposite things;
	// also, I'm not sure what TRIAC means and if it's actually an electric heater
	printHeaterActive(client)
	printHeaterVoltage(client)
	printTRIACActive(client)
	printTRIACVoltage(client)

	// Temperature values
	fmt.Println()
	printEcoEnabled(client)
	printEcoActive(client)
	printFreecoolingEnabled(client)
	printFreecoolingActive(client)
	printTempMode(client)
	printTempDemandPercentage(client)
	printTempTarget(client, "room")
	printTempTarget(client, "supply")
	printTemp(client, "OAT")
	printTemp(client, "SAT")
	printTemp(client, "EAT")
	printTemp(client, "OHT")
}

func printFanRPM(client *modbus.ModbusClient, fan string) {
	rpm, err := GetFanRPM(client, fan)
	if err != nil {
		fmt.Printf("%s RPM: error: %v\n", fan, err)
		return
	}
	fmt.Printf("%s: %d RPM\n", fan, rpm)
}

func printFanPercentage(client *modbus.ModbusClient, fan string) {
	pct, err := GetFanPercentage(client, fan)
	if err != nil {
		fmt.Printf("%s fan speed: error: %v\n", fan, err)
		return
	}
	fmt.Printf("%s: %d %%\n", fan, pct)
}

func printHumidity(client *modbus.ModbusClient, source string) {
	humidity, err := GetHumidity(client, source)
	if err != nil {
		fmt.Printf("Humidity %s: error: %v\n", source, err)
		return
	}
	if source == "sensor" {
		fmt.Printf("Humidity: %d %%\n", humidity)
	} else {
		fmt.Printf("Humidity demand: %d %%\n", humidity)
	}
}

func printIAQ(client *modbus.ModbusClient) {
	iaq, err := GetIAQ(client)
	if err != nil {
		fmt.Printf("IAQ: error: %v\n", err)
		return
	}
	fmt.Printf("IAQ: %s\n", iaq)
}

func printUsermode(client *modbus.ModbusClient) {
	mode, err := GetUsermode(client)
	if err != nil {
		fmt.Printf("Usermode: error: %v\n", err)
		return
	}
	fmt.Printf("Usermode: %s\n", mode)
}

func printUsermodeRemaining(client *modbus.ModbusClient) {
	remaining, err := GetUsermodeRemaining(client)
	if err != nil {
		fmt.Printf("Usermode remaining: error: %v\n", err)
		return
	}
	fmt.Printf("Usermode remaining: %v\n", remaining)
}

func printFilterRemaining(client *modbus.ModbusClient) {
	remaining, err := GetFilterRemaining(client)
	if err != nil {
		fmt.Printf("Filter remaining: error: %v\n", err)
		return
	}
	fmt.Printf("Filter remaining: %v\n", remaining)
}

func printHeatExchangerActive(client *modbus.ModbusClient) {
	active, err := GetHeatExchangerActive(client)
	if err != nil {
		fmt.Printf("Heat exchanger active: error: %v\n", err)
		return
	}
	fmt.Printf("Heat exchanger active: %t\n", active)
}

func printHeatExchangerVoltage(client *modbus.ModbusClient) {
	voltage, err := GetHeatExchangerVoltage(client)
	if err != nil {
		fmt.Printf("Heat exchanger voltage: error: %v\n", err)
		return
	}
	fmt.Printf("Heat exchanger voltage: %.1f V\n", voltage)
}

func printHeaterActive(client *modbus.ModbusClient) {
	active, err := GetHeaterActive(client)
	if err != nil {
		fmt.Printf("Electric heater active: error: %v\n", err)
		return
	}
	fmt.Printf("Electric heater active: %t\n", active)
}

func printHeaterVoltage(client *modbus.ModbusClient) {
	voltage, err := GetHeaterVoltage(client)
	if err != nil {
		fmt.Printf("Electric heater voltage: error: %v\n", err)
		return
	}
	fmt.Printf("Electric heater voltage: %.1f V\n", voltage)
}

func printTRIACActive(client *modbus.ModbusClient) {
	active, err := GetTRIACActive(client)
	if err != nil {
		fmt.Printf("TRIAC Electric heater active: error: %v\n", err)
		return
	}
	fmt.Printf("TRIAC Electric heater active: %t\n", active)
}

func printTRIACVoltage(client *modbus.ModbusClient) {
	voltage, err := GetTRIACVoltage(client)
	if err != nil {
		fmt.Printf("TRIAC Electric heater voltage: error: %v\n", err)
		return
	}
	fmt.Printf("TRIAC Electric heater voltage: %.1f V\n", voltage)
}

func printEcoEnabled(client *modbus.ModbusClient) {
	enabled, err := GetEcoEnabled(client)
	if err != nil {
		fmt.Printf("ECO mode enabled: error: %v\n", err)
		return
	}
	fmt.Printf("ECO mode enabled: %t\n", enabled)
}

func printEcoActive(client *modbus.ModbusClient) {
	active, err := GetEcoActive(client)
	if err != nil {
		fmt.Printf("ECO mode active: error: %v\n", err)
		return
	}
	fmt.Printf("ECO mode active: %t\n", active)
}

func printFreecoolingEnabled(client *modbus.ModbusClient) {
	enabled, err := GetFreecoolingEnabled(client)
	if err != nil {
		fmt.Printf("Freecooling enabled: error: %v\n", err)
		return
	}
	fmt.Printf("Freecooling enabled: %t \n", enabled)
}

func printFreecoolingActive(client *modbus.ModbusClient) {
	active, err := GetFreecoolingActive(client)
	if err != nil {
		fmt.Printf("Freecooling active: error: %v\n", err)
		return
	}
	fmt.Printf("Freecooling active: %t \n", active)
}

func printTempMode(client *modbus.ModbusClient) {
	mode, err := GetTempMode(client)
	if err != nil {
		fmt.Printf("Temp supply mode: error: %v\n", err)
		return
	}
	fmt.Printf("Temp supply mode: %s\n", mode)
}

func printTempDemandPercentage(client *modbus.ModbusClient) {
	pct, err := GetTempDemandPercentage(client)
	if err != nil {
		fmt.Printf("SATC controller output: error: %v\n", err)
		return
	}
	fmt.Printf("SATC controller output: %d %%\n", pct)
}

func printTempTarget(client *modbus.ModbusClient, target string) {
	temp, err := GetTempTarget(client, target)
	if err != nil {
		fmt.Printf("Target %s: error: %v\n", target, err)
		return
	}
	if target == "room" {
		fmt.Printf("Target room: %.1f C\n", temp)
	} else {
		fmt.Printf("Target supply: %.1f C\n", temp)
	}
}

func printTemp(client *modbus.ModbusClient, sensor string) {
	temp, err := GetTemp(client, sensor)
	if err != nil {
		fmt.Printf("%s: error: %v\n", sensor, err)
		return
	}
	fmt.Printf("%s: %.1f C\n", sensor, temp)
}
