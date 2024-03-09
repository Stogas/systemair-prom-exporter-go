// Package systemairmodbus provides various named functions
// for convenient querying of Systemair HVAC modbus registers.
//
// The provided functions wrap modbus.ModbusClient.ReadRegister()
// and other similar functions from the [modbus package].
//
// This allows querying by register name,
// instead of supplying a register address value,
// thus increasing readability.
//
// For example, one can use GetTempEAT(client)
// instead of client.ReadRegister(a int, t modbus.RegType).
//
// [modbus package]: https://github.com/simonvetter/modbus
package systemairmodbus

import (
	"fmt"
	"os"

	"github.com/simonvetter/modbus"
)

// readRegister16 is an internal wrapper function to modbus.ModbusClient.ReadRegister().
// handling read errors and returning a 16-bit unsigned integer.
func readRegister16(client *modbus.ModbusClient, address uint16, registerType modbus.RegType) uint16 {
	var reg16 uint16
	var err   error

	// We decrease the address by 1.
	// Unknown why, but that's the only way to get accurate values
	// according to what address the [documentation] specified
	// Assuming that this is because the [documentation]
	// starts counting from 1, while we count from 0
	//
	// [documentation]: https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20190116__REV__29_.PDF
	reg16, err  = client.ReadRegister(address - 1, registerType)

	// TODO: handle errors more gracefully:
	// Use a provided (or default) logger,
	// Do not crash the program on failure
	if err != nil {
		// error out
		fmt.Fprintf(os.Stderr, "Modbus reading register %d failed with error: %v\n", address, err)
		os.Exit(2)
	}

	return reg16
}

// readRegister16Signed is an internal wrapper function to modbus.ModbusClient.ReadRegister().
// handling read errors and returning a 16-bit signed integer.
func readRegister16Signed(client *modbus.ModbusClient, address uint16, registerType modbus.RegType) int16 {
	var reg16 uint16
	var err   error

	// We decrease the address by 1.
	// Unknown why, but that's the only way to get accurate values
	// according to what address the [documentation] specified
	// Assuming that this is because the [documentation]
	// starts counting from 1, while we count from 0
	//
	// [documentation]: https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20190116__REV__29_.PDF
	reg16, err  = client.ReadRegister(address - 1, registerType)

	// TODO: handle errors more gracefully:
	// Use a provided (or default) logger,
	// Do not crash the program on failure
	if err != nil {
		// error out
		fmt.Fprintf(os.Stderr, "Modbus reading register %d failed with error: %v\n", address, err)
		os.Exit(2)
	}

	return int16(reg16)
}

// readRegister32 is an internal wrapper function to modbus.ModbusClient.ReadUint32().
// handling read errors and returning a 32-bit unsigned integer.
func readRegister32(client *modbus.ModbusClient, address uint16, registerType modbus.RegType) uint32 {
	var reg32 uint32
	var err    error

	// We decrease the address by 1.
	// Unknown why, but that's the only way to get accurate values
	// according to what address the [documentation] specified
	// Assuming that this is because the [documentation]
	// starts counting from 1, while we count from 0
	//
	// [documentation]: https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20190116__REV__29_.PDF
	reg32, err  = client.ReadUint32(address - 1, registerType)

	// TODO: handle errors more gracefully:
	// Use a provided (or default) logger,
	// Do not crash the program on failure
	if err != nil {
		// error out
		fmt.Fprintf(os.Stderr, "Modbus reading register %d failed with error: %v\n", address, err)
		os.Exit(2)
	}

	return reg32
}
