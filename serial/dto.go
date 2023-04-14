package serial

import (
	"regexp"
)

type Dto interface {
	Append(buffer []byte) Dto
	IsComplete() bool
	BelongsTo(buffer []byte) bool
	GetObject() interface{}
}

// RegexpDto is used to match a series of buffers against a series of regular
// expressions. The DTO will be considered complete when all regular expressions
// have been matched.
type RegexpDto struct {
	matchedBuffer [][]byte
	Patterns      []*regexp.Regexp
}

func (s RegexpDto) Append(buffer []byte) Dto {
	s.matchedBuffer = append(s.matchedBuffer, buffer)
	return s
}

func (s RegexpDto) IsComplete() bool {
	return len(s.matchedBuffer) == len(s.Patterns)
}

func (s RegexpDto) BelongsTo(buffer []byte) bool {
	return s.Patterns[len(s.matchedBuffer)].Match(buffer)
}

func (s RegexpDto) GetObject() interface{} {
	return s.matchedBuffer
}

// A VoidDto is a DTO that does not expect any data to be returned from the
// serial connection. Write and forget.
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
