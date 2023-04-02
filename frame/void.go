package frame

import "bytes"

const InitCommand string = "port 0 ; select ; port 1 ; select ; port 2 ; select ; port 3 ; select ; echo 0"

type VoidFrame struct {
}

func (f *VoidFrame) IsEOF(buff []byte) bool {
	return bytes.Equal([]byte("\r\n"), buff)
}

func (f *VoidFrame) ParseBuffer(buff []byte) error {
	return nil
}

func (f *VoidFrame) GetContent() interface{} {
	return nil
}
