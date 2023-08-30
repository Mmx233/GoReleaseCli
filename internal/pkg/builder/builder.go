package builder

import (
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	"github.com/Mmx233/GoReleaseCli/tools"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"sync"
)

func NewBuilder(outputDir string) (*Builder, error) {
	outputName := LoadBinaryName()

	builder := goCMD.NewBuilder(global.Commands.Target)
	builder = builder.ProductionLdflags().TrimPath()
	if global.Commands.Ldflags != "" {
		builder = builder.Ldflags(global.Commands.Ldflags)
	}

	arch, e := MatchTargetArch()
	if e != nil {
		return nil, e
	}

	if e = PrepareDirs(outputDir); e != nil {
		return nil, e
	}

	return &Builder{
		OutputName: outputName,
		OutputDir:  outputDir,
		Arch:       arch,
		GoCMD:      builder,
	}, nil
}

type Builder struct {
	OutputName string
	OutputDir  string
	Arch       map[string][]string
	GoCMD      goCMD.BuildCommand

	WaitGroup *sync.WaitGroup
	// 编译并发限制器
	TreadChan chan bool
	// 失败编译收集
	FailedArchChan chan string
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
	env = append(env, "GOOS="+GOOS, "GOARCH="+GOARCH)
	cmd.Env = append(cmd.Environ(), env...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return buildName, fmt.Errorf("build error: %v", err)
	}

	if err = tools.MakeZip(outputPath, a.OutputName); err != nil {
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

func (a *Builder) BuildArches() {
	// 准备编译携程
	a.TreadChan = make(chan bool, int(global.Commands.Thread))
	a.WaitGroup = &sync.WaitGroup{}
	var count uint
	for GOOS, Arches := range a.Arch {
		for _, GOARCH := range Arches {
			a.NewBuildThread(GOOS, GOARCH)
			count++
			if global.Commands.SoftFloat && strings.Contains(GOARCH, "mips") {
				a.NewBuildThread(GOOS, GOARCH, "GOMIPS=softfloat")
				count++
			}
		}
	}
	a.FailedArchChan = make(chan string, count)

	log.Infof("found %d arches, building...", count)

	// 开始编译
	for i := uint16(0); i < global.Commands.Thread; i++ {
		a.TreadChan <- true
	}
	a.WaitGroup.Wait()

	// 打印编译结果
	if len(a.FailedArchChan) == 0 {
		log.Infoln("build completed successfully")
	} else {
		log.Infof("build %d arches succeed, %d arches failed", count-uint(len(a.FailedArchChan)), len(a.FailedArchChan))
		failedArches := make([]string, len(a.FailedArchChan))
		i := 0
		for name := range a.FailedArchChan {
			failedArches[i] = name
			i++
		}
		log.Infof("failed arches: %s", strings.Join(failedArches, ", "))
	}
}
