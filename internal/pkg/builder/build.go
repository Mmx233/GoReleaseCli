package builder

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Run() {
	PrepareDirs()

	var targetOS []string
	var targetArch []string
	if global.Commands.OS != "" {
		targetOS = strings.Split(global.Commands.OS, ",")
	}
	if global.Commands.Arch != "" {
		targetArch = strings.Split(global.Commands.Arch, ",")
	}

	binaryName := LoadBinaryName()

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
	runBuild := func(binaryName, GOOS, GOARCH string, env ...string) {
		if env != nil {
			log.Infof("building %s/%s %v", GOOS, GOARCH, env)
		} else {
			log.Infof("building %s/%s", GOOS, GOARCH)
		}

		cmd := builder.Exec()
		env = append(env, "GOOS="+GOOS, "GOARCH="+GOARCH)
		cmd.Env = append(cmd.Environ(), env...)
		output, e := cmd.Output()
		if e != nil {
			log.Errorf("build error: %v: %s", e, string(output))
			return
		}

		args := make([]string, len(env))
		args[0], args[1] = GOOS, GOARCH
		i := 2
		for _, v := range env[2:] {
			args[i] = strings.Split(v, "=")[1]
			i++
		}
		if binaryName, e = RenameBinary(binaryName, args...); e != nil {
			log.Fatalln(e)
		}
		log.Infof("build %s success", binaryName)

		if e = MakeZip(binaryName); e != nil {
			log.Fatalln("compress failed:", e)
		}
	}
	for GOOS, Arches := range arch {
		var binaryName = binaryName
		if GOOS == "windows" {
			binaryName += ".exe"
		}
		for _, GOARCH := range Arches {
			runBuild(binaryName, GOOS, GOARCH)
			if global.Commands.SoftFloat && strings.Contains(GOARCH, "mips") {
				runBuild(binaryName, GOOS, GOARCH, "GOMIPS=softfloat")
			}
		}
	}
}

func PrepareDirs() {
	_ = os.RemoveAll("build")
	_ = os.Mkdir("build", 0600)
}

func LoadBinaryName() string {
	target := strings.Replace(global.Commands.Target, `\`, "/", -1)
	target = strings.TrimSuffix(target, "/")
	temp := strings.Split(target, "/")
	return temp[len(temp)-1]
}

func RenameBinary(binaryName string, suffix ...string) (string, error) {
	ext := path.Ext(binaryName)
	name := "build/" + strings.TrimSuffix(binaryName, ext)
	for _, s := range suffix {
		name += "_" + s
	}
	name += ext
	return name, os.Rename(binaryName, name)
}

func MakeZip(binaryName string) error {
	name := strings.TrimSuffix(binaryName, path.Ext(binaryName))
	output, e := exec.Command("7z", "a", name+".zip", binaryName).Output()
	if e != nil {
		log.Debugln(string(output))
		return e
	}
	return nil
}
