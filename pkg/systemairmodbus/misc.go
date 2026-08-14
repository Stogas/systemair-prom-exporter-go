package systemairmodbus

import (
	"fmt"
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

const (
	regUsermodeRefreshTime = 1104
	regUsermodeCrowdedTime = 1105
	regUsermodeHMIRequest  = 1162

	maxRefreshDurationMinutes = 240
	maxCrowdedDurationHours   = 8

	// HMI write values for REG_USERMODE_HMI_CHANGE_REQUEST (1162); +1 vs read enum at 1161.
	hmiUsermodeAuto    = 1
	hmiUsermodeManual  = 2
	hmiUsermodeCrowded = 3
	hmiUsermodeRefresh = 4
)

// GetHumidity gets the "PDM RHS sensor value (standard)" and "Set point for RH demand control" as a percentage.
// To select the value, provide "sensor" or "demand" as the 'source' value
// Min 0 %, Max 100 %.
func GetHumidity(client *modbus.ModbusClient, source string) (uint16, error) {
	switch source {
	case "sensor":
		return readRegister16(client, 12136, modbus.HOLDING_REGISTER)
	case "demand":
		return readRegister16(client, 1011, modbus.HOLDING_REGISTER)
	}
	return 0, fmt.Errorf("unknown humidity source %q", source)
}

// GetIAQ gets the "Actual IAQ level" as a string.
// The possible levels are:
// 0 - Economic
// 1 - Good
// 2 - Improving
func GetIAQ(client *modbus.ModbusClient) (string, error) {
	reg, err := readRegister16(client, 1123, modbus.INPUT_REGISTER)
	if err != nil {
		return "", err
	}
	switch iaqLevel(reg) {
	case iaqEconomic:
		return "Economic", nil
	case iaqGood:
		return "Good", nil
	case iaqImproving:
		return "Improving", nil
	default:
		return unknownRegisterValue, nil
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
func GetUsermode(client *modbus.ModbusClient) (string, error) {
	reg, err := readRegister16(client, 1161, modbus.INPUT_REGISTER)
	if err != nil {
		return "", err
	}
	switch userMode(reg) {
	case userModeAuto:
		return "Auto", nil
	case userModeManual:
		return "Manual", nil
	case userModeCrowded:
		return "Crowded", nil
	case userModeRefresh:
		return "Refresh", nil
	case userModeFireplace:
		return "Fireplace", nil
	case userModeAway:
		return "Away", nil
	case userModeHoliday:
		return "Holiday", nil
	case userModeCookerHood:
		return "CookerHood", nil
	case userModeVacuumCleaner:
		return "VacuumCleaner", nil
	case userModeCDI1:
		return "CDI1", nil
	case userModeCDI2:
		return "CDI2", nil
	case userModeCDI3:
		return "CDI3", nil
	case userModePressureGuard:
		return "PressureGuard", nil
	default:
		return unknownRegisterValue, nil
	}
}

func activateTimedUserMode(
	client *modbus.ModbusClient,
	modeName string,
	durationReg uint16,
	duration uint16,
	hmiValue uint16,
) error {
	err := writeRegister16(client, durationReg, duration)
	if err != nil {
		return fmt.Errorf("activating %s while setting mode duration: %w", modeName, err)
	}
	err = writeRegister16(client, regUsermodeHMIRequest, hmiValue)
	if err != nil {
		return fmt.Errorf("activating %s while setting mode: %w", modeName, err)
	}
	return nil
}

func requestUserMode(client *modbus.ModbusClient, modeName string, hmiValue uint16) error {
	err := writeRegister16(client, regUsermodeHMIRequest, hmiValue)
	if err != nil {
		return fmt.Errorf("activating %s while setting mode: %w", modeName, err)
	}
	return nil
}

// ActivateAuto enables the Auto user mode.
func ActivateAuto(client *modbus.ModbusClient) error {
	return requestUserMode(client, "Auto", hmiUsermodeAuto)
}

// ActivateManual enables the Manual user mode.
func ActivateManual(client *modbus.ModbusClient) error {
	return requestUserMode(client, "Manual", hmiUsermodeManual)
}

// ActivateCrowded enables the Crowded user mode with the supplied duration, in hours.
func ActivateCrowded(client *modbus.ModbusClient, durationHours uint16) error {
	if durationHours < 1 || durationHours > maxCrowdedDurationHours {
		return fmt.Errorf("crowded mode duration must be between 1 and %d hours", maxCrowdedDurationHours)
	}
	return activateTimedUserMode(client, "Crowded", regUsermodeCrowdedTime, durationHours, hmiUsermodeCrowded)
}

// ActivateRefresh enables the "Refresh" user mode with the supplied duration, in minutes.
func ActivateRefresh(client *modbus.ModbusClient, durationMinutes uint16) error {
	if durationMinutes < 1 || durationMinutes > maxRefreshDurationMinutes {
		return fmt.Errorf("refresh mode duration must be between 1 and %d minutes", maxRefreshDurationMinutes)
	}
	return activateTimedUserMode(client, "Refresh", regUsermodeRefreshTime, durationMinutes, hmiUsermodeRefresh)
}

// GetUsermodeRemaining gets the "Remaining time for user mode state" as time.Duration.
func GetUsermodeRemaining(client *modbus.ModbusClient) (time.Duration, error) {
	seconds, err := readRegister32(client, 1111, modbus.INPUT_REGISTER)
	if err != nil {
		return 0, err
	}
	usermodeRemaining, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return 0, fmt.Errorf("parsing usermode remaining time: %w", err)
	}

	return usermodeRemaining, nil
}

// GetFilterRemaining gets the "Remaining time for filter" as time.Duration.
func GetFilterRemaining(client *modbus.ModbusClient) (time.Duration, error) {
	seconds, err := readRegister32(client, 7005, modbus.INPUT_REGISTER)
	if err != nil {
		return 0, err
	}
	filterRemaining, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return 0, fmt.Errorf("parsing filter remaining time: %w", err)
	}

	return filterRemaining, nil
}
