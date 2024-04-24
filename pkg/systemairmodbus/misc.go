package systemairmodbus

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/simonvetter/modbus"
)

// GetHumidity gets the "PDM RHS sensor value (standard)" and "Set point for RH demand control" as a percentage.
// To select the value, provide "sensor" or "demand" as the 'source' value
// Min 0 %, Max 100 %.
func GetHumidity(client *modbus.ModbusClient, source string) uint16 {
	switch source {
	case "sensor":
		return readRegister16(client, 12136, modbus.HOLDING_REGISTER)
	case "demand":
		return readRegister16(client, 1011, modbus.HOLDING_REGISTER)
	}
	return 0
}

// GetIAQ gets the "Actual IAQ level" as a string.
// The possible levels are:
// 0 - Economic
// 1 - Good
// 2 - Improving
func GetIAQ(client *modbus.ModbusClient) string {
	// TODO: use iota
	switch readRegister16(client, 1123, modbus.INPUT_REGISTER) {
	case 0:
		return "Economic"
	case 1:
		return "Good"
	case 2:
		return "Improving"
	}
	return "Error"
}

// GetUsermode gets the "Active user mode" as a string.
// The possible modes are:
// 0 - Auto
// 1 - Manual
// 2 - Crowded
// 3 - Refresh
// 4 - Fireplace
// 5 - Away
// 6 - Holiday
// 7 - CookerHood
// 8 - VacuumCleaner
// 9 - CDI1
// 10 - CDI2
// 11 - CDI3
// 12 - PressureGuard
func GetUsermode(client *modbus.ModbusClient) string {
	// TODO: use iota
	switch readRegister16(client, 1161, modbus.INPUT_REGISTER) {
	case 0:
		return "Auto"
	case 1:
		return "Manual"
	case 2:
		return "Crowded"
	case 3:
		return "Refresh"
	case 4:
		return "Fireplace"
	case 5:
		return "Away"
	case 6:
		return "Holiday"
	case 7:
		return "CookerHood"
	case 8:
		return "VacuumCleaner"
	case 9:
		return "CDI1"
	case 10:
		return "CDI2"
	case 11:
		return "CDI3"
	case 12:
		return "PressureGuard"
	}
	return "Error"
}

// ActivateRefresh enables the "Refresh" user mode with the supplied duration, in minutes.
func ActivateRefresh(client *modbus.ModbusClient, duration uint16) error {
	if duration > 240 {
		return errors.New("supplied refresh mode duration is too big, max 240min")
	}
	writeRegister16(client, 1104, duration)
	writeRegister16(client, 1162, 4) // modes are shifted +1 when writing
	return nil
}

// GetUsermodeRemaining gets the "Remaining time for user mode state" as time.Duration.
func GetUsermodeRemaining(client *modbus.ModbusClient) time.Duration {
	var err error
	var usermodeRemaining time.Duration

	usermodeRemaining, err = time.ParseDuration(fmt.Sprintf("%ds",readRegister32(client, 1111, modbus.INPUT_REGISTER)))
	if err != nil {
		// error out
		// TODO: handle errors more gracefully:
		// Use a provided (or default) logger,
		// Do not crash the program on failure
		fmt.Fprintf(os.Stderr, "Parsing time failed with error: %v\n", err)
		os.Exit(4)
	}

	return usermodeRemaining
}

// GetFilterRemaining gets the "Remaining time for filter" as time.Duration.
func GetFilterRemaining(client *modbus.ModbusClient) time.Duration {
	var err error
	var filterRemaining time.Duration

	filterRemaining, err = time.ParseDuration(fmt.Sprintf("%ds",readRegister32(client, 7005, modbus.INPUT_REGISTER)))
	if err != nil {
		// error out
		// TODO: handle errors more gracefully:
		// Use a provided (or default) logger,
		// Do not crash the program on failure
		fmt.Fprintf(os.Stderr, "Parsing time failed with error: %v\n", err)
		os.Exit(4)
	}
	
	return filterRemaining
}

// Not implemented:
// SetRefresh()
// Write address 1104 (set refresh time)
// Write address 1162 (set new desired user mode)