package buildhat

import "buildhat/device"

var hat *device.Hat

func Hat() *device.Hat {
	if hat == nil {
		hat = &device.Hat{}
	}

	return hat
}
