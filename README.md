# Systemair-prom-exporter-go

This is a Go rewrite of my Python-based [systemair-prom-exporter](https://gitlab.com/pabaisa/systemair-prom-exporter).

This rewrite is caused by:
- My want to reduce load on a Raspberry Pi Zero W, where the Python version pins the CPU to 100%
- My want to learn Go

Code might be shitty in places, but hey - it works fine. Feel free to submit PRs if you want! 😊 

As with the Python version, the tested set up is:
- Raspberry Pi Zero W
- A simple [USB to RS485 Adapter](https://web.archive.org/web/20180424082558/http://www.dx.com/p/usb-to-rs485-adapter-black-green-296620)
- [Systemair Save VTR 150/B](https://www.systemair.com/en/p/save-vtr-150-b-l-1000w-396937)

This includes the [2019 version of the Systemair Modbus reference map](https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20190116__REV__29_.PDF) and the [2021 version](https://shop.systemair.com/upload/assets/SAVE_MODBUS_VARIABLE_LIST_20210301_REV36.PDF)

#### Exit codes

- Exit code of `2` signals Modbus client creation, configuration, or connection failures at startup
- Exit code of `3` signals Prometheus exporter related issues

Register read/write failures during normal operation are logged and returned as errors; they do not terminate the process. Prometheus collectors skip failed series; REST `/hvac/status` returns HTTP 500 if any field cannot be read.

#### TODO

- Do not hardcode `hvac_` metric namespace prefix
- Do not hardcode port `9999`
- Do not hardcode HTTP metrics path `/metrics`
- Do not hardcode modbus config, including using `/dev/ttyUSB0` as the modbus parget path, speed, etc.
- Allow enabling/disabling specific metric subsystems (`temp`, etc.)
- Implement structured logging
- Expand functionality to include write capabilities (I want to enable refresh mode based on external decisions)