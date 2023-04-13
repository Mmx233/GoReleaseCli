package models

import "github.com/alecthomas/kingpin/v2"

type Commands struct {
	App       *kingpin.Application
	Target    string
	Ldflags   string
	SoftFloat bool
	OS        string
	Arch      string
	Output    string
	Name      string
}
