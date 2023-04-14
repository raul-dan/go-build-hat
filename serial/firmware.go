package serial

import (
	"strings"
)

const versionPrefix = "Firmware version: "
const bootloaderPrefix = "BuildHAT bootloader version "

type HatVersionDto struct {
	state *BootLoaderState
	Boot  bool
}

type BootLoaderState struct {
	RequiresFirmware bool
	version          string
}

func (v HatVersionDto) Append(buffer []byte) Dto {
	if v.state == nil {
		v.state = &BootLoaderState{}
	}

	if string(buffer) == "version" {
		return v
	}

	var prefix = ""

	if strings.HasPrefix(string(buffer), bootloaderPrefix) {
		v.state.RequiresFirmware = true
		prefix = bootloaderPrefix
	} else {
		prefix = versionPrefix
	}

	v.state.version = string(buffer[len(prefix):])
	return v
}

func (v HatVersionDto) IsComplete() bool {
	return v.state != nil && v.state.version != ""
}

func (v HatVersionDto) GetObject() interface{} {
	if v.Boot {
		return v.state
	}

	return v.state.version
}

func (v HatVersionDto) BelongsTo(buffer []byte) bool {
	return string(buffer) == "version" ||
		strings.HasPrefix(string(buffer), versionPrefix) ||
		strings.HasPrefix(string(buffer), bootloaderPrefix)
}
