package frame

import (
	"buildhat/frame/hat"
	"errors"
)

type Frame interface {
	IsEOF(buff []byte) bool
	ParseBuffer(buff []byte) error
	GetContent() interface{}
}

func NewFrame(cmd string) (Frame, error) {
	switch cmd {
	case hat.VersionCommand:
		return &hat.FirmwareVersionFrame{}, nil

	case hat.VoltageInCommand:
		return &hat.VoltageInFrame{}, nil

	case InitCommand:
		return &VoidFrame{}, nil

	default:
		return nil, errors.New("Unknown command: " + cmd)

	}
}
