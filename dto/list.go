package dto

import (
	"regexp"
	"strings"
)

var deviceMatchRegexp = regexp.MustCompile(
	`^P(?P<port>[0-3]): ((?:connected to active ID (?P<typeId>[0-9a-fA-F]+)\ntype (?P<typeIdConfirmation>[0-9a-fA-F]+)(?:.*?)position PID: (?:.*?))|(?:no device detected))$`,
)

type ListDevicesDto struct {
	dtos []*Dto
	data string
}

func (l ListDevicesDto) IngestBuffer(buffer []byte) Dto {
	l.data = l.data + "\n" + string(buffer)
	return l
}

func (l ListDevicesDto) IsComplete() bool {
	return strings.HasSuffix(l.data, "P3: no device detected")
}

func (l ListDevicesDto) Matches(buffer []byte) bool {
	return true
}

func (l ListDevicesDto) GetObject() interface{} {
	return l.data
}
