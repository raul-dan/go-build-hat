package buildhat

import "reflect"

type BatchFrame struct {
	FrameType   Frame
	frames      []Frame
	framesCount uint32
	countedEOFs uint32
}

func (f *BatchFrame) IsEOF(buff []byte) bool {
	if f.framesCount == 0 {
		panic("Batch frame is empty")
	}

	if f.framesCount == f.countedEOFs {
		panic("Batch frame overflow")
	}

	if f.FrameType.IsEOF(buff) {
		f.countedEOFs++
	}

	return f.framesCount == f.countedEOFs
}

func (f *BatchFrame) ParseBuffer(buff []byte) error {
	f.frames = make([]Frame, f.framesCount)

	for i := 0; i < int(f.framesCount); i++ {
		f.frames[i] = reflect.New(reflect.TypeOf(f.FrameType)).Interface().(Frame)

	}

	return nil
}

func (f *BatchFrame) GetContent() []interface{} {
	frameContents := make([]interface{}, f.framesCount)

	for i := 0; i < int(f.framesCount); i++ {
		frameContents[i] = f.frames[i].GetContent()
	}

	return frameContents
}
