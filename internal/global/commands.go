package global

import (
	"github.com/Mmx233/GoReleaseCli/internal/global/models"
	"github.com/alecthomas/kingpin/v2"
)

var Commands models.Commands

func InitCommands(version string) {
	Commands.App = kingpin.New("release", "Golang build production release helper.")
	Commands.App.Version(version)
	Commands.App.VersionFlag.Short('v')
	Commands.App.HelpFlag.Short('h')

	Commands.App.Flag("ldflags", "add custom ldflags").StringVar(&Commands.Ldflags)
	Commands.App.Flag("soft-float", "enable soft float for mips").BoolVar(&Commands.SoftFloat)
	Commands.App.Flag("output", "output name").Short('o').Default("main").StringVar(&Commands.OutputName)
	Commands.App.Flag("os", "target os").HintOptions("windows,linux").StringVar(&Commands.OS)
	Commands.App.Flag("arch", "target arch").HintOptions("386,amd64").StringVar(&Commands.Arch)
	Commands.App.Arg("target", "target package").Default(".").StringVar(&Commands.Target)
}
