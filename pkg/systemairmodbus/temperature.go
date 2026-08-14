package systemairmodbus

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

type tempMode uint16

const (
	tempModeSupply tempMode = iota
	tempModeRoom
	tempModeExtract
)

// GetTemp gets the Temperature Sensor values in Celsius, based on the supplied sensor name.
// OAT is the Outdoor Air Temperature.
// SAT is the Supply Air Temperature.
// EAT is the Extract Air Temperature.
// OHT is the Over Heat Temperature.
// Min -40 C, Max 80 C
func GetTemp(client *modbus.ModbusClient, sensor string) (float64, error) {
	switch sensor {
	case "OAT":
		reg, err := readRegister16Signed(client, 12102, modbus.HOLDING_REGISTER)
		if err != nil {
			return 0, err
		}
		return float64(reg) / 10, nil
	case "SAT":
		reg, err := readRegister16Signed(client, 12103, modbus.HOLDING_REGISTER)
		if err != nil {
			return 0, err
		}
		return float64(reg) / 10, nil
	case "EAT":
		reg, err := readRegister16(client, 12544, modbus.HOLDING_REGISTER)
		if err != nil {
			return 0, err
		}
		return float64(reg) / 10, nil
	case "OHT":
		reg, err := readRegister16Signed(client, 12108, modbus.HOLDING_REGISTER)
		if err != nil {
			return 0, err
		}
		return float64(reg) / 10, nil
	}
	return -255, fmt.Errorf("unknown temperature sensor %q", sensor)
}

// GetTempMode gets the "Unit temperature control mode" as a string.
// The possible modes are:
// 0 - Supply
// 1 - Room
// 2 - Extract
func GetTempMode(client *modbus.ModbusClient) (string, error) {
	reg, err := readRegister16(client, 2031, modbus.HOLDING_REGISTER)
	if err != nil {
		return "", err
	}
	switch tempMode(reg) {
	case tempModeSupply:
		return "Supply", nil
	case tempModeRoom:
		return "Room", nil
	case tempModeExtract:
		return "Extract", nil
	default:
		return unknownRegisterValue, nil
	}
}

// GetTempDemandPercentage gets the "Output of the SATC" in percentage
// Min 0 %, Max 100 %.
func GetTempDemandPercentage(client *modbus.ModbusClient) (uint16, error) {
	return readRegister16(client, 2055, modbus.INPUT_REGISTER)
}

// GetTempTarget gets the target temperatures for the chosen target type:
// room - target for the sensor used for "Temperature Control Mode", most commonly the room/extract temperature
// supply - target supply temperature the unit wants to use to get closer to the target "room" temperature
// Note: The target supply temperature might not be achieved, depending on the unit's configuration
// Min 12 C, Max 30 C
func GetTempTarget(client *modbus.ModbusClient, target string) (float64, error) {
	switch target {
	case "room":
		reg, err := readRegister16(client, 2001, modbus.HOLDING_REGISTER)
		if err != nil {
			return 0, err
		}
		return float64(reg) / 10, nil
	case "supply":
		reg, err := readRegister16(client, 2054, modbus.INPUT_REGISTER)
		if err != nil {
			return 0, err
		}
		return float64(reg) / 10, nil
	}
	return -255, fmt.Errorf("unknown temperature target %q", target)
}
