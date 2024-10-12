package global

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "15:04:05",
	})
}
