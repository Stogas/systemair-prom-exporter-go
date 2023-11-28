package systemairmodbus

import (
	"fmt"
	"os"

	"github.com/simonvetter/modbus"
)

func GetTempEAT(client *modbus.ModbusClient) float32 {
	return float32(readRegister16(client, 12543, modbus.HOLDING_REGISTER))
}

func readRegister16(client *modbus.ModbusClient, address uint16, registerType modbus.RegType) uint16 {
	// read a single 16-bit holding register at address
	var reg16 uint16
	var err   error

	reg16, err  = client.ReadRegister(address, registerType)
	if err != nil {
		// error out
		fmt.Fprintf(os.Stderr, "Modbus reading register %d failed with error: %v\n", address, err)
		os.Exit(2)
	}

	return reg16
}
