package main

import (
	"fmt"
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

	for {
		select {
		case <-ticker.C:
			currentHumidity := systemairmodbus.GetHumidity(m, "sensor")

			averageHumidityBefore := calculateAverage(humidityData)
			// threshold := float64(averageHumidity) * (1 + cfg.PercentageIncrease/100)
			if float64(currentHumidity) > float64(averageHumidityBefore)+cfg.PercentageIncrease {
				fmt.Printf("Humidity spike detected! Current: %.2f%%, Average: %.2f%%, PercentageIncreaseThreshold: %.2f%%\n", float64(currentHumidity), float64(averageHumidityBefore), cfg.PercentageIncrease)
				err := systemairmodbus.ActivateRefresh(m, uint16(cfg.RefreshDuration))
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
}

func calculateAverage(data []uint16) int {
	var sum uint16
	for _, v := range data {
		sum += v
	}
	return int(sum) / len(data)
}
