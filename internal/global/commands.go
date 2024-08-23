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
	app.Flag("perm", "Output file mode.").Default("0600").StringVar(&Config.Perm)

	app.Flag("mod-download-args", "custom args for go mod download.").HintOptions("-x").StringVar(&Config.ModDownloadArgs)

	app.Flag("ldflags", "Add custom ldflags.").StringVar(&Config.Ldflags)
	app.Flag("cgo", "Enable go cgo.").BoolVar(&Config.Cgo)
	app.Flag("os", "Target os.").HintOptions("windows,linux").StringVar(&Config.OS)
	app.Flag("arch", "Target arch.").HintOptions("386,amd64").StringVar(&Config.Arch)
	app.Flag("platforms", "Specify platforms").HintOptions("linux/386,windows/arm64").StringVar(&Config.Platforms)
	app.Flag("extra-arches", "Build all extra arches.").BoolVar(&Config.ExtraArches)
	app.Flag("extra-arches-show-default", "Show default extra arch name.").BoolVar(&Config.ExtraArchesShowDefault)

	app.Flag("output", "Output dir path.").Short('d').Default("build").StringVar(&Config.Output)
	app.Flag("name", "Output binary file name.").Short('o').StringVar(&Config.Name)
	app.Arg("target", "Target package.").Required().StringVar(&Config.Target)

	return app
}
