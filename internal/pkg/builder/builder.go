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
	e := cmd.Wait()
	if e != nil {
		return buildName, fmt.Errorf("build error: %v", e)
	}

	if e = tools.MakeZip(outputPath, a.OutputName); e != nil {
		return buildName, fmt.Errorf("compress %s failed: %v", outputPath, e)
	}

	return buildName, nil
}

func (a *Builder) NewBuildThread(wg *sync.WaitGroup, thread chan bool, GOOS, GOARCH string, env ...string) {
	wg.Add(1)
	go func() {
		<-thread
		if name, e := a.Build(GOOS, GOARCH, env...); e != nil {
			log.Errorf("error occur while building %s: %v", name, e)
		} else {
			log.Infof("build %s success", name)
		}
		wg.Done()
		thread <- true
	}()
}

func (a *Builder) BuildArches() {
	var threadChan = make(chan bool, int(global.Commands.Thread))
	var wg = &sync.WaitGroup{}
	var count uint
	for GOOS, Arches := range a.Arch {
		for _, GOARCH := range Arches {
			a.NewBuildThread(wg, threadChan, GOOS, GOARCH)
			count++
			if global.Commands.SoftFloat && strings.Contains(GOARCH, "mips") {
				a.NewBuildThread(wg, threadChan, GOOS, GOARCH, "GOMIPS=softfloat")
				count++
			}
		}
	}
	log.Infof("found %d arches, bulding...", count)

	for i := uint16(0); i < global.Commands.Thread; i++ {
		threadChan <- true
	}
	wg.Wait()
}
