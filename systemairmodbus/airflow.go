package systemairmodbus

import "github.com/simonvetter/modbus"

// GetFanSAF_RPM gets the "Supply Air Fan RPM indication from TACHO"
// Min 0 RPM, Max 5000 RPM.
func GetFanSAF_RPM(client *modbus.ModbusClient) uint16 {
	return readRegister16(client, 12401, modbus.HOLDING_REGISTER)
}

// GetFanEAF_RPM gets the "Extract Air Fan RPM indication from TACHO"
// Min 0 RPM, Max 5000 RPM.
func GetFanEAF_RPM(client *modbus.ModbusClient) uint16 {
	return readRegister16(client, 12402, modbus.HOLDING_REGISTER)
}

// GetFanSAFPercentage gets the "SAF fan speed" in percentage
// Min 0 %, Max 100 %.
func GetFanSAFPercentage(client *modbus.ModbusClient) uint16 {
	return readRegister16(client, 14001, modbus.HOLDING_REGISTER)
}

// GetFanEAFPercentage gets the "EAF fan speed" in percentage
// Min 0 %, Max 100 %.
func GetFanEAFPercentage(client *modbus.ModbusClient) uint16 {
	return readRegister16(client, 14002, modbus.HOLDING_REGISTER)
}