package hat

import "strings"

const VersionCommand string = "version"
const versionPrefix = "Firmware version: "

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
