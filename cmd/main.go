package main

import (
	"errors"
	"github.com/KennyMacCormik/HerdMaster/internal/config"
	myinit "github.com/KennyMacCormik/HerdMaster/internal/init"
	"os"
	"os/signal"
	"syscall"
)

const errExit = 1

func main() {
	conf, err := config.New()
	if err != nil {
		lg := myinit.Logger(conf)
		lg.Error("config init error", "error", errors.Unwrap(err).Error())
		os.Exit(errExit)
	}

	lg := myinit.Logger(conf)
	lg.Info("config init success")
	lg.Debug("dumping config", "config", conf)

	endpoint := myinit.Endpoint(conf, lg)
	defer endpoint.Close()
	go endpoint.Run()

	db, err := myinit.StorageDB(conf)
	if err != nil {
		lg.Error("database init error", "error", err.Error())
	}
	defer db.Close()

	// gracefully shutting down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	lg.Info("graceful shutdown done")
}
