package main

import (
	"fmt"
	"math"
	"time"

	"github.com/Stogas/systemair-prom-exporter-go/pkg/systemairmodbus"
	"github.com/simonvetter/modbus"
)

func monitorHumidity(cfg Config, m *modbus.ModbusClient) {
	fmt.Println("Humidity monitoring started.")

	humidityData := make([]uint16, 0, cfg.AveragePeriod)
	currentHumidity := systemairmodbus.GetHumidity(m, "sensor")
	humidityData = append(humidityData, currentHumidity)
	fmt.Printf("Current Humidity: %d%%\n", currentHumidity)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		currentHumidity := systemairmodbus.GetHumidity(m, "sensor")

		averageHumidityBefore := calculateAverage(humidityData)
		// threshold := float64(averageHumidity) * (1 + cfg.PercentageIncrease/100)
		if float64(currentHumidity) > float64(averageHumidityBefore)+cfg.PercentageIncrease {
			fmt.Printf("Humidity spike detected! Current: %.2f%%, Average: %.2f%%, PercentageIncreaseThreshold: %.2f%%\n", float64(currentHumidity), float64(averageHumidityBefore), cfg.PercentageIncrease)
			refreshDuration, err := refreshDurationMinutes(cfg.RefreshDuration)
			if err != nil {
				fmt.Printf("Invalid refresh duration %d: %v\n", cfg.RefreshDuration, err)
				continue
			}
			err = systemairmodbus.ActivateRefresh(m, refreshDuration)
			if err != nil {
				fmt.Printf("Failed to activate Refresh mode during a humidity spike: %v", err)
			}
		}

		if len(humidityData) >= cfg.AveragePeriod {
			humidityData = humidityData[1:] // Remove the oldest reading
		}
		humidityData = append(humidityData, currentHumidity)

		averageHumidity := calculateAverage(humidityData)
		fmt.Printf("Current Humidity: %d%%, Average Humidity: %d%%\n", currentHumidity, averageHumidity)
	}
}

func refreshDurationMinutes(duration uint) (uint16, error) {
	if duration > math.MaxUint16 {
		return 0, fmt.Errorf("duration %d overflows uint16", duration)
	}
	return uint16(duration), nil
}

func calculateAverage(data []uint16) int {
	var sum uint16
	for _, v := range data {
		sum += v
	}
	return int(sum) / len(data)
}
