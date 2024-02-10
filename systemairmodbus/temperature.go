package systemairmodbus

import (
	"github.com/simonvetter/modbus"
)

// GetTemp gets the Temperature Sensor values in Celsius, based on the supplied sensor name.
// OAT is the Outdoor Air Temperature.
// SAT is the Supply Air Temperature.
// EAT is the Extract Air Temperature.
// OHT is the Over Heat Temperature.
// Min -40 C, Max 80 C
func GetTemp(client *modbus.ModbusClient, sensor string) float64 {
	switch sensor {
	case "OAT":
		return float64(readRegister16(client, 12102, modbus.HOLDING_REGISTER)) / 10
	case "SAT":
		return float64(readRegister16(client, 12103, modbus.HOLDING_REGISTER)) / 10
	case "EAT":
		return float64(readRegister16(client, 12544, modbus.HOLDING_REGISTER)) / 10
	case "OHT":
		return float64(readRegister16(client, 12108, modbus.HOLDING_REGISTER)) / 10
	}
	return -255
}

// GetTempMode gets the "Unit temperature control mode" as a string.
// The possible modes are:
// 0 - Supply
// 1 - Room
// 2 - Extract
func GetTempMode(client *modbus.ModbusClient) string {
	// TODO: use iota
	switch readRegister16(client, 2031, modbus.HOLDING_REGISTER) {
	case 0:
		return "Supply"
	case 1:
		return "Room"
	case 2:
		return "Extract"
	}
	return "Error"
}

// GetTempDemandPercentage gets the "Output of the SATC" in percentage
// Min 0 %, Max 100 %.
func GetTempDemandPercentage(client *modbus.ModbusClient) uint16 {
	return readRegister16(client, 2055, modbus.INPUT_REGISTER)
}

// GetTempTarget gets the target temperatures for the chosen target type:
// room - target for the sensor used for "Temperature Control Mode", most commonly the room/extract temperature
// supply - target supply temperature the unit wants to use to get closer to the target "room" temperature
// Note: The target supply temperature might not be achieved, depending on the unit's configuration
// Min 12 C, Max 30 C
func GetTempTarget(client *modbus.ModbusClient, target string) float64 {
	switch target {
	case "room":
		return float64(readRegister16(client, 2001, modbus.HOLDING_REGISTER)) / 10
	case "supply":
		return float64(readRegister16(client, 2054, modbus.INPUT_REGISTER)) / 10
	}
	return -255
}