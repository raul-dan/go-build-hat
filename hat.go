package buildhat

import (
	"buildhat/device"
	"buildhat/frame/hat"
)

type Hat struct {
}

func (h *Hat) GetFirmwareVersion() string {
	return sExec(hat.VersionCommand).(string)
}

func (h *Hat) GetVin() float32 {
	return sExec(hat.VoltageInCommand).(float32)
}

func (h *Hat) GetDevices() []device.Type {
	return sExec("list").([]device.Type)
}
