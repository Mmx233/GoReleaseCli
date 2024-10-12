package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

var Version = "-.-.-"

func main() {
	kingpin.MustParse(global.NewCommands(Version).Parse(os.Args[1:]))
	builder.Run()
}
