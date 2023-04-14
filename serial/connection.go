package serial

import (
	"buildhat/logger"
	"bytes"
	goserial "go.bug.st/serial"
	"go.uber.org/zap"
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
	commandChannels      map[*Command]chan interface{}
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

	runWithTimeout(time.Second*10, func() interface{} {
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

		delimiter := []byte("\r\n")
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
				stalledSize := len(stalledBuff)

				if stalledSize < 2 || !bytes.Equal(delimiter, stalledBuff[stalledSize-2:]) {
					continue
				}

				for _, line := range bytes.Split(stalledBuff, delimiter) {
					if len(line) == 0 {
						continue
					}

					if c.broadcastBuffer(line) {
						// reset buffer
						stalledBuff = make([]byte, 0)
					}
				}

			}
		}

	}()
}

func (c *Connection) broadcastBuffer(buffer []byte) bool {
	fullyIngested := false

	for command, channel := range c.commandChannels {
		if (*command).dto.BelongsTo(buffer) {
			(*command).dto = command.dto.Append(buffer)

			if !(*command).dto.IsComplete() {
				continue
			}

			channel <- (*command).dto.GetObject()
			fullyIngested = true

			if !(*command).isSubscription {
				close(channel)
				delete(c.commandChannels, command)
			}
		}
	}

	return fullyIngested
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
	}, "Unable to pause serial connection read")
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
