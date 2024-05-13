package compress

import (
	"os/exec"
	"path"
	"strings"
)

type PathInfo struct {
	// src info
	BaseName string // filename without parent path
	Ext      string
	Name     string // basename without ext
	Dir      string

	// dist info
	OutputPath string
	TargetName string // rename src file
}

func (i PathInfo) Exec(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = i.Dir
	return cmd
}

func DecodePathInfo(filePath, targetName, targetFormat string) *PathInfo {
	baseName := path.Base(filePath)
	ext := path.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)
	dir := strings.TrimSuffix(filePath, baseName)

	return &PathInfo{
		BaseName: baseName,
		Ext:      ext,
		Name:     name,
		Dir:      dir,

		OutputPath: name + "." + targetFormat,
		TargetName: targetName + ext,
	}
}

type Make func(filePath, targetName string) error

func _NewMakeFn(targetFormat string, fn func(info *PathInfo) error) Make {
	return func(filePath, targetName string) error {
		return fn(DecodePathInfo(filePath, targetName, targetFormat))
	}
}
