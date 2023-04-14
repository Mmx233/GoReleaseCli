package builder

import (
	"github.com/Mmx233/GoReleaseCli/internal/global"
	log "github.com/sirupsen/logrus"
)

func Run() {
	builder, e := NewBuilder(global.Commands.Output)
	if e != nil {
		log.Fatalln(e)
	}
	builder.BuildArches()
}
