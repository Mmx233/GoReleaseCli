package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
	"os"
)

var Version = "-.-.-"

func main() {
	kingpin.MustParse(global.NewCommands(Version).Parse(os.Args[1:]))
	if err := global.Init(); err != nil {
		log.Fatalln(err)
	}
	builder.Run()
}
