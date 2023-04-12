package models

import "github.com/alecthomas/kingpin/v2"

type Commands struct {
	App        *kingpin.Application
	Target     string
	Ldflags    string
	SoftFloat  bool
	OutputName string
	OS         []string
	Arch       []string
}
