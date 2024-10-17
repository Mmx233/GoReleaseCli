package builder

import (
	"context"
	"errors"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	log "github.com/sirupsen/logrus"
)

func Run(ctx context.Context) {
	if err := DownloadGoMod(ctx); err != nil {
		log.Fatalln("download go mod failed:", err)
	}
	builder, err := NewBuilder(global.Config.Output)
	if err != nil {
		log.Fatalln(err)
	}
	if ctx.Err() != nil {
		return
	}
	if err = builder.BuildArches(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalln(err)
	}
}
