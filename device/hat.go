package device

import (
	"buildhat/dto"
	"buildhat/serial"
)

type Hat struct {
	version string
}

func (h *Hat) GetFirmwareVersion() string {
	if h.version == "" {
		h.version = serial.Execute("version", dto.HatVersionDto{}).(string)
	}

	return h.version
}

func (h *Hat) GetVin() float32 {
	return serial.Execute("vin", dto.VoltageDto{}).(float32)
}

func (h *Hat) GetDevices() ConnectedInventory {
	return serial.Execute("list", dto.ListDevicesDto{}).(ConnectedInventory)
}
