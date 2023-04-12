package goCMD

import "os/exec"

func NewBuilder(target string) BuildCommand {
	return BuildCommand{
		args:   []string{"build"},
		target: target,
	}
}

type BuildCommand struct {
	args   []string
	target string
}

func (c BuildCommand) Run() ([]byte, error) {

	return exec.Command("go", append(c.args, c.target)...).Output()
}

func (c BuildCommand) TrimPath() BuildCommand {
	c.args = append(c.args, "-gcflags=-trimpath=$GOPATH", "-asmflags=-trimpath=$GOPATH")
	return c
}

func (c BuildCommand) Ldflags(value string) BuildCommand {
	c.args = append(c.args, "-ldflags", value)
	return c
}

func (c BuildCommand) ProductionLdflags() BuildCommand {
	return c.Ldflags(`'-extldflags "-static" -s -w'`)
}

func (c BuildCommand) OutputName(name string) BuildCommand {
	c.args = append(c.args, "-o", name)
	return c
}
