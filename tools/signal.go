package tools

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func OsSignalCancel(cancel context.CancelFunc) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	log.Infoln("Cancel Tasks...")
	cancel()
}
