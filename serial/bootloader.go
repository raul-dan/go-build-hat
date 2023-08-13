package serial

import (
	"buildhat/dto"
	"buildhat/logger"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var isBooted = false

func (c *Connection) boot() {
	if isBooted {
		return
	}

	isBooted = true

	var bootLoaderState = Execute("version", dto.HatVersionDto{Boot: true}).(*dto.BootLoaderState)

	if !(*bootLoaderState).RequiresFirmware {
		logger.Instance.Debug("HAT already booted")
	} else {
		loadFirmware()
	}

	logger.Instance.Debug("HAT booted successfully")
}

func loadFirmware() {
	firmwarePath, _ := filepath.Abs("./data/firmware.bin")
	firmwareBinary, _ := os.ReadFile(firmwarePath)
	firmwareChecksum := checksum(firmwareBinary)

	logger.Instance.Debug("Booting HAT")

	Execute("clear", dto.RegexpDto{
		Patterns: []*regexp.Regexp{regexp.MustCompile(`(?i)^BHBL>(\s)?(version|clear)?$`)},
	})

	// prepare payload
	firmwarePayload := append([]byte{0x02}, append(firmwareBinary, []byte{0x03, '\r'}...)...)

	// suspend reading until full payload is sent. reading will resume automatically after write completes
	hatConnection.pauseRead()

	// send firmware size and checksum
	Execute(fmt.Sprintf("load %d %d", len(firmwareBinary), firmwareChecksum), dto.VoidDto{})

	// and finally send the payload
	Execute(firmwarePayload, dto.RegexpDto{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)^Image Received$"),
			regexp.MustCompile("(?i)^Checksum OK$"),
		},
	})

	// send signature
	writeFirmwareSignature()

	// and finally reboot
	Execute("reboot", dto.RegexpDto{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)^Done initialising ports$"),
		},
	})

	// sleep half of second to allow HAT to reboot
	time.Sleep(time.Microsecond * 500)
	logger.Instance.Debug("Firmware loaded successfully")
}

func writeFirmwareSignature() {
	signaturePath, _ := filepath.Abs("./data/signature.bin")
	signatureBinary, _ := os.ReadFile(signaturePath)

	// send signature size
	Execute(
		[]byte(fmt.Sprintf("signature %d\r", len(signatureBinary))),
		dto.VoidDto{},
	)

	// prepare payload
	signaturePayload := append([]byte{0x02}, append(signatureBinary, []byte{0x03, '\r'}...)...)

	// and finally send the payload
	Execute(signaturePayload, dto.RegexpDto{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)^Signature received$"),
		},
	})
}
