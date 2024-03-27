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

	Commands.App.Flag("thread", "How many threads to use for parallel compilation.").Short('j').Default(fmt.Sprint(runtime.NumCPU() + 1)).Uint16Var(&Commands.Thread)
	Commands.App.Flag("compress", "Compress the binary into the specified format of compressed file.").Short('c').HintOptions("zip", "tar.gz").StringVar(&Commands.Compress)
	Commands.App.Flag("disable-default-ldflags", "Disable ldflags added by default.").BoolVar(&Commands.DisableDefaultLdflags)

	Commands.App.Flag("ldflags", "Add custom ldflags.").StringVar(&Commands.Ldflags)
	Commands.App.Flag("soft-float", "Enable soft float for mips.").BoolVar(&Commands.SoftFloat)
	Commands.App.Flag("cgo", "Enable go cgo.").BoolVar(&Commands.Cgo)
	Commands.App.Flag("os", "Target os").HintOptions("windows,linux").StringVar(&Commands.OS)
	Commands.App.Flag("arch", "Target arch.").HintOptions("386,amd64").StringVar(&Commands.Arch)

	Commands.App.Flag("output", "Output dir path.").Short('d').Default("build").StringVar(&Commands.Output)
	Commands.App.Flag("name", "Output binary file name.").Short('o').StringVar(&Commands.Name)
	Commands.App.Arg("target", "Target package.").Required().StringVar(&Commands.Target)
}
