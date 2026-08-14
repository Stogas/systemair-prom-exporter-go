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

func ReadStatus(m *modbus.ModbusClient) (HVACStatus, error) {
	var s HVACStatus

	if err := readUserModeStatus(m, &s.UserMode); err != nil {
		return s, err
	}
	if err := readTemperatureStatus(m, &s.Temperatures); err != nil {
		return s, err
	}
	if err := readHumidityStatus(m, &s.RelativeHumidity); err != nil {
		return s, err
	}
	if err := readAirflowStatus(m, &s.Airflow.Supply, &s.Airflow.Extract); err != nil {
		return s, err
	}
	if err := readVoltageStatus(m, &s.Voltages); err != nil {
		return s, err
	}

	return s, nil
}

func readUserModeStatus(m *modbus.ModbusClient, userMode *HVACUserMode) error {
	name, err := systemairmodbus.GetUsermode(m)
	if err != nil {
		return err
	}
	userMode.Name = name

	remaining, err := systemairmodbus.GetUsermodeRemaining(m)
	if err != nil {
		return err
	}
	userMode.DurationNanoseconds = remaining.Nanoseconds()

	return nil
}

func readTemperatureStatus(m *modbus.ModbusClient, temps *HVACTemperature) error {
	oat, err := systemairmodbus.GetTemp(m, "OAT")
	if err != nil {
		return err
	}
	temps.OAT = oat

	sat, err := systemairmodbus.GetTemp(m, "SAT")
	if err != nil {
		return err
	}
	temps.SAT = sat

	eat, err := systemairmodbus.GetTemp(m, "EAT")
	if err != nil {
		return err
	}
	temps.EAT = eat

	oht, err := systemairmodbus.GetTemp(m, "OHT")
	if err != nil {
		return err
	}
	temps.OHT = oht

	supplyMode, err := systemairmodbus.GetTempMode(m)
	if err != nil {
		return err
	}
	temps.SupplyMode = supplyMode

	targetRoom, err := systemairmodbus.GetTempTarget(m, "room")
	if err != nil {
		return err
	}
	temps.TargetRoom = targetRoom

	targetSupply, err := systemairmodbus.GetTempTarget(m, "supply")
	if err != nil {
		return err
	}
	temps.TargetSupply = targetSupply

	return nil
}

func readHumidityStatus(m *modbus.ModbusClient, humidity *uint16) error {
	value, err := systemairmodbus.GetHumidity(m, "sensor")
	if err != nil {
		return err
	}
	*humidity = value
	return nil
}

func readAirflowStatus(m *modbus.ModbusClient, supply, extract *HVACAirflow) error {
	supplyPct, err := systemairmodbus.GetFanPercentage(m, "SAF")
	if err != nil {
		return err
	}
	supply.Percentage = supplyPct

	supplyRPM, err := systemairmodbus.GetFanRPM(m, "SAF")
	if err != nil {
		return err
	}
	supply.RPM = supplyRPM

	extractPct, err := systemairmodbus.GetFanPercentage(m, "EAF")
	if err != nil {
		return err
	}
	extract.Percentage = extractPct

	extractRPM, err := systemairmodbus.GetFanRPM(m, "EAF")
	if err != nil {
		return err
	}
	extract.RPM = extractRPM

	return nil
}

func readVoltageStatus(m *modbus.ModbusClient, voltages *HVACVoltage) error {
	heatExchanger, err := systemairmodbus.GetHeatExchangerVoltage(m)
	if err != nil {
		return err
	}
	voltages.HeatExchanger = heatExchanger

	heaterVoltage, err := systemairmodbus.GetHeaterVoltage(m)
	if err != nil {
		return err
	}
	triacVoltage, err := systemairmodbus.GetTRIACVoltage(m)
	if err != nil {
		return err
	}
	voltages.ElectricHeater = max(heaterVoltage, triacVoltage)

	return nil
}

func getStatusHandler(m *modbus.ModbusClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := ReadStatus(m)
		if err != nil {
			logModbusError("REST /hvac/status", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			fmt.Printf("Failed to encode modbus status into JSON while processing a GET REST call: %v", err)
		}
	}
}
