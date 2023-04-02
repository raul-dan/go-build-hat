package hat

import "strconv"

type VoltageInFrame struct {
	voltage float64
}

func (f *VoltageInFrame) IsEOF(buff []byte) bool {
	return string(buff[len(buff)-2:]) == " V"
}

func (f *VoltageInFrame) ParseBuffer(buff []byte) error {
	voltage, err := strconv.ParseFloat(string(buff[:len(buff)-2]), 64)
	f.voltage = voltage

	return err
}

func (f *VoltageInFrame) GetContent() interface{} {
	return f.voltage
}
