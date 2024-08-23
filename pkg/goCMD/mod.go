package goCMD

import "os/exec"

func DownloadMod(args ...string) *exec.Cmd {
	return exec.Command("go", append([]string{"mod", "download"}, args...)...)
}
