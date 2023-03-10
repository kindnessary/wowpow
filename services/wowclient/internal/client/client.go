package client

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/kindnessary/wowpow/services/pkg/message"
	"github.com/kindnessary/wowpow/services/pkg/pow"
	"github.com/kindnessary/wowpow/services/wowclient/internal/config"
)

const (
	maxServerPingAttempts   = 5
	serverPingRetryInterval = time.Second
)

func Run(cfg *config.Configuration, logger *zap.Logger) {
	wg := sync.WaitGroup{}
	wg.Add(cfg.General.NumOfClients)

	if err := pingServer(cfg, logger); err != nil {
		logger.Error("error ping server", zap.Error(err))
		return
	}

	for i := 0; i < cfg.General.NumOfClients; i++ {
		go func() {
			defer wg.Done()

			cl, err := new(cfg.General.ServerAddress, logger)
			if err != nil {
				logger.Error("error creating new client", zap.Error(err))
				return
			}

			cl.getQuote()
		}()
	}

	wg.Wait()
}

func pingServer(cfg *config.Configuration, logger *zap.Logger) error {
	for i := 0; i < maxServerPingAttempts; i++ {
		conn, err := net.Dial("tcp", cfg.General.ServerAddress)
		if err != nil {
			logger.Info("failed to connect to server", zap.Error(err))
			time.Sleep(serverPingRetryInterval)
			continue
		}

		n, err := conn.Write([]byte{message.MessagePing})
		if err != nil {
			return fmt.Errorf("error writing ping message: %w", err)
		}
		if n != 1 {
			return fmt.Errorf("wrong len ping message, expected 1, wrote %d", n)
		}

		response := make([]byte, 1)
		n, err = conn.Read(response)
		if err != nil {
			return fmt.Errorf("error reading pong response: %w", err)
		}
		if n != 1 {
			return fmt.Errorf("wrong len pong message, expected 1, read %d", n)
		}

		if response[0] != message.MessagePong {
			return fmt.Errorf("invalid pong response: %v", response[0])
		}

		return nil
	}

	return fmt.Errorf("failed to dial server")
}

type client struct {
	conn net.Conn

	logger *zap.Logger
}

func new(addr string, logger *zap.Logger) (*client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	return &client{
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *client) getQuote() {
	if err := c.resolvePOW(); err != nil {
		c.logger.Error("error resolving POW", zap.Error(err))
		return
	}

	quote, err := io.ReadAll(c.conn)
	if err != nil {
		c.logger.Error("error reading quote", zap.Error(err))
		return
	}

	c.logger.Info("received quote", zap.String("quote", string(quote)))
}

func (c *client) resolvePOW() error {
	n, err := c.conn.Write([]byte{message.MessageWOW})
	if err != nil {
		return fmt.Errorf("error writing wow message: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("wrong len wow message, expected 1, wrote %d", n)
	}

	metadata := make([]byte, 2)
	n, err = c.conn.Read(metadata)
	if err != nil {
		return fmt.Errorf("error reading metadata: %w", err)
	}
	if n != 2 {
		return fmt.Errorf("wrong len of metadata, expected 2, read %d", n)
	}

	difficulty, lenData := metadata[0], metadata[1]
	initialData := make([]byte, lenData)

	n, err = c.conn.Read(initialData)
	if err != nil {
		return fmt.Errorf("error reading initial data: %w", err)
	}
	if n != len(initialData) {
		return fmt.Errorf("wrong len of initial data, expected %d, read %d", len(initialData), n)
	}

	powCtrl := pow.New(difficulty, initialData)

	res := pow.ToHex(powCtrl.Calculate())
	n, err = c.conn.Write(res)
	if err != nil {
		return fmt.Errorf("error writing response token: %w", err)
	}
	if n != len(res) {
		return fmt.Errorf("wrong len result, expected %d, wrote %d", len(res), n)
	}
	return nil
}
