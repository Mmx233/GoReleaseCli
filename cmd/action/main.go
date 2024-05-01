package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func init() {
	if err := env.Parse(&global.Config, env.Options{
		Prefix: "INPUT_",
	}); err != nil {
		log.Fatalln(err)
	}

	if global.Config.Output == "" {
		global.Config.Output = "build"
	}
	if global.Config.Thread == 0 {
		global.Config.Thread = uint16(runtime.NumCPU() + 1)
	}
}

func main() {
	builder.Run()
}
