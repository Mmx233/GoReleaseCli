package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := env.Parse(&global.Config, env.Options{
		Prefix: "INPUT_",
	}); err != nil {
		log.Fatalln(err)
	}
	builder.Run()
}
