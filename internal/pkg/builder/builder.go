package builder

import (
	"container/list"
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

	platforms, err := MatchTargetPlatforms()
	if err != nil {
		return nil, err
	}

	builder := &Builder{
		OutputName: outputName,
		OutputDir:  outputDir,
		Platforms:  platforms,
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

	if err := PrepareDirs(outputDir); err != nil {
		return nil, err
	}

	return builder, nil
}

type Task struct {
	Name  string
	Build func() error
}

type Builder struct {
	OutputName string
	OutputDir  string
	Cgo        string
	Platforms  map[string][]string
	GoCMD      goCMD.BuildCommand

	WaitGroup      *sync.WaitGroup
	TaskChan       chan *Task
	FailedTaskChan chan string
	Compress       compress.Make
}

func (b *Builder) NewTask(GOOS, GOARCH string, nameSuffix, env []string) *Task {
	outputName := b.OutputName
	if GOOS == "windows" {
		outputName += ".exe"
	}
	buildName := BuildName(outputName, append([]string{GOOS, GOARCH}, nameSuffix...)...)
	outputPath := path.Join(b.OutputDir, buildName)

	cmd := b.GoCMD.OutputName(outputPath).Exec()
	cmd.Env = append(cmd.Environ(), env...)
	cmd.Env = append(cmd.Env, b.Cgo, "GOOS="+GOOS, "GOARCH="+GOARCH)
	cmd.Stderr = os.Stderr

	return &Task{
		Name: buildName,
		Build: func() error {
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("build error: %v", err)
			}
			if err = b.Compress(outputPath, b.OutputName); err != nil {
				return fmt.Errorf("compress %s failed: %v", outputPath, err)
			}
			return nil
		},
	}
}

func (b *Builder) BuildThread(stat *BuildStat) {
	for {
		task, ok := <-b.TaskChan
		if !ok {
			b.WaitGroup.Done()
			return
		}
		err := task.Build()
		stat.Done()
		logger := log.WithField("process", fmt.Sprintf("%.1f", stat.Percentage()))
		if err != nil {
			logger.Errorf("error occur while building %s: %v", task.Name, err)
			b.FailedTaskChan <- task.Name
		} else {
			logger.Infof("build %s success", task.Name)
		}
	}
}

func (b *Builder) CalcExtraArches(GOOS, GOARCH string, extraArches []goCMD.ArchExtra) *list.List {
	tasks := list.New()
	for _, extraArch := range extraArches {
		for _, value := range extraArch.Values {
			env := fmt.Sprintf("%s=%s", extraArch.EnvKey, value.Value)
			tasks.PushBack(b.NewTask(GOOS, GOARCH, value.Names(global.Config.ExtraArchesShowDefault), []string{env}))
			if value.ExtraFloat != "" {
				tasks.PushBack(b.NewTask(GOOS, GOARCH, value.NamesExtraFloat(global.Config.ExtraArchesShowDefault), []string{env + "," + value.ExtraFloat}))
			}
		}
	}
	return tasks
}

func (b *Builder) BuildArches() error {
	// match all build tasks
	var tasks = list.New()
	b.TaskChan = make(chan *Task)
	for GOOS, Arches := range b.Platforms {
		for _, GOARCH := range Arches {
			if global.Config.ExtraArches {
				extraArches, ok := goCMD.ExtraArches[GOARCH]
				if !ok {
					tasks.PushBack(b.NewTask(GOOS, GOARCH, nil, nil))
					continue
				}
				tasks.PushBackList(b.CalcExtraArches(GOOS, GOARCH, extraArches))
			} else {
				tasks.PushBack(b.NewTask(GOOS, GOARCH, nil, nil))
			}
		}
	}
	if tasks.Len() == 0 {
		return fmt.Errorf("no valid arch found")
	}

	log.Infof("found %d arches, building...", tasks.Len())

	// prepare channels and goroutines
	b.FailedTaskChan = make(chan string, tasks.Len())
	thread := int(global.Config.Thread)
	if tasks.Len() < thread {
		thread = tasks.Len()
	}
	b.WaitGroup = &sync.WaitGroup{}
	b.WaitGroup.Add(thread)
	stat := &BuildStat{
		Num: uint32(tasks.Len()),
	}
	for range thread {
		go b.BuildThread(stat)
	}

	// start compile
	for el := tasks.Front(); el != nil; el = el.Next() {
		b.TaskChan <- el.Value.(*Task)
	}
	close(b.TaskChan)
	b.WaitGroup.Wait()

	// print compile result
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
		log.Infof("failed arches: %s", strings.Join(failedArches, ", "))
	}
	return nil
}
