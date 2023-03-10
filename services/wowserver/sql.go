package wowserver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	"github.com/kindnessary/wowpow/services/wowserver/internal/config"
)

func New(cfg config.SQLDatabase, logger *zap.Logger) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	var err error
	var session *sql.DB

	for attempt := 1; attempt <= cfg.ConnectionAttempts; attempt++ {
		session, err = sql.Open(cfg.StorageDriver, url)
		if err != nil {
			logger.Error("unable to connect to DB. retrying...", zap.Error(err))
			time.Sleep(cfg.ConnectionInterval)
			continue
		}

		logger.Info("connected to DB")
		err = session.Ping()
		if err != nil {
			logger.Error("unable to ping DB. retrying....", zap.Error(err))
			time.Sleep(cfg.ConnectionInterval)
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("unable to connect to DB: %w", err)
	}

	session.SetMaxOpenConns(cfg.MaxOpenConns)
	session.SetMaxIdleConns(cfg.MaxIdleConns)
	session.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return session, nil
}
