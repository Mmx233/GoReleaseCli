package builder

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func DownloadGoMod() error {
	log.Infoln("downloading go mods...")

	var args []string
	if global.Config.ModDownloadArgs != "" {
		args = strings.Split(global.Config.ModDownloadArgs, " ")
	}
	cmd := goCMD.DownloadMod(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
