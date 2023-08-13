package serial

import (
	serialdto "buildhat/dto"
	"buildhat/logger"
	"bytes"
	goserial "go.bug.st/serial"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"math"
	"time"
)

type Connection struct {
	portTTY              string
	port                 goserial.Port
	portMode             *goserial.Mode
	isConnected          bool
	readingStarted       bool
	readingPaused        bool
	readInterruptSignals chan bool
	commands             []*Command
}

func (c *Connection) connect() {
	if c.isConnected {
		return
	}

	var err error
	c.port, err = goserial.Open("/dev/serial0", c.portMode)

	if err != nil {
		panic("Unable to open serial port. Received error: " + err.Error())
	}

	err = c.port.SetReadTimeout(time.Millisecond * 100)

	if err != nil {
		panic("Unable to set serial port read timeout. Received error: " + err.Error())
	}

	c.isConnected = true

	runWithTimeout(time.Minute*20, func() interface{} {
		c.boot()
		return nil
	}, "Unable to boot HAT")
}

func (c *Connection) startRead() {
	if c.readingStarted || c.readingPaused {
		return
	}

	go func() {
		c.readingStarted = true

		stalledBuff := make([]byte, 0)

		for {
			select {
			case <-c.readInterruptSignals:
				c.readingStarted = false
				c.readingPaused = true
				return

			default:
				buff := make([]byte, 32)
				size, err := c.port.Read(buff)

				if err != nil {
					panic("Unable to read from serial port. Received error: " + err.Error())
				}

				if size == 0 {
					continue
				}

				logger.Instance.Debug(
					"Read data from serial connection",
					zap.Int("size", size),
					zap.String("buff", string(buff[:size])),
				)

				stalledBuff = append(stalledBuff, buff[:size]...)
				stalledBuff = c.ingestBuffer(stalledBuff)
			}
		}

	}()
}

func (c *Connection) ingestBuffer(buffer []byte) []byte {
	delimiter := []byte("\r\n")

	if len(buffer) < 2 || !bytes.Contains(buffer, delimiter) {
		return buffer
	}

	unprocessedBuffer := make([]byte, 0)
	inputLines := bytes.Split(buffer, delimiter)

	if len(inputLines[len(inputLines)-1]) > 0 {
		// the last line is not complete, we're going to keep it aside
		unprocessedBuffer = inputLines[len(inputLines)-1]
		inputLines = inputLines[:len(inputLines)-1]
	}

	if len(inputLines) >= 1 && len(inputLines[0]) == 0 {
		// also remove any potential leading blank line
		inputLines = inputLines[1:]
	}

	inputLines, commandsPendingRelease := c.pipeToMultiLineCommands(inputLines)
	inputLines, successfulCommands := c.pipeToLineByLineCommands(inputLines)
	commandsPendingRelease = append(commandsPendingRelease, successfulCommands...)

	for _, command := range commandsPendingRelease {
		(*command).channel <- (*command).dto.GetObject()
		logger.Instance.Debug("Releasing command", zap.String("command", string((*command).cmd)))

		if serialdto.IsSubscription((*command).dto) {
			(*command).dto = (*command).dto.(serialdto.SubscriptionDto).Reset()
		} else {
			c.removeCommand(command)
		}
	}

	return bytes.Join(append(inputLines, unprocessedBuffer), delimiter)
}

func (c *Connection) pipeToMultiLineCommands(lines [][]byte) ([][]byte, []*Command) {
	remainingBuffer := make([][]byte, 0)
	var commands []*Command

	for _, line := range lines {
		remainingBuffer = append(remainingBuffer, line)
		bufferWithNl := bytes.Join(remainingBuffer, []byte("\n"))

		for _, command := range c.commands {
			if serialdto.IsLineByLine((*command).dto) || !(*command).dto.Matches(bufferWithNl) {
				continue
			}

			(*command).dto = command.dto.IngestBuffer(bufferWithNl)
			commands = append(commands, command)
			remainingBuffer = make([][]byte, 0)
		}
	}

	return remainingBuffer, commands
}

func (c *Connection) pipeToLineByLineCommands(lines [][]byte) ([][]byte, []*Command) {
	remainingBuffer := make([][]byte, 0)
	var completeCommands []*Command

	for _, line := range lines {
		ingested := false

		for _, command := range c.commands {
			if slices.Contains(completeCommands, command) || !(*command).dto.Matches(line) {
				continue
			}

			(*command).dto = command.dto.IngestBuffer(line)
			ingested = true

			if (*command).dto.(serialdto.LineByLineDto).IsComplete() {
				completeCommands = append(completeCommands, command)
			}
		}

		if !ingested {
			remainingBuffer = append(remainingBuffer, line)
		}
	}

	return remainingBuffer, completeCommands
}

func (c *Connection) pauseRead() {
	if !c.readingStarted || c.readingPaused {
		return
	}

	runWithTimeout(time.Second*2, func() interface{} {
		c.readInterruptSignals <- true

		for {
			if c.readingPaused {
				break
			}
		}

		return nil
	}, "Unable to pause serial hatConnection read")
}

func (c *Connection) resumeRead() {
	c.readingPaused = false
	c.startRead()
}

func (c *Connection) write(data []byte) {
	// write bytes to serial port
	_, err := c.port.Write(data)

	if err != nil {
		panic("Unable to write to serial port. Received error: " + err.Error())
	}

	logger.Instance.Debug(
		"Wrote data to serial connection",
		zap.ByteString("data", data[:int(math.Min(float64(len(data)), 20))]),
	)

	c.startRead()
}
