package builder

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/compress"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

type Config struct {
	OutputDir       string
	OutputFormatter OutputFormatter
}

func NewBuilder(conf Config) (*Builder, error) {
	outputName := LoadBinaryName()

	goBuilder := goCMD.NewBuilder(global.Config.Target).TrimPath()
	if !global.Config.DisableDefaultLdflags {
		goBuilder = goBuilder.ProductionLdflags()
	}
	if global.Config.Ldflags != "" {
		goBuilder = goBuilder.Ldflags(global.Config.Ldflags)
	}

	platforms, err := MatchTargetPlatforms()
	if err != nil {
		return nil, err
	}

	builder := &Builder{
		OutputName:      outputName,
		OutputFormatter: conf.OutputFormatter,
		OutputDir:       conf.OutputDir,
		Platforms:       platforms,
		GoCMD:           goBuilder,
	}
	if builder.OutputFormatter == nil {
		builder.OutputFormatter = DefaultOutputFormatter{}
	}

	if global.Config.Compress == "" {
		builder.Compress = func(_ context.Context, _, _ string) error {
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

	if err := PrepareDirs(conf.OutputDir); err != nil {
		return nil, err
	}

	return builder, nil
}

type Builder struct {
	OutputName      string
	OutputFormatter OutputFormatter
	OutputDir       string
	Cgo             string
	Platforms       map[string][]string
	GoCMD           goCMD.BuildCommand

	WaitGroup      *sync.WaitGroup
	TaskChan       chan *Task
	FailedTaskChan chan string
	Compress       compress.Make
}

type Task struct {
	Name string

	// status of whole builder
	Stat    *BuildStat
	builder *Builder

	ctx        context.Context
	cmd        *exec.Cmd
	outputPath string
}

func (task Task) Build() error {
	defer task.Stat.Done()

	err := task.cmd.Run()
	if err != nil {
		return fmt.Errorf("build error: %v", err)
	}
	if err = task.builder.Compress(task.ctx, task.outputPath, task.builder.OutputName); err != nil {
		return fmt.Errorf("compress %s failed: %v", task.outputPath, err)
	}
	return nil
}

func (b *Builder) NewTask(ctx TaskContext) *Task {
	outputName := b.OutputName
	if ctx.GOOS == "windows" {
		outputName += ".exe"
	}
	buildName := b.OutputFormatter.Format(outputName, global.Config.Divider, append([]string{ctx.GOOS, ctx.GOARCH}, ctx.NameSuffix...)...)
	outputPath := path.Join(b.OutputDir, buildName)

	cmd := b.GoCMD.OutputName(outputPath).ExecContext(ctx)
	cmd.Env = append(cmd.Environ(), ctx.Env...)
	cmd.Env = append(cmd.Env, b.Cgo, "GOOS="+ctx.GOOS, "GOARCH="+ctx.GOARCH)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Task{
		ctx:        ctx,
		cmd:        cmd,
		builder:    b,
		outputPath: outputPath,
		Name:       buildName,
		Stat:       ctx.Stat,
	}
}

func (b *Builder) BuildThread() {
	for {
		task, ok := <-b.TaskChan
		if !ok {
			b.WaitGroup.Done()
			return
		}
		if task.ctx.Err() != nil {
			continue
		}
		err := task.Build()
		logger := log.WithField("process", task.Stat.PercentageString())
		if err != nil {
			logger.Errorf("error occur while building %s: %v", task.Name, err)
			b.FailedTaskChan <- task.Name
		} else {
			logger.Infof("build %s success", task.Name)
		}
	}
}

func (b *Builder) CalcExtraArches(ctx ArchContext, GOOS, GOARCH string, extraArches []goCMD.ArchExtra) *list.List {
	tasks := list.New()
	for _, extraArch := range extraArches {
		for _, value := range extraArch.Values {
			env := fmt.Sprintf("%s=%s", extraArch.EnvKey, value.Value)
			tasks.PushBack(b.NewTask(TaskContext{
				ArchContext: ctx,
				GOOS:        GOOS,
				GOARCH:      GOARCH,
				NameSuffix:  value.Names(global.Config.ExtraArchesShowDefault),
				Env:         []string{env},
			}))
			if value.ExtraFloat != "" {
				tasks.PushBack(b.NewTask(TaskContext{
					ArchContext: ctx,
					GOOS:        GOOS,
					GOARCH:      GOARCH,
					NameSuffix:  value.NamesExtraFloat(global.Config.ExtraArchesShowDefault),
					Env:         []string{env + "," + value.ExtraFloat},
				}))
			}
		}
	}
	return tasks
}

func (b *Builder) BuildArches(ctx context.Context) error {
	// match all build tasks
	b.TaskChan = make(chan *Task)
	tasks := list.New()
	archCtx := ArchContext{
		Context: ctx,
		Stat:    &BuildStat{},
	}
	for GOOS, Arches := range b.Platforms {
		for _, GOARCH := range Arches {
			extraArches, ok := goCMD.ExtraArches[GOARCH]
			if global.Config.ExtraArches && ok {
				extraArchList := b.CalcExtraArches(archCtx, GOOS, GOARCH, extraArches)
				tasks.PushBackList(extraArchList)
			} else {
				tasks.PushBack(b.NewTask(TaskContext{
					ArchContext: archCtx,
					GOOS:        GOOS,
					GOARCH:      GOARCH,
				}))
			}
		}
	}
	if tasks.Len() == 0 {
		return fmt.Errorf("no valid arch found")
	}
	archCtx.Stat.SetNum(uint32(tasks.Len()))

	log.Infof("found %d arches, building...", tasks.Len())

	// prepare channels and goroutines
	b.FailedTaskChan = make(chan string, tasks.Len())
	thread := int(global.Config.Thread)
	if tasks.Len() < thread {
		thread = tasks.Len()
		log.Debugf("scale down build thread to %d", thread)
	}
	b.WaitGroup = &sync.WaitGroup{}
	b.WaitGroup.Add(thread)
	for range thread {
		go b.BuildThread()
	}

	// start compile
	for el := tasks.Front(); el != nil; el = el.Next() {
		select {
		case <-ctx.Done():
			break
		case b.TaskChan <- el.Value.(*Task):
		}
	}
	close(b.TaskChan)
	b.WaitGroup.Wait()

	// print compile result
	if ctx.Err() != nil {
		log.Errorln("build process interrupted")
		return ctx.Err()
	}
	if len(b.FailedTaskChan) == 0 {
		log.Infoln("completed successfully")
	} else if len(b.FailedTaskChan) == tasks.Len() {
		return errors.New("all arches build failed")
	} else {
		log.Infof("completed: %d arches succeed, %d arches failed", tasks.Len()-len(b.FailedTaskChan), len(b.FailedTaskChan))
		failedArches := make([]string, len(b.FailedTaskChan))
		for i := len(b.FailedTaskChan) - 1; i >= 0; i-- {
			failedArches[i] = <-b.FailedTaskChan
		}
		log.Warnf("failed arches: %s", strings.Join(failedArches, ", "))
	}
	return nil
}
