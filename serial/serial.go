package serial

import (
	goserial "go.bug.st/serial"
	"reflect"
)

var connection *Connection = &Connection{
	portTTY: "/dev/serial0",
	portMode: &goserial.Mode{
		BaudRate: 115200,
	},
	channels: map[*Command]chan interface{}{},
}

func IsConnected() bool {
	return connection.isConnected
}

func Execute(command string, dto Dto, cCallback CommandCallback) interface{} {
	cmd := Command{
		cmd: command,
		dto: dto,
	}

	if reflect.TypeOf(dto).String() == "serial.SimpleDto" {
		cmd.isVoid = true
	}
	if cCallback != nil {
		cmd.isSubscription = true
		cmd.callback = cCallback
	}

	return connection.execute(cmd)
}
