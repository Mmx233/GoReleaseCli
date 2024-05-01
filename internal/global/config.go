package global

var Config _Config

type _Config struct {
	Target    string
	Ldflags   string
	SoftFloat bool
	Cgo       bool
	OS        string
	Arch      string
	Output    string
	Name      string

	Compress              string
	Thread                uint16
	DisableDefaultLdflags bool
}
