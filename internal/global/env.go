package global

import (
	"github.com/caarlos0/env/v11"
	"runtime"
)

func ParseConfigFromEnv() error {
	if err := env.ParseWithOptions(&Config, env.Options{
		Prefix: "INPUT_",
	}); err != nil {
		return err
	}

	if Config.Thread == 0 {
		Config.Thread = uint16(runtime.NumCPU() + 1)
	}
	return nil
}
