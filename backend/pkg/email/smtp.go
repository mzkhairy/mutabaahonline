package email

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string // e.g. "smtp.gmail.com"
	Port     int    // e.g. 587
	Username string
	Password string
}

type SMTPSender struct {
	cfg SMTPConfig
}

func NewSMTPSender(cfg SMTPConfig) *SMTPSender { return &SMTPSender{cfg: cfg} }

func (s *SMTPSender) Send(ctx context.Context, msg Message) error {
	if err := msg.Validate(); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	data := buildMIME(&msg)
	// net/smtp does not accept context; for timeouts, use a context-aware dialer in a custom impl.
	return smtp.SendMail(addr, auth, msg.From, allRecipients(&msg), data)
}
