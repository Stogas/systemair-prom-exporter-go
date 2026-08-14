package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"
	"github.com/simonvetter/modbus"
)

type setUserModeRequest struct {
	Name     string `json:"name"`
	Duration *uint  `json:"duration"`
}

func setModeHandler(m *modbus.ModbusClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req setUserModeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		err = applyUserMode(m, req)
		if err != nil {
			if isModbusWriteError(err) {
				logModbusError("REST /hvac/mode", unwrapModbusWriteError(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func applyUserMode(m *modbus.ModbusClient, req setUserModeRequest) error {
	switch req.Name {
	case "Auto":
		if req.Duration != nil {
			return fmt.Errorf("duration is not allowed for mode %q", req.Name)
		}
		return wrapModbusWriteError(systemairmodbus.ActivateAuto(m))
	case "Manual":
		if req.Duration != nil {
			return fmt.Errorf("duration is not allowed for mode %q", req.Name)
		}
		return wrapModbusWriteError(systemairmodbus.ActivateManual(m))
	case "Refresh":
		duration, err := requiredDuration(req)
		if err != nil {
			return err
		}
		if duration > 240 {
			return fmt.Errorf("duration for mode %q must be between 1 and 240 minutes", req.Name)
		}
		return wrapModbusWriteError(systemairmodbus.ActivateRefresh(m, uint16(duration)))
	case "Crowded":
		duration, err := requiredDuration(req)
		if err != nil {
			return err
		}
		if duration > 8 {
			return fmt.Errorf("duration for mode %q must be between 1 and 8 hours", req.Name)
		}
		return wrapModbusWriteError(systemairmodbus.ActivateCrowded(m, uint16(duration)))
	default:
		return fmt.Errorf("unknown mode %q", req.Name)
	}
}

func requiredDuration(req setUserModeRequest) (uint, error) {
	if req.Duration == nil {
		return 0, fmt.Errorf("duration is required for mode %q", req.Name)
	}
	if *req.Duration == 0 {
		return 0, fmt.Errorf("duration must be greater than zero for mode %q", req.Name)
	}
	return *req.Duration, nil
}

type modbusWriteError struct {
	cause error
}

func (e modbusWriteError) Error() string {
	return e.cause.Error()
}

func wrapModbusWriteError(err error) error {
	if err == nil {
		return nil
	}
	return modbusWriteError{cause: err}
}

func isModbusWriteError(err error) bool {
	var writeErr modbusWriteError
	return errors.As(err, &writeErr)
}

func unwrapModbusWriteError(err error) error {
	var writeErr modbusWriteError
	if errors.As(err, &writeErr) {
		return writeErr.cause
	}
	return err
}
