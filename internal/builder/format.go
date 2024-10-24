package builder

import (
	"path"
	"strings"
)

type OutputFormatter interface {
	Format(name, divider string, args ...string) (string, error)
}

type DefaultOutputFormatter struct{}

func (DefaultOutputFormatter) Format(filename, divider string, args ...string) (string, error) {
	ext := path.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	for _, s := range args {
		name += divider + s
	}
	name += ext
	return name, nil
}

type PostOutputFormatter struct{}

func (PostOutputFormatter) Format(filename, divider string, args ...string) (string, error) {
	var name string
	for _, s := range args {
		name += s + divider
	}
	return name + filename, nil
}
