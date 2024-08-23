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

// MatchTargetPlatforms can only be executed once because permanent modify of global arch map
func MatchTargetPlatforms() (map[string][]string, error) {
	var targetPlatforms map[string][]string

	if global.Config.Platforms == "" || global.Config.OS != "" || global.Config.Arch != "" {
		var targetOS []string
		var targetArch []string
		if global.Config.OS != "" {
			targetOS = strings.Split(global.Config.OS, ",")
		}
		if global.Config.Arch != "" {
			targetArch = strings.Split(global.Config.Arch, ",")
		}

		targetPlatforms = make(map[string][]string, len(targetOS))

		// match GOOS
		if len(targetOS) == 0 {
			targetPlatforms = goCMD.Platforms
		} else {
			for _, GOOS := range targetOS {
				if GOARCH, ok := goCMD.Platforms[GOOS]; ok {
					targetPlatforms[GOOS] = GOARCH
				}
			}
			if len(targetPlatforms) == 0 {
				return nil, errors.New("no valid os found")
			}
		}

		// match GOARCH
		var keepArch = make(map[string]int, len(targetOS))
		if len(targetArch) != 0 {
			for GOOS, Arches := range targetPlatforms {
				archCounter := 0
				for i, GOARCH := range Arches {
					for _, GOARCHExist := range targetArch {
						if GOARCH == GOARCHExist {
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
					delete(targetPlatforms, GOOS)
					continue
				}

				newArches := make([]string, count)
				i := 0
				for _, Arch := range targetPlatforms[GOOS] {
					if Arch != "" {
						newArches[i] = Arch
						i++
					}
				}
				targetPlatforms[GOOS] = newArches
			}
		}
	}

	// add platforms
	if global.Config.Platforms != "" {
		if targetPlatforms == nil {
			targetPlatforms = make(map[string][]string)
		}
		platforms := strings.Split(global.Config.Platforms, ",")
		for _, platform := range platforms {
			splitPlatform := strings.Split(platform, "/")
			if len(splitPlatform) != 2 {
				return nil, fmt.Errorf("invalid platform: %s", platform)
			}
			platformGOOS, platformGOARCH := splitPlatform[0], splitPlatform[1]
			Arches, _ := targetPlatforms[platformGOOS]
			for _, GOARCH := range Arches {
				if GOARCH == platformGOARCH {
					goto nextPlatform
				}
			}
			targetPlatforms[platformGOOS] = append(Arches, platformGOARCH)
		nextPlatform:
		}
	}

	if len(targetPlatforms) == 0 {
		return nil, errors.New("no valid platform found")
	}

	return targetPlatforms, nil
}
