package wowclient

import (
	"log"

	"go.uber.org/zap"

	"github.com/kindnessary/wowpow/services/wowclient/internal/client"
	"github.com/kindnessary/wowpow/services/wowclient/internal/config"
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
		log.Println(err)
		return
	}

	client.Run(cfg, logger)
}
