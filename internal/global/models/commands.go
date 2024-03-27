package models

import "github.com/alecthomas/kingpin/v2"

type Commands struct {
	App       *kingpin.Application
	Target    string
	Ldflags   string
	SoftFloat bool
	Cgo       bool
	OS        string
	Arch      string
	Output    string
	Name      string
	Thread    uint16
}
