package tools

import (
	"os/exec"
	"path"
	"strings"
)

func MakeZip(filePath string, targetName string) error {
	baseName := path.Base(filePath)
	ext := path.Ext(baseName)
	name := strings.TrimSuffix(baseName, ext)
	dir := strings.TrimSuffix(filePath, baseName)

	outputPath := name + ".zip"

	cmd := exec.Command("7z", "a", "-tzip", outputPath, baseName, "-mx9")
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("7z", "rn", outputPath, baseName, targetName+ext)
	cmd.Dir = dir
	return cmd.Run()
}
