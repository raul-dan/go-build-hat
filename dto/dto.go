package dto

import (
	"regexp"
)

type Dto interface {
	Matches(buffer []byte) bool
	IngestBuffer(buffer []byte) Dto
	GetObject() interface{}
}

type LineByLineDto interface {
	IsComplete() bool
	Dto
}

type SubscriptionDto interface {
	Reset() Dto
	Callback(object interface{})
	LineByLineDto
}

func IsSubscription(dto Dto) bool {
	_, isSubscription := dto.(SubscriptionDto)
	return isSubscription
}

func IsLineByLine(dto Dto) bool {
	_, ok := dto.(LineByLineDto)
	return ok
}

// RegexpDto is used to match a series of buffers against a series of regular
// expressions. The DTO will be considered complete when all regular expressions
// have been matched.
type RegexpDto struct {
	matchedBuffer [][]byte
	Patterns      []*regexp.Regexp
}

func (s RegexpDto) IngestBuffer(buffer []byte) Dto {
	s.matchedBuffer = append(s.matchedBuffer, buffer)
	return s
}

func (s RegexpDto) IsComplete() bool {
	return len(s.matchedBuffer) == len(s.Patterns)
}

func (s RegexpDto) Matches(buffer []byte) bool {
	return s.Patterns[len(s.matchedBuffer)].Match(buffer)
}

func (s RegexpDto) GetObject() interface{} {
	return s.matchedBuffer
}

// A VoidDto is a DTO that does not expect any data to be returned from the
// serial Connection. Write and forget.
type VoidDto struct {
}

func (v VoidDto) IngestBuffer(buffer []byte) Dto {
	panic("Void DTO should not be processed")
}

func (v VoidDto) IsComplete() bool {
	panic("Void DTO should not be processed")
}

func (v VoidDto) Matches(buffer []byte) bool {
	panic("Void DTO should not be processed")
}

func (v VoidDto) GetObject() interface{} {
	panic("Void DTO should not be processed")
}
