package device

import (
	"buildhat/device/dto"
	"buildhat/serial"
)

type Hat struct {
	version string
}

func (h *Hat) GetFirmwareVersion() string {
	if h.version == "" {
		h.version = serial.Execute("version", serial.HatVersionDto{}, nil).(string)
	}

	return h.version
}

func (h *Hat) GetVin() float32 {
	return serial.Execute("vin", dto.VoltageDto{}, nil).(float32)
}

func (h *Hat) GetDevices() string {
	return serial.Execute("list", dto.ListDevicesDto{}, nil).(string)
}
