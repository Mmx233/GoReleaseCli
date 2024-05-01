package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"path"
	"runtime"
)

func init() {
	if err := env.Parse(&global.Config, env.Options{
		Prefix: "INPUT_",
	}); err != nil {
		log.Fatalln(err)
	}
	if global.Config.Thread == 0 {
		global.Config.Thread = uint16(runtime.NumCPU() + 1)
	}
	global.Config.Output = path.Join("/github/workspace", global.Config.Output)
	if err := global.Init(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	builder.Run()
}
