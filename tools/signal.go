package tools

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func OsSignalCancel(cancel context.CancelFunc) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	log.Infoln("Cancel Tasks...")
	cancel()
}
