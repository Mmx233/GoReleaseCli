package builder

import (
	"errors"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	"os"
	"path"
	"strings"
)

func PrepareDirs(outputDir string) error {
	_ = os.RemoveAll(outputDir)
	return os.MkdirAll(outputDir, 0600)
}

func LoadBinaryName() string {
	if global.Commands.Name != "" {
		return global.Commands.Name
	}
	target := strings.Replace(global.Commands.Target, `\`, "/", -1)
	target = strings.TrimSuffix(target, "/")
	temp := strings.Split(target, "/")
	return temp[len(temp)-1]
}

func BuildName(binaryName string, suffix ...string) string {
	ext := path.Ext(binaryName)
	name := strings.TrimSuffix(binaryName, ext)
	for _, s := range suffix {
		name += "_" + s
	}
	name += ext
	return name
}

// MatchTargetArch 由于没有进行深拷贝，本方法只能执行一次
func MatchTargetArch() (map[string][]string, error) {
	var targetOS []string
	var targetArch []string
	if global.Commands.OS != "" {
		targetOS = strings.Split(global.Commands.OS, ",")
	}
	if global.Commands.Arch != "" {
		targetArch = strings.Split(global.Commands.Arch, ",")
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
			return nil, errors.New("no valid os found")
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
			return nil, errors.New("no valid arch found")
		}
	}
	return arch, nil
}
