package tools

import (
	"os/exec"
	"path"
	"strings"
)

type CompressInfo struct {
	// src info
	BaseName string
	Ext      string
	Name     string
	Dir      string

	// dist info
	OutputPath string
	TargetName string // rename src file
}

func DecodeCompressPath(filePath, targetName, targetFormat string) *CompressInfo {
	baseName := path.Base(filePath)
	ext := path.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)
	dir := strings.TrimSuffix(filePath, baseName)

	return &CompressInfo{
		BaseName: baseName,
		Ext:      ext,
		Name:     name,
		Dir:      dir,

		OutputPath: name + "." + targetFormat,
		TargetName: targetName + ext,
	}
}

func MakeZip(filePath, targetName string) error {
	info := DecodeCompressPath(filePath, targetName, "zip")

	cmd := exec.Command("7z", "a", "-tzip", info.OutputPath, info.BaseName, "-mx9")
	cmd.Dir = info.Dir
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("7z", "rn", info.OutputPath, info.BaseName, info.TargetName)
	cmd.Dir = info.Dir
	return cmd.Run()
}
