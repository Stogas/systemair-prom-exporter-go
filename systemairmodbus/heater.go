package systemairmodbus

import (
	"fmt"
	"os"

	"github.com/simonvetter/modbus"
)

// GetHeaterActive gets the "Heater DO state" as a boolean.
// This shows whether the Electric Heater is active.
func GetHeaterActive(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 14102, modbus.INPUT_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetHeaterEactive failed to get the status of the heater")
	return false // I don't want to do proper error handling in this case
}

// GetHeaterVoltage gets the "Heater AO state".
// This is the Voltage applied to the Electric Heater.
// Min 0 V, Max 10 V
func GetHeaterVoltage(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 14101, modbus.INPUT_REGISTER)) / 10
}

// GetTRIACActive gets the "TRIAC control signal" as a boolean.
// This shows whether the TRIAC Electric Heater is active.
func GetTRIACActive(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 14381, modbus.INPUT_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetTRIACActive failed to get the status of the heater")
	return false // I don't want to do proper error handling in this case
}

// GetTRIACVoltage gets the "TRIAC after manual override".
// This is the Voltage applied to the TRIAC Electric Heater.
// Min 0 V, Max 10 V
func GetTRIACVoltage(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 2149, modbus.INPUT_REGISTER)) / 10
}

// GetHeatExchangerActive gets the "TRIAC control signal" as a boolean.
// This shows whether the Heat Exchanger is active.
func GetHeatExchangerActive(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 14104, modbus.INPUT_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetTRIACActive failed to get the status of the heater")
	return false // I don't want to do proper error handling in this case
}

// GetHeatExchangerVoltage gets the "Heat Exchanger AO state".
// This is the Voltage applied to the Heat Exchanger.
// Min 0 V, Max 10 V
func GetHeatExchangerVoltage(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 14103, modbus.INPUT_REGISTER)) / 10
}