package serial

import "bytes"

type Dto interface {
	Append(buffer []byte) Dto
	IsComplete() bool
	BelongsTo(buffer []byte) bool
	GetObject() interface{}
}

type SimpleDto struct {
	ExpectedReply []byte
}

func (s SimpleDto) Append(buffer []byte) Dto {
	return s
}

func (s SimpleDto) IsComplete() bool {
	return true
}

func (s SimpleDto) BelongsTo(buffer []byte) bool {
	return bytes.Equal(s.ExpectedReply, buffer)
}

func (s SimpleDto) GetObject() interface{} {
	return s.ExpectedReply
}

type VoidDto struct {
}

func (v VoidDto) Append(buffer []byte) Dto {
	panic("Void DTO should not be processed")
}

func (v VoidDto) IsComplete() bool {
	panic("Void DTO should not be processed")
}

func (v VoidDto) BelongsTo(buffer []byte) bool {
	panic("Void DTO should not be processed")
}

func (v VoidDto) GetObject() interface{} {
	panic("Void DTO should not be processed")
}
