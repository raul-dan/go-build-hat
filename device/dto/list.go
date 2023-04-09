package dto

import (
	"buildhat/serial"
	"strings"
)

type ListDevicesDto struct {
	data string
}

func (l ListDevicesDto) Append(buffer []byte) serial.Dto {
	l.data = l.data + "\n" + string(buffer)
	return l
}

func (l ListDevicesDto) IsComplete() bool {
	return strings.HasSuffix(l.data, "P3: no device detected")
}

func (l ListDevicesDto) BelongsTo(buffer []byte) bool {
	return true
}

func (l ListDevicesDto) GetObject() interface{} {
	return l.data
}
