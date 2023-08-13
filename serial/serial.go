package serial

import (
	serialdto "buildhat/dto"
	goserial "go.bug.st/serial"
	"reflect"
)

var hatConnection = &Connection{
	portTTY: "/dev/serial0",
	portMode: &goserial.Mode{
		BaudRate: 115200,
	},
	commands:             make([]*Command, 0),
	readInterruptSignals: make(chan bool, 1),
}

func IsConnected() bool {
	return hatConnection.isConnected
}

func Execute(command interface{}, dto serialdto.Dto) interface{} {
	if reflect.TypeOf(command).String() == "string" {
		command = append([]byte(command.(string)), '\r')
	}

	cmd := Command{
		cmd:     command.([]byte),
		dto:     dto,
		channel: make(chan interface{}, 1),
	}

	if reflect.TypeOf(dto).String() == reflect.TypeOf(serialdto.VoidDto{}).String() {
		cmd.isVoid = true
	}

	hatConnection.execute(&cmd)

	if cmd.isVoid {
		close(cmd.channel)
		return nil
	}

	if serialdto.IsSubscription(cmd.dto) {
		unsubscribe := make(chan bool, 1)

		go func() {
			for {
				select {
				case data, ok := <-cmd.channel:
					if !ok {
						continue
					}
					cmd.dto.(serialdto.SubscriptionDto).Callback(data)
				case <-unsubscribe:

					return
				}
			}
		}()

		return func() {
			unsubscribe <- true
		}
	}

	return <-cmd.channel
}
