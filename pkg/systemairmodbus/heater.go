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
func GetHeaterVoltage(client *modbus.ModbusClient) float64 {
	return float64(readRegister16(client, 14101, modbus.INPUT_REGISTER)) / 10
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
func GetTRIACVoltage(client *modbus.ModbusClient) float64 {
	return float64(readRegister16(client, 2149, modbus.INPUT_REGISTER)) / 10
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
func GetHeatExchangerVoltage(client *modbus.ModbusClient) float64 {
	return float64(readRegister16(client, 14103, modbus.INPUT_REGISTER)) / 10
}

// GetEcoEnabled gets the "Enabling of eco mode" as a boolean.
// This shows whether the ECO mode is enabled by the user.
func GetEcoEnabled(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 2505, modbus.HOLDING_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetEcoEnabled failed to get the status of the ECO mode")
	return false // I don't want to do proper error handling in this case
}

// GetEcoActive gets the "Indication if conditions for ECO mode are active (low temperature)" as a boolean.
// This shows whether the ECO mode is active at this moment.
func GetEcoActive(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 2506, modbus.INPUT_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetEcoActive failed to get the status of the ECO mode")
	return false // I don't want to do proper error handling in this case
}

// GetFreecoolingEnabled gets the "if free cooling is enabled" as a boolean.
// This shows whether the Freecooling mode is enabled by the user.
func GetFreecoolingEnabled(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 4101, modbus.HOLDING_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetFreecoolingEnabled failed to get the status of the Freecooling mode")
	return false // I don't want to do proper error handling in this case
}

// GetFreecoolingActive gets the "if free cooling is being performed" as a boolean.
// This shows whether the Freecooling mode is active at this moment.
func GetFreecoolingActive(client *modbus.ModbusClient) bool {
	switch readRegister16(client, 4111, modbus.INPUT_REGISTER) {
	case 0:
		return false
	case 1:
		return true
	}
	fmt.Fprintf(os.Stderr, "systemairmodbus.GetFreecoolingActive failed to get the status of the Freecooling mode")
	return false // I don't want to do proper error handling in this case
}
