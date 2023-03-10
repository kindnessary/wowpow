package wowserver

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/kindnessary/wowpow/services/wowserver/internal/config"
	"github.com/kindnessary/wowpow/services/wowserver/internal/repository"
	"github.com/kindnessary/wowpow/services/wowserver/internal/server"
)

func Run() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println(err)
		return
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("error loading config", zap.Error(err))
		return
	}

	db, err := New(cfg.SQLDatabase, logger)
	if err != nil {
		logger.Error("error creating db", zap.Error(err))
		return
	}

	srv, err := server.New(cfg.General, repository.New(db), logger)
	if err != nil {
		logger.Error("error creating server", zap.Error(err))
		return
	}

	logger.Info("starting server", zap.String("address", cfg.General.Address))
	go srv.ListenAndServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("stopping server")
	srv.Stop()
}
