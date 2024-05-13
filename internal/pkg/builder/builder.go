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

func (a *Builder) Build(GOOS, GOARCH string, env ...string) (string, error) {
	args := make([]string, len(env)+2)
	args[0], args[1] = GOOS, GOARCH
	i := 2
	for _, v := range env {
		args[i] = strings.Split(v, "=")[1]
		i++
	}
	outputName := a.OutputName
	if GOOS == "windows" {
		outputName += ".exe"
	}
	buildName := BuildName(outputName, args...)
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

func (a *Builder) NewBuildThread(GOOS, GOARCH string, env ...string) {
	a.WaitGroup.Add(1)
	go func() {
		<-a.TreadChan
		name, err := a.Build(GOOS, GOARCH, env...)
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
	// 准备编译携程
	a.TreadChan = make(chan bool, int(global.Config.Thread))
	a.WaitGroup = &sync.WaitGroup{}
	var count int
	for GOOS, Arches := range a.Arch {
		for _, GOARCH := range Arches {
			a.NewBuildThread(GOOS, GOARCH)
			count++
			if global.Config.SoftFloat && strings.Contains(GOARCH, "mips") {
				a.NewBuildThread(GOOS, GOARCH, "GOMIPS=softfloat")
				count++
			}
		}
	}
	a.FailedArchChan = make(chan string, count)

	log.Infof("found %d arches, building...", count)

	// 开始编译
	for i := uint16(0); i < global.Config.Thread; i++ {
		a.TreadChan <- true
	}
	a.WaitGroup.Wait()

	// 打印编译结果
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
