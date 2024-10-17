package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/builder"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := global.ParseConfigFromEnv(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	builder.Run()
}
