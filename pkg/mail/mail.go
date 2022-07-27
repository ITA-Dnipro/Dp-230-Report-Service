package mail

import (
	"errors"
	"net/smtp"

	"go.uber.org/zap"
)

type MailConfiguration struct {
	Host     string
	Port     string
	Mail     string
	Password string
}

type MailSender struct {
	name     string
	host     string
	port     string
	mail     string
	password string
	logger   *zap.Logger
}

func NewMailSender(name string, config interface{}, log *zap.Logger) (*MailSender, error) {
	cfg, ok := config.(*MailConfiguration)
	if !ok || cfg == nil {
		return nil, errors.New("invalid mail sender config")
	}
	return &MailSender{
		name:     name,
		logger:   log,
		host:     cfg.Host,
		port:     cfg.Port,
		mail:     cfg.Mail,
		password: cfg.Password,
	}, nil
}

func (s *MailSender) Send(to string, message []byte) error {
	toEmails := []string{to}
	auth := smtp.PlainAuth("", s.mail, s.password, s.host)
	address := s.host + ":" + s.port
	return smtp.SendMail(address, auth, s.mail, toEmails, message)
}
