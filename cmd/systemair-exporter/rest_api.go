package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"
	"github.com/simonvetter/modbus"
)

func RegisterAPI(cfg Config, m *modbus.ModbusClient) {
	http.HandleFunc("/hvac/status", getStatusHandler(m))
	fmt.Printf("Registered REST API status HTTP handler on: %v\n", "/hvac/status")

	// http.HandleFunc("/hvac/mode", setModeHandler(m))
	// fmt.Printf("Registered REST API mode set HTTP handler on: %v\n", "/hvac/mode")
}

func ReadStatus(m *modbus.ModbusClient) (s HVACStatus) {
	s.UserMode.Name = systemairmodbus.GetUsermode(m)
	s.UserMode.DurationNanoseconds = systemairmodbus.GetUsermodeRemaining(m).Nanoseconds()
	s.Temperatures.OAT = systemairmodbus.GetTemp(m, "OAT")
	s.Temperatures.SAT = systemairmodbus.GetTemp(m, "SAT")
	s.Temperatures.EAT = systemairmodbus.GetTemp(m, "EAT")
	s.Temperatures.OHT = systemairmodbus.GetTemp(m, "OHT")
	s.Temperatures.SupplyMode = systemairmodbus.GetTempMode(m)
	s.Temperatures.TargetRoom = systemairmodbus.GetTempTarget(m, "room")
	s.Temperatures.TargetSupply = systemairmodbus.GetTempTarget(m, "supply")
	s.RelativeHumidity = systemairmodbus.GetHumidity(m, "sensor")
	s.Airflow.Supply.Percentage = systemairmodbus.GetFanPercentage(m, "SAF")
	s.Airflow.Supply.RPM = systemairmodbus.GetFanRPM(m, "SAF")
	s.Airflow.Extract.Percentage = systemairmodbus.GetFanPercentage(m, "EAF")
	s.Airflow.Extract.RPM = systemairmodbus.GetFanRPM(m, "EAF")
	s.Voltages.HeatExchanger = systemairmodbus.GetHeatExchangerVoltage(m)
	s.Voltages.ElectricHeater = max(systemairmodbus.GetHeaterVoltage(m), systemairmodbus.GetTRIACVoltage(m))
	return
}

func SetMode(m *modbus.ModbusClient, mode string, ttl int64) error {
	// Implementation to set mode on Modbus
	return nil
}

func getStatusHandler(m *modbus.ModbusClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := ReadStatus(m)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func setModeHandler(m *modbus.ModbusClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request HVACUserMode
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = SetMode(m, request.Name, request.DurationNanoseconds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
