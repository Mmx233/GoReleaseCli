package main

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/internal/pkg/builder"
	"github.com/Mmx233/GoReleaseCli/tools"
	"github.com/alecthomas/kingpin/v2"
	"os"
)

var Version = "-.-.-"

func main() {
	tools.MustSevenZip()
	global.InitCommands(Version)
	kingpin.MustParse(global.Commands.App.Parse(os.Args[1:]))
	builder.Run()
}
