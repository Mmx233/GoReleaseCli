package builder

import (
	"path"
	"strings"
)

type OutputFormatter interface {
	Format(name, divider string, args ...string) string
}

type DefaultOutputFormatter struct{}

func (DefaultOutputFormatter) Format(filename, divider string, args ...string) string {
	ext := path.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	for _, s := range args {
		name += divider + s
	}
	name += ext
	return name
}

type PostOutputFormatter struct{}

func (PostOutputFormatter) Format(filename, divider string, args ...string) string {
	var name string
	for _, s := range args {
		name += s + divider
	}
	return name + filename
}
