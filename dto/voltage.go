package dto

import (
	"strconv"
)

type VoltageDto struct {
	voltage float32
}

func (v VoltageDto) IngestBuffer(buffer []byte) Dto {
	voltage, err := strconv.ParseFloat(string(buffer[:len(buffer)-2]), 32)

	if err != nil {
		panic("Unexpected error while parsing voltage: " + err.Error())
	}

	v.voltage = float32(voltage)

	return v
}

func (v VoltageDto) IsComplete() bool {
	return v.voltage != 0
}

func (v VoltageDto) Matches(buffer []byte) bool {
	return len(buffer) > 2 && string(buffer[len(buffer)-2:]) == " V"
}

func (v VoltageDto) GetObject() interface{} {
	return v.voltage
}
