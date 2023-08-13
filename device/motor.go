package device

type Motor struct {
	BaseDevice
}

func test() Device {
	x := &Motor{
		BaseDevice{
			ModelId: 0,
		},
	}

	return x
}
