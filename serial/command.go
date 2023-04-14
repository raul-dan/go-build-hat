package serial

type Command struct {
	cmd            []byte
	isVoid         bool // true if the command does not return any data
	isSubscription bool // true if the command is a subscription
	dto            Dto  // the Dto to use to parse the response
	callback       CommandCallback
}

type CommandCallback *func(data interface{})

func (c *Connection) registerCommand(cmd *Command) {
	c.connect()

	if cmd.isVoid {
		return
	}

	c.commandChannels[cmd] = make(chan interface{}, 1)
}

func (c *Connection) execute(cmd Command) interface{} {
	c.registerCommand(&cmd)
	c.write(cmd.cmd)

	if cmd.isSubscription {
		go func(cmd *Command) {
			callback := *cmd.callback
			select {
			case data := <-c.commandChannels[cmd]:
				callback(data)
			default:

			}
		}(&cmd)
	}

	if !cmd.isVoid && c.readingPaused {
		c.resumeRead()
	}

	if cmd.isSubscription || cmd.isVoid {
		return nil
	}

	return <-c.commandChannels[&cmd]
}
