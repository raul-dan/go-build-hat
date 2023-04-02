package hat

import "strconv"

const VoltageInCommand string = "vin"

type VoltageInFrame struct {
	voltage float32
}

func (f *VoltageInFrame) IsEOF(buff []byte) bool {
	return string(buff[len(buff)-2:]) == " V"
}

func (f *VoltageInFrame) ParseBuffer(buff []byte) error {
	voltage, err := strconv.ParseFloat(string(buff[:len(buff)-2]), 32)
	f.voltage = float32(voltage)

	return err
}

func (f *VoltageInFrame) GetContent() interface{} {
	return f.voltage
}
