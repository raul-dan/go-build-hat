package hat

import "buildhat/device"

type DevicesFrame struct {
	devices []device.Type
}

func (f *DevicesFrame) IsEOF(buff []byte) bool {
	return string(buff) == "OK"
}

func (f *DevicesFrame) ParseBuffer(buff []byte) error {
	return nil
}

func (f *DevicesFrame) GetContent() interface{} {
	return f.devices
}
