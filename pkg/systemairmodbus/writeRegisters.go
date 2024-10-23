package systemairmodbus

import (
	"fmt"
	"os"

	"github.com/simonvetter/modbus"
)

// readRegister16 is an internal wrapper function to modbus.ModbusClient.ReadRegister().
// handling read errors and returning a 16-bit unsigned integer.
func writeRegister16(client *modbus.ModbusClient, address uint16, value uint16) error {
	// We decrease the address by 1.
	// Unknown why, but that's the only way to get accurate values
	// according to what address the [documentation] specified
	// Assuming that this is because the [documentation]
	// starts counting from 1, while we count from 0
	//
	// [documentation]: https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20190116__REV__29_.PDF
	err := client.WriteRegister(address-1, value)

	// TODO: handle errors more gracefully:
	// Use a provided (or default) logger,
	// Do not crash the program on failure
	if err != nil {
		// error out
		fmt.Fprintf(os.Stderr, "Modbus writing register %d failed with error: %v\n", address, err)
		os.Exit(2)
	}

	return err
}
