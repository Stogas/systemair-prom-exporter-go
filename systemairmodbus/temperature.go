package systemairmodbus

import "github.com/simonvetter/modbus"

// GetTempOAT gets the "Outdoor Air Temperature sensor value (standard)".
// This is the Outdoor Air Temperature in Celsius.
// Min -40 C, Max 80 C
func GetTempOAT(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 12102, modbus.HOLDING_REGISTER)) / 10
}

// GetTempSAT gets the "Supply Air Temperature sensor value (standard)".
// This is the Supply Air Temperature in Celsius.
// Min -40 C, Max 80 C
func GetTempSAT(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 12103, modbus.HOLDING_REGISTER)) / 10
}

// GetTempEAT gets the "PDM EAT sensor value (standard)".
// This is the Extract Air Temperature in Celsius.
// Min -40 C, Max 80 C
func GetTempEAT(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 12544, modbus.HOLDING_REGISTER)) / 10
}

// GetTempOHT gets the "Over Heat Temperature sensor (Electrical Heater)".
// This is the Over Heat Temperature in Celsius.
// Min -40 C, Max 80 C
func GetTempOHT(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 12108, modbus.HOLDING_REGISTER)) / 10
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

// GetTempTargetRoom gets the "Temperature setpoint for the supply air temperature".
// This is the Target Room Temperature in Celsius.
// Min 12 C, Max 30 C
func GetTempTargetRoom(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 2001, modbus.HOLDING_REGISTER)) / 10
}

// GetTempTargetSupply gets the "Temperature setpoint for the supply air temperature".
// This is the Target Supply Temperature in Celsius.
// Min 12 C, Max 30 C
func GetTempTargetSupply(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 2054, modbus.INPUT_REGISTER)) / 10
}