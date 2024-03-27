package global

import (
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global/models"
	"github.com/alecthomas/kingpin/v2"
	"runtime"
)

var Commands models.Commands

func InitCommands(version string) {
	Commands.App = kingpin.New("release", "Golang build production release helper.")
	Commands.App.Version(version)
	Commands.App.VersionFlag.Short('v')
	Commands.App.HelpFlag.Short('h')

	Commands.App.Flag("thread", "use how many thread to build").Short('j').Default(fmt.Sprint(runtime.NumCPU() + 1)).Uint16Var(&Commands.Thread)
	Commands.App.Flag("compress", "compress binary").Short('c').HintOptions("zip", "tar.gz").StringVar(&Commands.Compress)

	Commands.App.Flag("ldflags", "add custom ldflags").StringVar(&Commands.Ldflags)
	Commands.App.Flag("soft-float", "enable soft float for mips").BoolVar(&Commands.SoftFloat)
	Commands.App.Flag("cgo", "enable cgo").BoolVar(&Commands.Cgo)
	Commands.App.Flag("os", "target os").HintOptions("windows,linux").StringVar(&Commands.OS)
	Commands.App.Flag("arch", "target arch").HintOptions("386,amd64").StringVar(&Commands.Arch)

	Commands.App.Flag("output", "output dir path").Short('d').Default("build").StringVar(&Commands.Output)
	Commands.App.Flag("name", "output binary file name").Short('o').StringVar(&Commands.Name)
	Commands.App.Arg("target", "target package").Required().StringVar(&Commands.Target)

}
