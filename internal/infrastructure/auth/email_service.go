package auth

import (
	"fmt"
	"net/smtp"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(to, subject, body string) error
}

type smtpEmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewSMTPEmailService creates a new SMTP-based EmailService instance
func NewSMTPEmailService(host string, port int, username, password, from string) EmailService {
	return &smtpEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendEmail sends an email
func (s *smtpEmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
