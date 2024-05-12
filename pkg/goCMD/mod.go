package goCMD

import "os/exec"

func DownloadMod() *exec.Cmd {
	return exec.Command("go", "mod", "download", "-x")
}
