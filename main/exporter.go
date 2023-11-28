package main

import (
    "fmt"
    "github.com/simonvetter/modbus"
)

func main() {
    fmt.Println("Starting")

    var client  *modbus.ModbusClient
    var err      error

    // for an RTU (serial) device/bus
    client, err = modbus.NewClient(&modbus.ClientConfiguration{
        URL:      "rtu:///dev/ttyUSB0",
        Speed:    19200,                   // default
        DataBits: 8,                       // default, optional
        Parity:   modbus.PARITY_NONE,      // default, optional
        StopBits: 1,                       // default if no parity, optional
        Timeout:  500 * time.Millisecond,
    })

    if err != nil {
        fmt.Println("Error - modbus client creation failed")
        // error out if client creation failed
    }
    fmt.Println("Client created")

    // now that the client is created and configured, attempt to connect
    err = client.Open()
    if err != nil {
        fmt.Println("Error - modbus client Open() failed")
        // error out if we failed to connect/open the device
        // note: multiple Open() attempts can be made on the same client until
        // the connection succeeds (i.e. err == nil), calling the constructor again
        // is unnecessary.
        // likewise, a client can be opened and closed as many times as needed.
    }
    fmt.Println("Client opened")

    // read a single 16-bit holding register at address 12543
    var reg16   uint16
    reg16, err  = client.ReadRegister(12543, modbus.HOLDING_REGISTER)
    if err != nil {
        fmt.Println("Error - modbus client reading failed")
      // error out
    } else {
      // use value
      fmt.Printf("value: %v", reg16)        // as unsigned integer
      fmt.Printf("value: %v", int16(reg16)) // as signed integer
    }

    // close the TCP connection/serial port
    client.Close()
}
