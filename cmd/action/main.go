package main

import (
	"context"
	"github.com/Mmx233/GoReleaseCli/internal/builder"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/tools"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := global.ParseConfigFromEnv(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go tools.OsSignalCancel(cancel)
	builder.Run(ctx)
}
