package global

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"runtime"
)

func NewCommands(version string) *kingpin.Application {
	app := kingpin.New("release", "Golang build production release helper.")
	app.Version(version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	app.Flag("thread", "How many threads to use for parallel compilation.").Short('j').Default(fmt.Sprint(runtime.NumCPU() + 1)).Uint16Var(&Config.Thread)
	app.Flag("compress", "Compress the binary into the specified format of compressed file.").Short('c').HintOptions("zip", "tar.gz").StringVar(&Config.Compress)
	app.Flag("disable-default-ldflags", "Disable ldflags added by default.").BoolVar(&Config.DisableDefaultLdflags)

	app.Flag("ldflags", "Add custom ldflags.").StringVar(&Config.Ldflags)
	app.Flag("soft-float", "Enable soft float for mips.").BoolVar(&Config.SoftFloat)
	app.Flag("cgo", "Enable go cgo.").BoolVar(&Config.Cgo)
	app.Flag("os", "Target os").HintOptions("windows,linux").StringVar(&Config.OS)
	app.Flag("arch", "Target arch.").HintOptions("386,amd64").StringVar(&Config.Arch)

	app.Flag("output", "Output dir path.").Short('d').Default("build").StringVar(&Config.Output)
	app.Flag("name", "Output binary file name.").Short('o').StringVar(&Config.Name)
	app.Arg("target", "Target package.").Required().StringVar(&Config.Target)

	return app
}
