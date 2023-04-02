package buildhat

import "errors"

type Frame interface {
	IsEOF(buff []byte) bool
	ParseBuffer(buff []byte) error
	GetContent() interface{}
}

func NewFrame(cmd string) (Frame, error) {
	switch cmd {
	case "version":
		return &FirmwareVersionFrame{}, nil

	case "vin":
		return &VoltageInFrame{}, nil

	default:
		return nil, errors.New("Unknown command: " + cmd)

	}
}
