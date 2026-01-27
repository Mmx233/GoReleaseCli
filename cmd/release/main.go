package main

import (
	"context"
	"os"

	"github.com/Mmx233/GoReleaseCli/internal/builder"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/tools"
	"github.com/alecthomas/kingpin/v2"
)

var Version = "-.-.-"

func main() {
	kingpin.MustParse(global.NewCommands(Version).Parse(os.Args[1:]))
	ctx, cancel := context.WithCancel(context.Background())
	go tools.OsSignalCancel(cancel)
	builder.Run(ctx)
}
