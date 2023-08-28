package builder

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	log "github.com/sirupsen/logrus"
)

func Run() {
	builder, err := NewBuilder(global.Commands.Output)
	if err != nil {
		log.Fatalln(err)
	}
	builder.BuildArches()
}
