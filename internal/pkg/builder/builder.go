package builder

import (
	"errors"
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/compress"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"sync"
)

func NewBuilder(outputDir string) (*Builder, error) {
	outputName := LoadBinaryName()

	goBuilder := goCMD.NewBuilder(global.Config.Target).TrimPath()
	if !global.Config.DisableDefaultLdflags {
		goBuilder = goBuilder.ProductionLdflags()
	}
	if global.Config.Ldflags != "" {
		goBuilder = goBuilder.Ldflags(global.Config.Ldflags)
	}

	arch, err := MatchTargetArch()
	if err != nil {
		return nil, err
	}

	builder := &Builder{
		OutputName: outputName,
		OutputDir:  outputDir,
		Arch:       arch,
		GoCMD:      goBuilder,
	}

	if global.Config.Compress == "" {
		builder.Compress = func(_, _ string) error {
			return nil
		}
	} else {
		var compressor compress.Make
		if compress.SevenZipAvailable() {
			switch global.Config.Compress {
			case "zip":
				compressor = compress.SevenZipMakeZip
			case "tar.gz":
				compressor = compress.SevenZipMakeTarGzip
			}
		} else {
			switch global.Config.Compress {
			case "zip":
				if compress.ZipAvailable() {
					compressor = compress.ZipMakeZip
				}
			case "tar.gz":
				if compress.TarAvailable() {
					compressor = compress.TarMakeTarGzip
				}
			}
		}
		if compressor == nil {
			log.Fatalf("compression library is missing or the compression format (%s) is not supported", global.Config.Compress)
		}
		builder.Compress = compressor
	}

	if global.Config.Cgo {
		builder.Cgo = "CGO_ENABLED=1"
	} else {
		builder.Cgo = "CGO_ENABLED=0"
	}

	if err = PrepareDirs(outputDir); err != nil {
		return nil, err
	}

	return builder, nil
}

type Builder struct {
	OutputName string
	OutputDir  string
	Cgo        string
	Arch       map[string][]string
	GoCMD      goCMD.BuildCommand

	WaitGroup *sync.WaitGroup
	// 编译并发限制器
	TreadChan chan bool
	// 失败编译收集
	FailedArchChan chan string
	Compress       compress.Make
}

func (a *Builder) Build(GOOS, GOARCH string, nameSuffix, env []string) (string, error) {
	outputName := a.OutputName
	if GOOS == "windows" {
		outputName += ".exe"
	}
	buildName := BuildName(outputName, append([]string{GOOS, GOARCH}, nameSuffix...)...)
	outputPath := path.Join(a.OutputDir, buildName)

	cmd := a.GoCMD.OutputName(outputPath).Exec()
	cmd.Env = append(cmd.Environ(), env...)
	cmd.Env = append(cmd.Env, a.Cgo, "GOOS="+GOOS, "GOARCH="+GOARCH)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return buildName, fmt.Errorf("build error: %v", err)
	}

	if err = a.Compress(outputPath, a.OutputName); err != nil {
		return buildName, fmt.Errorf("compress %s failed: %v", outputPath, err)
	}

	return buildName, nil
}

func (a *Builder) NewBuildThread(GOOS, GOARCH string, nameSuffix, env []string) {
	a.WaitGroup.Add(1)
	go func() {
		<-a.TreadChan
		name, err := a.Build(GOOS, GOARCH, nameSuffix, env)
		if err != nil {
			log.Errorf("error occur while building %s: %v", name, err)
		} else {
			log.Infof("build %s success", name)
		}
		a.WaitGroup.Done()
		a.TreadChan <- true

		// 报告编译错误架构
		if err != nil {
			a.FailedArchChan <- name
		}
	}()
}

func (a *Builder) BuildArches() error {
	// prepare compile channels and goroutines
	a.TreadChan = make(chan bool, int(global.Config.Thread))
	a.WaitGroup = &sync.WaitGroup{}
	var count int
	for GOOS, Arches := range a.Arch {
		for _, GOARCH := range Arches {
			if global.Config.ExtraArches {
				extraArches, ok := goCMD.ExtraArches[GOARCH]
				if !ok {
					a.NewBuildThread(GOOS, GOARCH, nil, nil)
					count++
					continue
				}
				for _, extraArch := range extraArches {
					for _, value := range extraArch.Values {
						env := fmt.Sprintf("%s=%s", extraArch.EnvKey, value.Value)
						a.NewBuildThread(GOOS, GOARCH, value.Names(global.Config.ExtraArchesShowDefault), []string{env})
						count++
						if value.ExtraFloat != "" {
							a.NewBuildThread(GOOS, GOARCH, value.NamesExtraFloat(global.Config.ExtraArchesShowDefault), []string{env + "," + value.ExtraFloat})
							count++
						}
					}
				}
			}
		}
	}
	a.FailedArchChan = make(chan string, count)

	log.Infof("found %d arches, building...", count)

	// start compile
	for i := uint16(0); i < global.Config.Thread; i++ {
		a.TreadChan <- true
	}
	a.WaitGroup.Wait()

	// print compile result
	if len(a.FailedArchChan) == 0 {
		log.Infoln("completed successfully")
	} else if len(a.FailedArchChan) == count {
		return errors.New("all arches build failed")
	} else {
		log.Infof("completed: %d arches succeed, %d arches failed", count-len(a.FailedArchChan), len(a.FailedArchChan))
		failedArches := make([]string, len(a.FailedArchChan))
		for i := len(a.FailedArchChan) - 1; i >= 0; i-- {
			failedArches[i] = <-a.FailedArchChan
		}
		log.Infof("failed arches: %s", strings.Join(failedArches, ", "))
	}
	return nil
}
