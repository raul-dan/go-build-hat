package buildhat

import (
	frame "buildhat/frame"
	"bytes"
	"go.bug.st/serial"
	"go.uber.org/zap"
)

var Serial = &SerialConnection{
	portTTY: "/dev/serial0",
	portMode: &serial.Mode{
		BaudRate: 115200,
	},
}

func sExec(cmd string) interface{} {
	frame, err := Serial.Execute(cmd)

	if err != nil {
		panic(err)
	}

	return frame.GetContent()
}

type SerialConnection struct {
	portTTY       string
	port          serial.Port
	portMode      *serial.Mode
	isOpened      bool
	pending       []pendingFrame
	readerStarted bool
}

type pendingFrame struct {
	c chan frame.Frame
	f frame.Frame
}

func (s *SerialConnection) open() {
	if !s.isOpened {
		s.portMode = &serial.Mode{
			BaudRate: 115200,
		}
		var err error
		s.port, err = serial.Open("/dev/serial0", s.portMode)

		if err != nil {
			panic("Unable to open serial port. Received error: " + err.Error())
		}

		s.isOpened = true
	}
}

func (s *SerialConnection) read() {
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
				currentFrame.f.ParseBuffer(pageContent)
				page = make([]byte, 0)
				currentFrame.c <- currentFrame.f
				close(currentFrame.c)
			}
		}

	}()
}

func (s *SerialConnection) write(cmd string, output chan frame.Frame) {
	newFrame, err := frame.NewFrame(cmd)

	if err != nil {
		panic("Unable to create frame for command " + cmd + ". Received error: " + err.Error())
	}

	s.open()
	s.pending = append(s.pending, pendingFrame{c: output, f: newFrame})
	s.port.Write(append([]byte(cmd), '\r'))
	s.read()
}

func (s *SerialConnection) Execute(cmd string) (frame.Frame, error) {
	logger.Debug("Executing command", zap.String("cmd", cmd))

	output := make(chan frame.Frame)

	go func() {
		s.write(cmd, output)
	}()

	return <-output, nil
}
