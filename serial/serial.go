package serial

import (
	goserial "go.bug.st/serial"
	"reflect"
)

var connection = &Connection{
	portTTY: "/dev/serial0",
	portMode: &goserial.Mode{
		BaudRate: 115200,
	},
	commandChannels:      map[*Command]chan interface{}{},
	readInterruptSignals: make(chan bool, 1),
}

func IsConnected() bool {
	return connection.isConnected
}

func Execute(command interface{}, dto Dto, cCallback CommandCallback) interface{} {
	if reflect.TypeOf(command).String() == "string" {
		command = append([]byte(command.(string)), '\r')
	}

	cmd := Command{
		cmd: command.([]byte),
		dto: dto,
	}

	if reflect.TypeOf(dto).String() == reflect.TypeOf(VoidDto{}).String() {
		cmd.isVoid = true
	}
	if cCallback != nil {
		cmd.isSubscription = true
		cmd.callback = cCallback
	}

	return connection.execute(cmd)
}
