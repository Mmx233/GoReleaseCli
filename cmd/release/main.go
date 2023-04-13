package main

import (
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

const Version = "-.-.-"

func main() {
	global.InitCommands(Version)
	kingpin.MustParse(global.Commands.App.Parse(os.Args[1:]))

	var targetOS []string
	var targetArch []string
	if global.Commands.OS != "" {
		targetOS = strings.Split(global.Commands.OS, ",")
	}
	if global.Commands.Arch != "" {
		targetArch = strings.Split(global.Commands.Arch, ",")
	}

	builder := goCMD.NewBuilder(global.Commands.Target)
	builder = builder.ProductionLdflags().TrimPath().OutputName(global.Commands.OutputName)
	if global.Commands.Ldflags != "" {
		builder = builder.Ldflags(global.Commands.Ldflags)
	}

	var arch = make(map[string][]string, len(targetOS))

	// pair GOOS
	if len(targetOS) == 0 {
		arch = goCMD.Arch
	} else {
		for _, GOOS := range targetOS {
			if GOARCH, ok := goCMD.Arch[GOOS]; ok {
				arch[GOOS] = GOARCH
			}
		}
		if len(arch) == 0 {
			log.Fatalln("no valid os found")
		}
	}

	// pair GOARCH
	var keepArch = make(map[string]int, len(targetOS))
	if len(targetArch) != 0 {
		for GOOS, Arches := range arch {
			archCounter := 0
			for i, GOARCH := range Arches {
				for _, ArchEX := range targetArch {
					if GOARCH == ArchEX {
						archCounter++
						goto nextArch
					}
				}
				Arches[i] = ""
			nextArch:
			}
			keepArch[GOOS] = archCounter
		}
		for GOOS, count := range keepArch {
			if count == 0 {
				delete(arch, GOOS)
				continue
			}

			newARCH := make([]string, count)
			i := 0
			for _, Arch := range arch[GOOS] {
				if Arch != "" {
					newARCH[i] = Arch
					i++
				}
			}
			arch[GOOS] = newARCH
		}
		if len(arch) == 0 {
			log.Fatalln("no valid arch found")
		}
	}

	// print pair result
	var archOutput string
	for GOOS, Arches := range arch {
		for _, GOARCH := range Arches {
			archOutput += fmt.Sprintf("%s/%s ", GOOS, GOARCH)
		}
	}
	log.Infof("building platform: %s", archOutput)

	// build
	runBuild := func(env ...string) {
		cmd := builder.Exec()
		cmd.Env = append(cmd.Environ(), env...)
		output, e := cmd.Output()
		if e != nil {
			log.Fatalln("build error:", e)
		} else if len(output) != 0 {
			fmt.Println(string(output))
		}
	}
	for GOOS, Arches := range arch {
		for _, GOARCH := range Arches {
			runBuild("GOOS="+GOOS, "GOARCH="+GOARCH)

			if global.Commands.SoftFloat && strings.Contains(GOARCH, "mips") {
				runBuild("GOOS="+GOOS, "GOARCH="+GOARCH, "GOMIPS=soft-float")
			}
		}
	}
}
