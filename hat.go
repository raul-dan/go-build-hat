package buildhat

var hat *IHat

type IHat struct {
}

func (h *IHat) GetFirmwareVersion() string {
	return sExec("version").(string)
}

func (h *IHat) GetVin() float64 {
	return sExec("vin").(float64)
}

func Hat() *IHat {
	if hat == nil {
		hat = new(IHat)
	}

	return hat
}
