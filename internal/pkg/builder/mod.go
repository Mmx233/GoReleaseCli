package builder

import (
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	"os"
)

func DownloadGoMod() error {
	cmd := goCMD.DownloadMod()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
