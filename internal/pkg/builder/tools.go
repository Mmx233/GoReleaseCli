package builder

import (
	"errors"
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	"os"
	"path"
	"strings"
)

func PrepareDirs(outputDir string) error {
	_ = os.RemoveAll(outputDir)
	return os.MkdirAll(outputDir, global.Perm)
}

func LoadBinaryName() string {
	if global.Config.Name != "" {
		return global.Config.Name
	}
	target := strings.Replace(global.Config.Target, `\`, "/", -1)
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

// MatchTargetArch can only be executed once because permanent modify of global arch map
func MatchTargetArch() (map[string][]string, error) {
	var arch map[string][]string

	if global.Config.Platforms == "" || global.Config.OS != "" || global.Config.Arch != "" {
		var targetOS []string
		var targetArch []string
		if global.Config.OS != "" {
			targetOS = strings.Split(global.Config.OS, ",")
		}
		if global.Config.Arch != "" {
			targetArch = strings.Split(global.Config.Arch, ",")
		}

		arch = make(map[string][]string, len(targetOS))

		// match GOOS
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

		// match GOARCH
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
		}
	}

	// add platforms
	if global.Config.Platforms != "" {
		if arch == nil {
			arch = make(map[string][]string)
		}
		platforms := strings.Split(global.Config.Platforms, ",")
		for _, platform := range platforms {
			splitPlatform := strings.Split(platform, "/")
			if len(splitPlatform) != 2 {
				return nil, fmt.Errorf("invalid platform: %s", platform)
			}
			platformGOOS, platformGOARCH := splitPlatform[0], splitPlatform[1]
			Arches, _ := arch[platformGOOS]
			for _, GOARCH := range Arches {
				if GOARCH == platformGOARCH {
					goto nextPlatform
				}
			}
			arch[platformGOOS] = append(Arches, platformGOARCH)
		nextPlatform:
		}
	}

	if len(arch) == 0 {
		return nil, errors.New("no valid arch found")
	}

	return arch, nil
}
