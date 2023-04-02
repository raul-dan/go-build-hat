package buildhat

import (
	"buildhat/frame"
	"buildhat/frame/hat"
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

		logger.Debug("Opening serial port", zap.String("port", s.portTTY), zap.Int("baud", s.portMode.BaudRate))

		var err error
		s.port, err = goserial.Open("/dev/serial0", s.portMode)

		if err != nil {
			panic("Unable to open serial port. Received error: " + err.Error())
		}

		s.isOpened = true
		s.init()
	}
}

func (s *serialConnection) init() {
	logger.Debug("Initializing HAT")
	logger.Debug("Hat version: " + sExec(hat.VersionCommand).(string))
	sExec(frame.InitCommand)
}

func (s *serialConnection) read() {
	if s.readerStarted {
		return
	}

	s.readerStarted = true

	go func() {
		delimiter := []byte("\r\n")
		page := make([]byte, 0)

		for {
			if len(s.pending) == 0 {
				continue
			}

			buff := make([]byte, 100)
			size, err := s.port.Read(buff)

			currentFrame := s.pending[0]
			_, currentFrameIsVoid := currentFrame.f.(*frame.VoidFrame)

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

			if (pageSize <= 4 && !currentFrameIsVoid) ||
				!bytes.Equal(delimiter, page[:2]) ||
				!bytes.Equal(delimiter, page[pageSize-2:pageSize]) {
				continue
			}

			if !currentFrameIsVoid {
				page = page[2 : pageSize-2]
			}

			if currentFrame.f.IsEOF(page) {
				s.pending = s.pending[1:]
				err := currentFrame.f.ParseBuffer(page)

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

func (s *serialConnection) write(data []byte) {
	logger.Debug("Writing data to serial port", zap.String("data", string(data)))
	s.port.Write(data)
}

func (s *serialConnection) readFrame(cmd string, output chan frame.Frame) {
	s.open()

	newFrame, err := frame.NewFrame(cmd)

	if err != nil {
		panic("Unable to create frame for command " + cmd + ". Received error: " + err.Error())
	}

	s.pending = append(s.pending, pendingFrame{c: output, f: newFrame})
	s.write(append([]byte(cmd), '\r'))
	s.read()
}

func (s *serialConnection) Execute(cmd string) (frame.Frame, error) {
	logger.Debug("Executing command", zap.String("cmd", cmd))

	output := make(chan frame.Frame)

	go func() {
		s.readFrame(cmd, output)
	}()

	return <-output, nil
}
