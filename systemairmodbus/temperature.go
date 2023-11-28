package systemairmodbus

import "github.com/simonvetter/modbus"

// GetTempEAT gets the "PDM EAT sensor value (standard)".
// This is the Extract Air Temperature in Celsius.
// Min -40 C, Max 80 C
func GetTempEAT(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 12544, modbus.HOLDING_REGISTER)) / 10
}
