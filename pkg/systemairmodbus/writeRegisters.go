package systemairmodbus

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

// writeRegister16 is an internal wrapper function to modbus.ModbusClient.WriteRegister().
func writeRegister16(client *modbus.ModbusClient, address uint16, value uint16) error {
	// We decrease the address by 1.
	// Unknown why, but that's the only way to get accurate values
	// according to what address the [documentation] specified
	// Assuming that this is because the [documentation]
	// starts counting from 1, while we count from 0
	//
	// [documentation]: https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20190116__REV__29_.PDF
	err := client.WriteRegister(address-1, value)
	if err != nil {
		return fmt.Errorf("modbus write register %d: %w", address, err)
	}

	return nil
}
