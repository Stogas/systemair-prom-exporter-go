package systemairmodbus

import "github.com/simonvetter/modbus"

// GetFanRPM gets the "Supply/Extract Air Fan RPM indication from TACHO"
// Provide "SAF" or "EAF" as a value to select the wanted fan sensor.
// Min 0 RPM, Max 5000 RPM.
func GetFanRPM(client *modbus.ModbusClient, fan string) uint16 {
	switch fan {
	case "SAF":
		return readRegister16(client, 12401, modbus.HOLDING_REGISTER)
	case "EAF":
		return readRegister16(client, 12402, modbus.HOLDING_REGISTER)
	}
	return 0
}

// GetFanPercentage gets the "SAF/EAF fan speed" in percentage
// Provide "SAF" or "EAF" as a value to select the wanted fan sensor.
// Min 0 %, Max 100 %.
func GetFanPercentage(client *modbus.ModbusClient, fan string) uint16 {
	switch fan {
	case "SAF":
		return readRegister16(client, 14001, modbus.HOLDING_REGISTER)
	case "EAF":
		return readRegister16(client, 14002, modbus.HOLDING_REGISTER)
	}
	return 0
}
