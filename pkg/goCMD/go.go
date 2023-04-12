package goCMD

import (
	"os/exec"
	"strings"
)

func NewBuilder(target string) BuildCommand {
	return BuildCommand{
		args:   []string{"build"},
		target: target,
	}
}

type BuildCommand struct {
	args    []string
	target  string
	ldflags []string
}

func (c BuildCommand) Exec() *exec.Cmd {
	if len(c.ldflags) != 0 {
		c.args = append(c.args, "-ldflags", strings.Join(c.ldflags, " "))
	}
	return exec.Command("go", append(c.args, c.target)...)
}

func (c BuildCommand) Run() ([]byte, error) {
	return c.Exec().Output()
}

func (c BuildCommand) TrimPath() BuildCommand {
	c.args = append(c.args, "-gcflags=-trimpath=$GOPATH", "-asmflags=-trimpath=$GOPATH")
	return c
}

func (c BuildCommand) Ldflags(value string) BuildCommand {
	c.ldflags = append(c.ldflags, value)
	return c
}

func (c BuildCommand) ProductionLdflags() BuildCommand {
	return c.Ldflags(`-extldflags "-static" -s -w`)
}

func (c BuildCommand) OutputName(name string) BuildCommand {
	c.args = append(c.args, "-o", name)
	return c
}
