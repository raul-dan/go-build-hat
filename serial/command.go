package serial

import (
	"buildhat/dto"
)

type Command struct {
	cmd     []byte
	isVoid  bool    // true if the command does not return any data
	dto     dto.Dto // the Dto to use to parse the response
	channel chan interface{}
}

type CommandCallback *func(data interface{})

func (c *Connection) registerCommand(cmd *Command) {
	c.connect()

	if cmd.isVoid {
		return
	}

	c.commands = append(c.commands, cmd)
}

func (c *Connection) removeCommand(command *Command) {
	for i, cmd := range c.commands {
		if cmd == command {
			c.commands = append(c.commands[:i], c.commands[i+1:]...)
			close(cmd.channel)
			break
		}
	}
}

func (c *Connection) execute(cmd *Command) {
	c.registerCommand(cmd)
	c.write((*cmd).cmd)

	if !cmd.isVoid && c.readingPaused {
		c.resumeRead()
	}
}
