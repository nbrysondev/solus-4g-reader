// build
// env GOOS=linux GOARCH=arm GOARM=5 go build
package SolusReader

import (
    "time"
    "encoding/json"
    "github.com/simonvetter/modbus"
)

func GetSolusReadings() string {
    var client  *modbus.ModbusClient
    var err      error

    // for an RTU (serial) device/bus
    client, err = modbus.NewClient(&modbus.ClientConfiguration{
        URL:      "rtu:///dev/ttyUSB0",
        Speed:    9600,                   // default
        DataBits: 8,                       // default, optional
        Parity:   modbus.PARITY_NONE,      // default, optional
        StopBits: 1,                       // default if no parity, optional
        Timeout:  300 * time.Millisecond,
    })

    if err != nil {
        // error out if client creation failed
    }

    // now that the client is created and configured, attempt to connect
    err = client.Open()
    if err != nil {
        // error out if we failed to connect/open the device
        // note: multiple Open() attempts can be made on the same client until
        // the connection succeeds (i.e. err == nil), calling the constructor again
        // is unnecessary.
        // likewise, a client can be opened and closed as many times as needed.
    }

    type SolarData struct {
        RealTime uint32 `json:"realtime"`
        Today float32 `json:"today"`
        Yesterday float32 `json:"yesterday"`
        Month uint32 `json:"month"`
        Year uint32 `json:"year"`
    }

    var realtime   uint32
    var today uint16
    var yesterday uint16
    var thismonth uint32
    var thisyear uint32
    
	today, _ = client.ReadRegister(3014, modbus.INPUT_REGISTER)
    yesterday, _ = client.ReadRegister(3015, modbus.INPUT_REGISTER)
    thismonth, _ = client.ReadUint32(3010, modbus.INPUT_REGISTER)
    thisyear, _ = client.ReadUint32(3016, modbus.INPUT_REGISTER)
    realtime, _ = client.ReadUint32(3004, modbus.INPUT_REGISTER)
    client.Close()
    
    data := SolarData{}
    data.RealTime = realtime
    data.Today = float32(today) / 10
    data.Yesterday = float32(yesterday) / 10
    data.Month = thismonth
    data.Year = thisyear
    json, err := json.Marshal(data)
    return string(json)
}