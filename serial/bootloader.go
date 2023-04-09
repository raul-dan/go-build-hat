package serial

import (
	"buildhat/logger"
	"fmt"
	"os"
	"path/filepath"
)

var isBooted = false

func (c *Connection) boot() {
	if isBooted {
		return
	}

	isBooted = true

	var dtoResult = Execute("version", HatVersionDto{Boot: true}, nil).(*BootLoaderState)

	if !(*dtoResult).RequiresFirmware {
		logger.Instance.Debug("HAT already booted")
		return
	}

	logger.Instance.Debug("Booting HAT")
	Execute("clear", SimpleDto{ExpectedReply: []byte("BHBL> clear")}, nil)

	var firmwarePath, _ = filepath.Abs("./data/firmware.bin")
	var firmwareBinary, _ = os.ReadFile(firmwarePath)

	/**
	 * Load firmware
	 */
	connection.write([]byte{0x02})
	connection.write(firmwareBinary)
	connection.write([]byte{0x03, '\r'})

	fmt.Println(dtoResult)
	fmt.Println(dtoResult)

}
