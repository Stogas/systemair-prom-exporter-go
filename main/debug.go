package main

import (
	"fmt"
	"systemair-prom-exporter-go/systemairmodbus"

	"github.com/simonvetter/modbus"
)

func PrintModbusRegisters(client *modbus.ModbusClient) {
  // Airflow values
	fmt.Println()
	fmt.Printf("SAF: %d RPM\n", systemairmodbus.GetFanRPM(client, "SAF")) // hvac_fan_speed_rpm{fan='SAF'}
	fmt.Printf("EAF: %d RPM\n", systemairmodbus.GetFanRPM(client, "EAF"))// hvac_fan_speed_rpm{fan='EAF'}
	fmt.Printf("SAF: %d %%\n", systemairmodbus.GetFanPercentage(client, "SAF")) // hvac_fan_speed_percentage{fan='SAF'}
	fmt.Printf("EAF: %d %%\n", systemairmodbus.GetFanPercentage(client, "EAF")) // hvac_fan_speed_percentage{fan='EAF'}

	// Misc values
	fmt.Println()
	fmt.Printf("Humidity: %d %%\n", systemairmodbus.GetHumidity(client)) // hvac_humidity_percentage{type='sensor'}
	fmt.Printf("Humidity demand: %d %%\n", systemairmodbus.GetHumidityDemand(client)) // hvac_humidity_percentage{type='demand'}
	fmt.Printf("IAQ: %s\n", systemairmodbus.GetIAQ(client)) // hvac_iaq_level{level='<result>'}=1
	fmt.Printf("Usermode: %s\n", systemairmodbus.GetUsermode(client)) // hvac_usermode_enabled{mode='<result>'}=1
	fmt.Printf("Usermode remaining: %v\n", systemairmodbus.GetUsermodeRemaining(client))
	fmt.Printf("Filter remaining: %v\n", systemairmodbus.GetFilterRemaining(client)) // hvac_airfilter_remaining_seconds

  // Heater values
	fmt.Println()
	fmt.Printf("Heat exchanger active: %t\n", systemairmodbus.GetHeatExchangerActive(client)) // hvac_heater_active
	fmt.Printf("Heat exchanger voltage: %.1f V\n", systemairmodbus.GetHeatExchangerVoltage(client)) // hvac_heater_voltage
	// TODO: heater & TRIAC correctness is unverified - these two might mean opposite things;
	// also, I'm not sure what TRIAC means and if it's actually an electric heater
	fmt.Printf("Electric heater active: %t\n", systemairmodbus.GetHeaterActive(client)) // hvac_heater_active
	fmt.Printf("Electric heater voltage: %.1f V\n", systemairmodbus.GetHeaterVoltage(client)) // hvac_heater_voltage
	fmt.Printf("TRIAC Electric heater active: %t\n", systemairmodbus.GetTRIACActive(client)) // hvac_triac_active
	fmt.Printf("TRIAC Electric heater voltage: %.1f V\n", systemairmodbus.GetTRIACVoltage(client)) // hvac_triac_voltage

	// Temperature values
	fmt.Println()
	fmt.Printf("ECO mode enabled: %t\n", systemairmodbus.GetEcoEnabled(client)) // hvac_eco{state='enabled'}
	fmt.Printf("ECO mode active: %t\n", systemairmodbus.GetEcoActive(client)) // hvac_eco{state='active'}
	fmt.Printf("Freecooling enabled: %t \n", systemairmodbus.GetFreecoolingEnabled(client)) // hvac_freecooling{state='enabled'}
	fmt.Printf("Freecooling active: %t \n", systemairmodbus.GetFreecoolingActive(client)) // hvac_freecooling{state='active'}
	fmt.Printf("Temp supply mode: %s\n", systemairmodbus.GetTempMode(client)) // hvac_temp_mode_enabled
	fmt.Printf("SATC controller output: %d %%\n", systemairmodbus.GetTempDemandPercentage(client)) // hvac_temp_controller_percentage
	fmt.Printf("Target room: %.1f C\n", systemairmodbus.GetTempTarget(client, "room")) // hvac_temp_target_degrees{type='room'}
	fmt.Printf("Target supply: %.1f C\n", systemairmodbus.GetTempTarget(client, "supply")) // hvac_temp_target_degrees{type='supply'}
	fmt.Printf("OAT: %.1f C\n", systemairmodbus.GetTemp(client, "OAT")) // hvac_temp_degrees{sensor='OAT'}
	fmt.Printf("SAT: %.1f C\n", systemairmodbus.GetTemp(client, "SAT")) // hvac_temp_degrees{sensor='SAT'}
	fmt.Printf("EAT: %.1f C\n", systemairmodbus.GetTemp(client, "EAT")) // hvac_temp_degrees{sensor='EAT'}
	fmt.Printf("OHT: %.1f C\n", systemairmodbus.GetTemp(client, "OHT")) // hvac_temp_degrees{sensor='OHT'}
}