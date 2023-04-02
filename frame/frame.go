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
	case "version":
		return &hat.FirmwareVersionFrame{}, nil

	case "vin":
		return &hat.VoltageInFrame{}, nil

	default:
		return nil, errors.New("Unknown command: " + cmd)

	}
}
