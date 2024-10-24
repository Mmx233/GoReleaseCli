package builder

import (
	"context"
	"errors"
	"github.com/Mmx233/GoReleaseCli/internal/global"
	log "github.com/sirupsen/logrus"
)

func Run(ctx context.Context) {
	conf := Config{
		OutputDir: global.Config.Output,
	}
	switch global.Config.OutputFormat {
	case "":
	case "post":
		conf.OutputFormatter = PostOutputFormatter{}
	default:
		log.Warnln("Unsupported output format:", global.Config.OutputFormat)
	}

	if err := DownloadGoMod(ctx); err != nil {
		log.Fatalln("download go mod failed:", err)
	}
	builder, err := NewBuilder(conf)
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
