package builder

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	"github.com/Mmx233/GoReleaseCli/pkg/goCMD"
	log "github.com/sirupsen/logrus"
)

func Run() {
	if err := goCMD.DownloadMod().Run(); err != nil {
		log.Fatalln("download go mod failed:", err)
	}
	builder, err := NewBuilder(global.Config.Output)
	if err != nil {
		log.Fatalln(err)
	}
	if err = builder.BuildArches(); err != nil {
		log.Fatalln(err)
	}
}
