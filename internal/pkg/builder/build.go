package builder

import (
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Run() {
	var targetOS []string
	var targetArch []string
	if global.Commands.OS != "" {
		targetOS = strings.Split(global.Commands.OS, ",")
	}
	if global.Commands.Arch != "" {
		targetArch = strings.Split(global.Commands.Arch, ",")
	}

	builder := goCMD.NewBuilder(global.Commands.Target)
	builder = builder.ProductionLdflags().TrimPath()
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

	// build
	runBuild := func(GOOS, GOARCH string, env ...string) {
		log.Infof("building %s/%s", GOOS, GOARCH)

		cmd := builder.Exec()
		env = append(env, "GOOS="+GOOS, "GOARCH="+GOARCH)
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
			runBuild(GOOS, GOARCH)

			if global.Commands.SoftFloat && strings.Contains(GOARCH, "mips") {
				runBuild(GOOS, GOARCH, "GOMIPS=soft-float")
			}
		}
	}
}
