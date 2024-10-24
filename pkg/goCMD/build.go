package goCMD

import (
	"context"
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

func (c BuildCommand) _ExecArgs() []string {
	if len(c.ldflags) != 0 {
		c.args = append(c.args, "-ldflags", strings.Join(c.ldflags, " "))
	}
	return append(c.args, c.target)
}

func (c BuildCommand) Exec() *exec.Cmd {
	return exec.Command("go", c._ExecArgs()...)
}

func (c BuildCommand) ExecContext(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, "go", c._ExecArgs()...)
}

func (c BuildCommand) Run() ([]byte, error) {
	return c.Exec().Output()
}

func (c BuildCommand) TrimPath() BuildCommand {
	c.args = append(c.args, "-trimpath")
	return c
}

func (c BuildCommand) Ldflags(value string) BuildCommand {
	c.ldflags = append(c.ldflags, value)
	return c
}

func (c BuildCommand) ProductionLdflags() BuildCommand {
	return c.Ldflags(`-extldflags "-static -fpic" -s -w`)
}

func (c BuildCommand) OutputName(name string) BuildCommand {
	c.args = append(c.args, "-o", name)
	return c
}
