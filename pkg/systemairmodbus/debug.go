package systemairmodbus

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

// PrintModbusRegisters prints all of the supported modbus registers values to standard output in a human-readable format.
func PrintModbusRegisters(client *modbus.ModbusClient) {
  // Airflow values
	fmt.Println()
	fmt.Printf("SAF: %d RPM\n", GetFanRPM(client, "SAF")) // hvac_fan_speed_rpm{fan='SAF'}
	fmt.Printf("EAF: %d RPM\n", GetFanRPM(client, "EAF"))// hvac_fan_speed_rpm{fan='EAF'}
	fmt.Printf("SAF: %d %%\n", GetFanPercentage(client, "SAF")) // hvac_fan_speed_percentage{fan='SAF'}
	fmt.Printf("EAF: %d %%\n", GetFanPercentage(client, "EAF")) // hvac_fan_speed_percentage{fan='EAF'}

	// Misc values
	fmt.Println()
	fmt.Printf("Humidity: %d %%\n", GetHumidity(client, "sensor")) // hvac_humidity_percentage{type='sensor'}
	fmt.Printf("Humidity demand: %d %%\n", GetHumidity(client, "demand")) // hvac_humidity_percentage{type='demand'}
	fmt.Printf("IAQ: %s\n", GetIAQ(client)) // hvac_iaq_level{level='<result>'}=1
	fmt.Printf("Usermode: %s\n", GetUsermode(client)) // hvac_usermode_enabled{mode='<result>'}=1
	fmt.Printf("Usermode remaining: %v\n", GetUsermodeRemaining(client)) // hvac_usermode_remaining_seconds
	fmt.Printf("Filter remaining: %v\n", GetFilterRemaining(client)) // hvac_airfilter_remaining_seconds

  // Heater values
	fmt.Println()
	fmt.Printf("Heat exchanger active: %t\n", GetHeatExchangerActive(client)) // hvac_heater_active
	fmt.Printf("Heat exchanger voltage: %.1f V\n", GetHeatExchangerVoltage(client)) // hvac_heater_voltage
	// TODO: heater & TRIAC correctness is unverified - these two might mean opposite things;
	// also, I'm not sure what TRIAC means and if it's actually an electric heater
	fmt.Printf("Electric heater active: %t\n", GetHeaterActive(client)) // hvac_heater_active
	fmt.Printf("Electric heater voltage: %.1f V\n", GetHeaterVoltage(client)) // hvac_heater_voltage
	fmt.Printf("TRIAC Electric heater active: %t\n", GetTRIACActive(client)) // hvac_triac_active
	fmt.Printf("TRIAC Electric heater voltage: %.1f V\n", GetTRIACVoltage(client)) // hvac_triac_voltage

	// Temperature values
	fmt.Println()
	fmt.Printf("ECO mode enabled: %t\n", GetEcoEnabled(client)) // hvac_eco{state='enabled'}
	fmt.Printf("ECO mode active: %t\n", GetEcoActive(client)) // hvac_eco{state='active'}
	fmt.Printf("Freecooling enabled: %t \n", GetFreecoolingEnabled(client)) // hvac_freecooling{state='enabled'}
	fmt.Printf("Freecooling active: %t \n", GetFreecoolingActive(client)) // hvac_freecooling{state='active'}
	fmt.Printf("Temp supply mode: %s\n", GetTempMode(client)) // hvac_temp_mode_enabled
	fmt.Printf("SATC controller output: %d %%\n", GetTempDemandPercentage(client)) // hvac_temp_controller_percentage
	fmt.Printf("Target room: %.1f C\n", GetTempTarget(client, "room")) // hvac_temp_target_degrees{type='room'}
	fmt.Printf("Target supply: %.1f C\n", GetTempTarget(client, "supply")) // hvac_temp_target_degrees{type='supply'}
	fmt.Printf("OAT: %.1f C\n", GetTemp(client, "OAT")) // hvac_temp_degrees{sensor='OAT'}
	fmt.Printf("SAT: %.1f C\n", GetTemp(client, "SAT")) // hvac_temp_degrees{sensor='SAT'}
	fmt.Printf("EAT: %.1f C\n", GetTemp(client, "EAT")) // hvac_temp_degrees{sensor='EAT'}
	fmt.Printf("OHT: %.1f C\n", GetTemp(client, "OHT")) // hvac_temp_degrees{sensor='OHT'}
}