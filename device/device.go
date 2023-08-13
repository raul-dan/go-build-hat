package device

type Type string

const (
	PassiveMotorDeviceType        Type = "PassiveMotor"
	LightDeviceType                    = "Light"
	TiltSensorDeviceType               = "TiltSensor"
	MotionSensorDeviceType             = "MotionSensor"
	ColorDistanceSensorDeviceType      = "ColorDistanceSensor"
	ColorSensorDeviceType              = "ColorSensor"
	DistanceSensorDeviceType           = "DistanceSensor"
	ForceSensorDeviceType              = "ForceSensor"
	MatrixDeviceType                   = "Matrix"
	MotorDeviceType                    = "Motor"
)

type knownDeviceType struct {
	Type Type
	Name string
}

type Device interface {
	Name() string
	Type() Type
}

type BaseDevice struct {
	ModelId uint32
}

func (d *BaseDevice) Name() string {
	device := getKnownDevice(d.ModelId)

	if device == nil {
		return "Unknown Device"
	}

	return device.Name
}

func (d *BaseDevice) Type() Type {
	device := getKnownDevice(d.ModelId)

	if device == nil {
		return "Unknown Device"
	}

	return device.Type
}

func getKnownDevice(modelId uint32) *knownDeviceType {
	knownDevice, ok := knownDevices[modelId]

	if !ok {
		return nil
	}

	return knownDevice
}

var knownDevices = map[uint32]*knownDeviceType{
	0x1:  {Type: PassiveMotorDeviceType, Name: "PassiveMotorDeviceType"},
	0x2:  {Type: PassiveMotorDeviceType, Name: "PassiveMotorDeviceType"},
	0x8:  {Type: LightDeviceType, Name: "LightDeviceType"},                             // 88005
	0x22: {Type: TiltSensorDeviceType, Name: "WeDo 2.0 Tilt Sensor"},                   // 45305
	0x23: {Type: MotionSensorDeviceType, Name: "MotionSensorDeviceType"},               // 45304
	0x25: {Type: ColorDistanceSensorDeviceType, Name: "Color & Distance Sensor"},       // 88007
	0x3D: {Type: ColorSensorDeviceType, Name: "Color Sensor"},                          // 45605
	0x3E: {Type: DistanceSensorDeviceType, Name: "Distance Sensor"},                    // 45604
	0x3F: {Type: ForceSensorDeviceType, Name: "Force Sensor"},                          // 45606
	0x40: {Type: MatrixDeviceType, Name: "3x3 Color LightDeviceType MatrixDeviceType"}, // 45608
	0x26: {Type: MotorDeviceType, Name: "Medium Linear MotorDeviceType"},               // 88008
	0x2E: {Type: MotorDeviceType, Name: "Large MotorDeviceType"},                       // 88013
	0x2F: {Type: MotorDeviceType, Name: "XL MotorDeviceType"},                          // 88014
	0x30: {Type: MotorDeviceType, Name: "Medium Angular MotorDeviceType (Cyan)"},       // 45603
	0x31: {Type: MotorDeviceType, Name: "Large Angular MotorDeviceType (Cyan)"},        // 45602
	0x41: {Type: MotorDeviceType, Name: "Small Angular MotorDeviceType"},               // 45607
	0x4B: {Type: MotorDeviceType, Name: "Medium Angular MotorDeviceType (Grey)"},       // 88018
	0x4C: {Type: MotorDeviceType, Name: "Large Angular MotorDeviceType (Grey)"},        // 88017
}
