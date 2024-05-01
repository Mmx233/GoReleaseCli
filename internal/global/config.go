package global

var Config _Config

type _Config struct {
	Target    string `env:"TARGET,required"`
	Ldflags   string `env:"LDFLAGS"`
	SoftFloat bool   `env:"SOFT-FLOAT"`
	Cgo       bool   `env:"CGO"`
	OS        string `env:"OS"`
	Arch      string `env:"ARCH"`
	Output    string `env:"OUTPUT"`
	Name      string `env:"NAME"`

	Compress              string `env:"COMPRESS"`
	Thread                uint16 `env:"TREAD"`
	DisableDefaultLdflags bool   `env:"DISABLE-DEFAULT-LDFLAGS"`
}
