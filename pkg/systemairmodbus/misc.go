package systemairmodbus

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/simonvetter/modbus"
)

type iaqLevel uint16

const (
	iaqEconomic iaqLevel = iota
	iaqGood
	iaqImproving
)

type userMode uint16

const (
	userModeAuto userMode = iota
	userModeManual
	userModeCrowded
	userModeRefresh
	userModeFireplace
	userModeAway
	userModeHoliday
	userModeCookerHood
	userModeVacuumCleaner
	userModeCDI1
	userModeCDI2
	userModeCDI3
	userModePressureGuard
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
	switch iaqLevel(readRegister16(client, 1123, modbus.INPUT_REGISTER)) {
	case iaqEconomic:
		return "Economic"
	case iaqGood:
		return "Good"
	case iaqImproving:
		return "Improving"
	default:
		return unknownRegisterValue
	}
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
	switch userMode(readRegister16(client, 1161, modbus.INPUT_REGISTER)) {
	case userModeAuto:
		return "Auto"
	case userModeManual:
		return "Manual"
	case userModeCrowded:
		return "Crowded"
	case userModeRefresh:
		return "Refresh"
	case userModeFireplace:
		return "Fireplace"
	case userModeAway:
		return "Away"
	case userModeHoliday:
		return "Holiday"
	case userModeCookerHood:
		return "CookerHood"
	case userModeVacuumCleaner:
		return "VacuumCleaner"
	case userModeCDI1:
		return "CDI1"
	case userModeCDI2:
		return "CDI2"
	case userModeCDI3:
		return "CDI3"
	case userModePressureGuard:
		return "PressureGuard"
	default:
		return unknownRegisterValue
	}
}

// ActivateRefresh enables the "Refresh" user mode with the supplied duration, in minutes.
func ActivateRefresh(client *modbus.ModbusClient, duration uint16) error {
	if duration > 240 {
		return errors.New("supplied refresh mode duration is too big, max 240min")
	}
	err := writeRegister16(client, 1104, duration)
	if err != nil {
		// error out
		// TODO: handle errors more gracefully:
		// Use a provided (or default) logger,
		// Do not crash the program on failure
		fmt.Fprintf(os.Stderr, "Failed activating Refresh while setting mode duration: %v\n", err)
		os.Exit(4)
	}
	err = writeRegister16(client, 1162, 4) // modes are shifted +1 when writing
	if err != nil {
		// error out
		// TODO: handle errors more gracefully:
		// Use a provided (or default) logger,
		// Do not crash the program on failure
		fmt.Fprintf(os.Stderr, "Failed activating Refresh while setting mode: %v\n", err)
		os.Exit(4)
	}
	return nil
}

// GetUsermodeRemaining gets the "Remaining time for user mode state" as time.Duration.
func GetUsermodeRemaining(client *modbus.ModbusClient) time.Duration {
	usermodeRemaining, err := time.ParseDuration(fmt.Sprintf("%ds", readRegister32(client, 1111, modbus.INPUT_REGISTER)))
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
	filterRemaining, err := time.ParseDuration(fmt.Sprintf("%ds", readRegister32(client, 7005, modbus.INPUT_REGISTER)))
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
