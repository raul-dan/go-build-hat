package buildhat

import (
	"buildhat/frame"
	"bytes"
	goserial "go.bug.st/serial"
	"go.uber.org/zap"
)

var serial = &serialConnection{
	portTTY: "/dev/serial0",
	portMode: &goserial.Mode{
		BaudRate: 115200,
	},
}

func sExec(cmd string) interface{} {
	resultFrame, err := serial.Execute(cmd)

	if err != nil {
		panic(err)
	}

	return resultFrame.GetContent()
}

type serialConnection struct {
	portTTY       string
	port          goserial.Port
	portMode      *goserial.Mode
	isOpened      bool
	pending       []pendingFrame
	readerStarted bool
}

type pendingFrame struct {
	c chan frame.Frame
	f frame.Frame
}

func (s *serialConnection) open() {
	if !s.isOpened {
		s.portMode = &goserial.Mode{
			BaudRate: 115200,
		}
		var err error
		s.port, err = goserial.Open("/dev/serial0", s.portMode)

		if err != nil {
			panic("Unable to open serial port. Received error: " + err.Error())
		}

		s.isOpened = true
	}
}

func (s *serialConnection) read() {
	if s.readerStarted {
		return
	}

	s.readerStarted = true

	go func() {
		s.open()
		delimiter := []byte("\r\n")
		page := make([]byte, 0)

		for {
			if len(s.pending) == 0 {
				continue
			}

			buff := make([]byte, 100)
			size, err := s.port.Read(buff)

			currentFrame := s.pending[0]

			if err != nil {
				logger.Error("Error reading from serial port", zap.Error(err))
				break
			}

			logger.Debug(
				"Read data from serial connection",
				zap.Int("size", size),
				zap.String("buff", string(buff[:size])),
			)

			page = append(page, buff[:size]...)
			pageSize := len(page)

			if pageSize <= 4 || !bytes.Equal(delimiter, page[:2]) || !bytes.Equal(delimiter, page[pageSize-2:pageSize]) {
				continue
			}

			pageContent := page[2 : pageSize-2]

			if currentFrame.f.IsEOF(pageContent) {
				s.pending = s.pending[1:]
				err := currentFrame.f.ParseBuffer(pageContent)

				if err != nil {
					panic("Unable to parse frame. Received error: " + err.Error())
				}

				page = make([]byte, 0)
				currentFrame.c <- currentFrame.f
				close(currentFrame.c)
			}
		}

	}()
}

func (s *serialConnection) write(cmd string, output chan frame.Frame) {
	newFrame, err := frame.NewFrame(cmd)

	if err != nil {
		panic("Unable to create frame for command " + cmd + ". Received error: " + err.Error())
	}

	s.open()
	s.pending = append(s.pending, pendingFrame{c: output, f: newFrame})
	s.port.Write(append([]byte(cmd), '\r'))
	s.read()
}

func (s *serialConnection) Execute(cmd string) (frame.Frame, error) {
	logger.Debug("Executing command", zap.String("cmd", cmd))

	output := make(chan frame.Frame)

	go func() {
		s.write(cmd, output)
	}()

	return <-output, nil
}
