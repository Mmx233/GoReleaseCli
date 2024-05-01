package global

import (
	"os"
	"strconv"
)

// Init call after Config is set.
func Init() error {
	perm, err := strconv.ParseUint(Config.Perm, 8, 32)
	if err != nil {
		return err
	}
	Perm = os.FileMode(perm)
	return nil
}

var Config _Config

var Perm os.FileMode

type _Config struct {
	Target    string `env:"TARGET,required,notEmpty"`
	Ldflags   string `env:"LDFLAGS"`
	SoftFloat bool   `env:"SOFT-FLOAT"`
	Cgo       bool   `env:"CGO"`
	OS        string `env:"OS"`
	Arch      string `env:"ARCH"`
	Output    string `env:"OUTPUT,notEmpty" envDefault:"build"`
	Name      string `env:"NAME"`

	Compress              string `env:"COMPRESS"`
	Thread                uint16 `env:"TREAD"`
	DisableDefaultLdflags bool   `env:"DISABLE-DEFAULT-LDFLAGS"`
	Perm                  string `env:"PERM,notEmpty" envDefault:"0777"`
}
