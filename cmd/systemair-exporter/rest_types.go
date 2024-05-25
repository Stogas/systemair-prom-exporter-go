package main

type HVACStatus struct {
	UserMode         HVACUserMode    `json:"user_mode"`
	Temperatures     HVACTemperature `json:"temperatures"`
	RelativeHumidity uint16          `json:"relative_humidity"`
	Airflow          struct {
		Supply HVACAirflow `json:"supply"`
		Extract HVACAirflow `json:"extract"`
	}                                `json:"airflow"`
	Voltages         HVACVoltage     `json:"voltages"`
}

type HVACUserMode struct {
	Name                string        `json:"name"`
	DurationNanoseconds int64         `json:"duration_ns"`
}

type HVACTemperature struct {
	SupplyMode string `json:"supply_mode"`
	Target     float64 `json:"target"`
	OAT        float64 `json:"oat"`
	OHT        float64 `json:"oht"`
	SAT        float64 `json:"sat"`
	EAT        float64 `json:"eat"`
}

type HVACVoltage struct {
	HeatExchanger  float64 `json:"heat_exchanger"`
	ElectricHeater float64 `json:"electric_heater"`
}

type HVACAirflow struct {
	RPM        uint16 `json:"rpm"`
	Percentage uint16 `json:"percentage"`
}