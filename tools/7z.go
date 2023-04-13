package tools

import (
	"os/exec"
	"path"
	"strings"
)

func MakeZip(filePath string) error {
	baseName := path.Base(filePath)
	ext := path.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)
	dir := strings.TrimSuffix(filePath, baseName)

	cmd := exec.Command("7z", "a", name+".zip", baseName)
	cmd.Dir = dir
	_, e := cmd.Output()
	return e
}
