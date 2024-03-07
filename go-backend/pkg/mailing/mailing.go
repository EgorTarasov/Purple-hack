package mailing

import (
	"crypto/tls"
	"fmt"
	"github.com/yogenyslav/logger"
	"net/smtp"
)

type MailServer struct {
	client *smtp.Client
	sender string
	Addr   string
}

func NewMailServer(cfg *Config) *MailServer {
	var err error

	tlsConfig := &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: true,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		logger.Panic(err)
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		logger.Panic(err)
	}

	if err = client.Auth(smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.Host)); err != nil {
		logger.Panic(err)
	}

	return &MailServer{
		client: client,
		Addr:   addr,
		sender: cfg.Email,
	}
}

func (s *MailServer) Send(to string, msg []byte) error {
	var err error

	if err = s.client.Mail(s.sender); err != nil {
		return err
	}

	if err = s.client.Rcpt(to); err != nil {
		return err
	}

	wc, err := s.client.Data()
	if err != nil {
		return err
	}

	_, err = wc.Write(msg)
	if err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return err
	}

	return nil
}
