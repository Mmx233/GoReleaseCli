package builder

import (
	"context"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func DownloadGoMod(ctx context.Context) error {
	log.Infoln("downloading go mods...")

	var args []string
	if global.Config.ModDownloadArgs != "" {
		args = strings.Split(global.Config.ModDownloadArgs, " ")
	}
	cmd := goCMD.DownloadMod(ctx, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
