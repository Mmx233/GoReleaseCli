package builder

import (
	"context"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	log "github.com/sirupsen/logrus"
)

func Run() {
	if err := DownloadGoMod(); err != nil {
		log.Fatalln("download go mod failed:", err)
	}
	builder, err := NewBuilder(global.Config.Output)
	if err != nil {
		log.Fatalln(err)
	}
	if err = builder.BuildArches(context.TODO()); err != nil {
		log.Fatalln(err)
	}
}
