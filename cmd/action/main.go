package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
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
