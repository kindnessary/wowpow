package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"go.uber.org/zap"

	"github.com/kindnessary/wowpow/services/pkg/pow"
)

func (s *Server) handleWOWPOW(conn net.Conn) error {
	pow, err := s.sendTask(conn)
	if err != nil {
		return fmt.Errorf("error sending pow task: %w", err)
	}

	if err := s.checkResult(conn, pow); err != nil {
		return fmt.Errorf("error checking result: %w", err)
	}

	if err := s.respondWithQuote(conn); err != nil {
		return fmt.Errorf("error responding with quote: %w", err)
	}

	return nil
}

func (s *Server) sendTask(conn net.Conn) (*pow.Processor, error) {
	initialData := pow.ToHex(s.randSrc.Uint64())
	lenData := len(initialData)
	pow := pow.New(s.difficulty, initialData)

	message := append([]byte{s.difficulty, byte(lenData)}, initialData...)
	n, err := conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("error writing initial task")
	}
	if n != len(message) {
		return nil, fmt.Errorf("invalid message len wrote, expected %d, wrote %d", len(message), n)
	}

	return pow, nil
}

func (s *Server) checkResult(conn net.Conn, pow *pow.Processor) error {
	token := make([]byte, 8)
	n, err := conn.Read(token)
	if err != nil {
		return fmt.Errorf("error reading token: %w", err)
	}
	if n != len(token) {
		return fmt.Errorf("invalid result token size")
	}

	valid := pow.Validate(binary.BigEndian.Uint64(token))
	if !valid {
		return fmt.Errorf("invalid result token")
	}

	return nil
}

func (s *Server) respondWithQuote(conn net.Conn) error {
	quote, err := s.repo.GetQuote(context.Background(), s.randSrc.Intn(s.numOfQuotes)+1)
	if err != nil {
		return fmt.Errorf("error getting quote: %w", err)
	}

	s.logger.Info("responding quote", zap.Int("quote_id", quote.ID))

	message := []byte(quote.Text)
	n, err := conn.Write(message)
	if err != nil {
		return fmt.Errorf("error writing quote: %w", err)
	}
	if n != len(message) {
		return fmt.Errorf("invalid message len wrote, expected %d, wrote %d", len(message), n)
	}

	return nil
}
