package serial

import (
	"buildhat/logger"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var isBooted = false

func (c *Connection) boot() {
	if isBooted {
		return
	}

	isBooted = true

	var bootLoaderState = Execute("version", HatVersionDto{Boot: true}, nil).(*BootLoaderState)

	if !(*bootLoaderState).RequiresFirmware {
		logger.Instance.Debug("HAT already booted")
		return
	}

	loadFirmware()
}

func loadFirmware() {
	firmwarePath, _ := filepath.Abs("./data/firmware.bin")
	firmwareBinary, _ := os.ReadFile(firmwarePath)
	firmwareChecksum := checksum(firmwareBinary)

	promptDto := RegexpDto{
		Patterns: []*regexp.Regexp{regexp.MustCompile(`(?i)^BHBL>(\s)?(version|clear)?$`)},
	}

	logger.Instance.Debug("Booting HAT")
	Execute("clear", promptDto, nil)

	// prepare payload
	firmwarePayload := append([]byte{0x02}, append(firmwareBinary, []byte{0x03, '\r'}...)...)

	// suspend reading until full payload is sent. reading will resume automatically after write completes
	connection.pauseRead()

	// send firmware size and checksum
	Execute(fmt.Sprintf("load %d %d", len(firmwareBinary), firmwareChecksum), VoidDto{}, nil)

	// and finally send the payload
	Execute(firmwarePayload, RegexpDto{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)^Image Received$"),
			regexp.MustCompile("(?i)^Checksum OK$"),
		},
	}, nil)

	// send signature
	writeFirmwareSignature()

	// and finally reboot
	Execute("reboot", RegexpDto{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)^Done initialising ports$"),
		},
	}, nil)

	logger.Instance.Debug("Firmware loaded successfully")
}

func writeFirmwareSignature() {
	signaturePath, _ := filepath.Abs("./data/signature.bin")
	signatureBinary, _ := os.ReadFile(signaturePath)

	// send signature size
	Execute(
		[]byte(fmt.Sprintf("signature %d\r", len(signatureBinary))),
		VoidDto{}, nil,
	)

	// prepare payload
	signaturePayload := append([]byte{0x02}, append(signatureBinary, []byte{0x03, '\r'}...)...)

	// and finally send the payload
	Execute(signaturePayload, RegexpDto{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)^Signature received$"),
		},
	}, nil)
}
