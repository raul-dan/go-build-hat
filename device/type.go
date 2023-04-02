package device

const (
	PassiveMotor        string = "PassiveMotor"
	Light                      = "Light"
	TiltSensor                 = "TiltSensor"
	MotionSensor               = "MotionSensor"
	ColorDistanceSensor        = "ColorDistanceSensor"
	ColorSensor                = "ColorSensor"
	DistanceSensor             = "DistanceSensor"
	ForceSensor                = "ForceSensor"
	Matrix                     = "Matrix"
	Motor                      = "Motor"
)

type Type struct {
	TypeId uint32
	Type   string
	Name   string
}

var Types = map[uint32]Type{
	0x1:  {TypeId: 0x1, Type: PassiveMotor, Name: "PassiveMotor"},
	0x2:  {TypeId: 0x2, Type: PassiveMotor, Name: "PassiveMotor"},
	0x8:  {TypeId: 0x8, Type: Light, Name: "Light"},                                  // 88005
	0x22: {TypeId: 0x22, Type: TiltSensor, Name: "WeDo 2.0 Tilt Sensor"},             // 45305
	0x23: {TypeId: 0x23, Type: MotionSensor, Name: "MotionSensor"},                   // 45304
	0x25: {TypeId: 0x25, Type: ColorDistanceSensor, Name: "Color & Distance Sensor"}, // 88007
	0x3D: {TypeId: 0x3D, Type: ColorSensor, Name: "Color Sensor"},                    // 45605
	0x3E: {TypeId: 0x3E, Type: DistanceSensor, Name: "Distance Sensor"},              // 45604
	0x3F: {TypeId: 0x3F, Type: ForceSensor, Name: "Force Sensor"},                    // 45606
	0x40: {TypeId: 0x40, Type: Matrix, Name: "3x3 Color Light Matrix"},               // 45608
	0x26: {TypeId: 0x26, Type: Motor, Name: "Medium Linear Motor"},                   // 88008
	0x2E: {TypeId: 0x2E, Type: Motor, Name: "Large Motor"},                           // 88013
	0x2F: {TypeId: 0x2F, Type: Motor, Name: "XL Motor"},                              // 88014
	0x30: {TypeId: 0x30, Type: Motor, Name: "Medium Angular Motor (Cyan)"},           // 45603
	0x31: {TypeId: 0x31, Type: Motor, Name: "Large Angular Motor (Cyan)"},            // 45602
	0x41: {TypeId: 0x41, Type: Motor, Name: "Small Angular Motor"},                   // 45607
	0x4B: {TypeId: 0x4B, Type: Motor, Name: "Medium Angular Motor (Grey)"},           // 88018
	0x4C: {TypeId: 0x4C, Type: Motor, Name: "Large Angular Motor (Grey)"},            // 88017
}
