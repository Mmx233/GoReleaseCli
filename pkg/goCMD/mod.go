package goCMD

import (
	"context"
	"os/exec"
)

func DownloadMod(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "go", append([]string{"mod", "download"}, args...)...)
}
