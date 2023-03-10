package server

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/kindnessary/wowpow/services/pkg/message"
	"github.com/kindnessary/wowpow/services/wowserver/internal/config"
	"github.com/kindnessary/wowpow/services/wowserver/internal/repository"
)

type Server struct {
	wg            *sync.WaitGroup
	listener      net.Listener
	aliveListener net.Listener

	numOfQuotes        int
	difficulty         uint8
	connectionLifetime time.Duration
	repo               *repository.Repository

	logger *zap.Logger

	randSrc *rand.Rand
}

func New(cfg config.General, repo *repository.Repository, logger *zap.Logger) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, err
	}
	aliveListener, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		return nil, err
	}

	return &Server{
		wg:                 &sync.WaitGroup{},
		listener:           listener,
		aliveListener:      aliveListener,
		numOfQuotes:        cfg.NumOfQuotes,
		difficulty:         cfg.Difficulty,
		connectionLifetime: cfg.ConnectionLifetime,
		repo:               repo,
		logger:             logger,
		randSrc:            rand.New(rand.NewSource(time.Now().Unix())),
	}, nil
}

func (s *Server) Stop() {
	err := s.listener.Close()
	if err != nil {
		s.logger.Error("error closing listener", zap.Error(err))
	}
	s.wg.Wait()
	err = s.aliveListener.Close()
	if err != nil {
		s.logger.Error("error closing alive listener", zap.Error(err))
	}
}

func (s *Server) ListenAndServer() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Info("closing listener")
				return
			}
			s.logger.Error("error accepting connection", zap.Error(err))
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	s.logger.Info("received new connection", zap.String("remote address", conn.RemoteAddr().String()))

	if err := conn.SetDeadline(time.Now().Add(s.connectionLifetime)); err != nil {
		s.logger.Error("error setting connection lifetime", zap.Error(err))
		return
	}

	msg := make([]byte, 1)
	n, err := conn.Read(msg)
	if err != nil {
		s.logger.Error("error reading initial message", zap.Error(err))
		return
	}
	if n != 1 {
		s.logger.Error("invalid message len", zap.Int("read bytes", n))
		return
	}

	switch msg[0] {
	case message.MessagePing:
		err = s.handlerPing(conn)
	case message.MessageWOW:
		err = s.handleWOWPOW(conn)
	default:
		err = fmt.Errorf("invalid initial message: %v", msg[0])
	}

	if err != nil {
		s.logger.Error("error handling connection", zap.Error(err))
		return
	}

	s.logger.Info("connection successfully handled", zap.String("remote address", conn.RemoteAddr().String()))
}

func (s *Server) handlerPing(conn net.Conn) error {
	n, err := conn.Write([]byte{message.MessagePong})
	if err != nil {
		return fmt.Errorf("error writing pong response: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("invalid message len wrote, expected 1, wrote %d", n)
	}

	return nil
}
