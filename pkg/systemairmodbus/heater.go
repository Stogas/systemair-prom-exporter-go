package systemairmodbus

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

func registerBool(value uint16) (bool, error) {
	switch value {
	case 0:
		return false, nil
	case 1:
		return true, nil
	}
	return false, fmt.Errorf("unexpected register value %d", value)
}

// GetHeaterActive gets the "Heater DO state" as a boolean.
// This shows whether the Electric Heater is active.
func GetHeaterActive(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 14102, modbus.INPUT_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}

// GetHeaterVoltage gets the "Heater AO state".
// This is the Voltage applied to the Electric Heater.
// Min 0 V, Max 10 V
func GetHeaterVoltage(client *modbus.ModbusClient) (float64, error) {
	reg, err := readRegister16(client, 14101, modbus.INPUT_REGISTER)
	if err != nil {
		return 0, err
	}
	return float64(reg) / 10, nil
}

// GetTRIACActive gets the "TRIAC control signal" as a boolean.
// This shows whether the TRIAC Electric Heater is active.
func GetTRIACActive(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 14381, modbus.INPUT_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}

// GetTRIACVoltage gets the "TRIAC after manual override".
// This is the Voltage applied to the TRIAC Electric Heater.
// Min 0 V, Max 10 V
func GetTRIACVoltage(client *modbus.ModbusClient) (float64, error) {
	reg, err := readRegister16(client, 2149, modbus.INPUT_REGISTER)
	if err != nil {
		return 0, err
	}
	return float64(reg) / 10, nil
}

// GetHeatExchangerActive gets the "TRIAC control signal" as a boolean.
// This shows whether the Heat Exchanger is active.
func GetHeatExchangerActive(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 14104, modbus.INPUT_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}

// GetHeatExchangerVoltage gets the "Heat Exchanger AO state".
// This is the Voltage applied to the Heat Exchanger.
// Min 0 V, Max 10 V
func GetHeatExchangerVoltage(client *modbus.ModbusClient) (float64, error) {
	reg, err := readRegister16(client, 14103, modbus.INPUT_REGISTER)
	if err != nil {
		return 0, err
	}
	return float64(reg) / 10, nil
}

// GetEcoEnabled gets the "Enabling of eco mode" as a boolean.
// This shows whether the ECO mode is enabled by the user.
func GetEcoEnabled(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 2505, modbus.HOLDING_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}

// GetEcoActive gets the "Indication if conditions for ECO mode are active (low temperature)" as a boolean.
// This shows whether the ECO mode is active at this moment.
func GetEcoActive(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 2506, modbus.INPUT_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}

// GetFreecoolingEnabled gets the "if free cooling is enabled" as a boolean.
// This shows whether the Freecooling mode is enabled by the user.
func GetFreecoolingEnabled(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 4101, modbus.HOLDING_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}

// GetFreecoolingActive gets the "if free cooling is being performed" as a boolean.
// This shows whether the Freecooling mode is active at this moment.
func GetFreecoolingActive(client *modbus.ModbusClient) (bool, error) {
	reg, err := readRegister16(client, 4111, modbus.INPUT_REGISTER)
	if err != nil {
		return false, err
	}
	return registerBool(reg)
}
