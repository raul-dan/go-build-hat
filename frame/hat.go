package buildhat

import (
	"strconv"
	"strings"
)

const versionPrefix = "Firmware version: "

/************* Frame FirmwareVersionFrame *************/

type FirmwareVersionFrame struct {
	version string
}

func (f *FirmwareVersionFrame) IsEOF(buff []byte) bool {
	if !strings.HasPrefix(string(buff), "Firmware version:") {
		panic("Buffer is not a valid firmware version frame. Received: " + string(buff))
	}

	return true
}

func (f *FirmwareVersionFrame) ParseBuffer(buff []byte) error {
	f.version = string(buff[len(versionPrefix):])
	return nil
}

func (f *FirmwareVersionFrame) GetContent() interface{} {
	return f.version
}

/************* Frame Voltage IN *************/

type VoltageInFrame struct {
	voltage float64
}

func (f *VoltageInFrame) IsEOF(buff []byte) bool {
	return string(buff[len(buff)-2:]) == " V"
}

func (f *VoltageInFrame) ParseBuffer(buff []byte) error {
	voltage, err := strconv.ParseFloat(string(buff[:len(buff)-2]), 64)
	f.voltage = voltage

	return err
}

func (f *VoltageInFrame) GetContent() interface{} {
	return f.voltage
}
