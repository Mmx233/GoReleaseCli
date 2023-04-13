package builder

import (
	log "github.com/sirupsen/logrus"
)

const (
	StoreDir = "build"
)

func Run() {
	builder, e := NewBuilder(StoreDir)
	if e != nil {
		log.Fatalln(e)
	}
	builder.BuildArches()
}
