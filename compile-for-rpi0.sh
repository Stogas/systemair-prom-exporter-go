echo "Building for Raspberry Pi Zero (ARMv6).."
env GOOS=linux GOARCH=arm GOARM=6 go build -o bin/systemair-exporter main/exporter.go
