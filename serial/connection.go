package serial

import (
	"buildhat/logger"
	"bytes"
	goserial "go.bug.st/serial"
	"go.uber.org/zap"
	"math"
	"reflect"
)

type Connection struct {
	portTTY        string
	port           goserial.Port
	portMode       *goserial.Mode
	isConnected    bool
	readingStarted bool
	channels       map[*Command]chan interface{}
	commands       []*Command
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

	c.isConnected = true
	c.boot()
}

func (c *Connection) continuousRead() {
	if c.readingStarted {
		return
	}

	c.readingStarted = true

	go func() {
		delimiter := []byte("\r\n")
		stalledBuff := make([]byte, 0)

		for {
			buff := make([]byte, 100)
			size, err := c.port.Read(buff)

			if err != nil {
				panic("Unable to read from serial port. Received error: " + err.Error())
			}

			stalledBuff = append(stalledBuff, buff[:size]...)
			stalledSize := len(stalledBuff)

			if stalledSize < 2 || !bytes.Equal(delimiter, stalledBuff[stalledSize-2:]) {
				continue
			}

			logger.Instance.Debug(
				"Read data from serial connection",
				zap.Int("size", stalledSize),
				zap.String("buff", string(stalledBuff)),
			)

			for _, line := range bytes.Split(stalledBuff, delimiter) {
				if len(line) == 0 {
					continue
				}

				for command, channel := range c.channels {
					if (*command).dto.BelongsTo(line) {
						(*command).dto = command.dto.Append(line)

						if !(*command).dto.IsComplete() {
							continue
						}

						channel <- (*command).dto.GetObject()

						if !command.isSubscription {
							close(channel)
							delete(c.channels, command)
						}
					}
				}
			}

			// reset buffer
			stalledBuff = make([]byte, 0)
		}
	}()
}

func (c *Connection) write(data interface{}) {
	if reflect.TypeOf(data).String() == "string" {
		// add separator if we are writing a string
		data = append([]byte(data.(string)), '\r')
	}

	// write bytes to serial port
	_, err := c.port.Write(data.([]byte))

	if err != nil {
		panic("Unable to write to serial port. Received error: " + err.Error())
	}

	logger.Instance.Debug(
		"Wrote data to serial connection",
		zap.ByteString("data", data.([]byte)[:int(math.Min(float64(len(data.([]byte))), 10))]),
	)

	c.continuousRead()
}
