package services

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func NewEmailService(host, port, user, pass, from string) *EmailService {
	return &EmailService{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		From: from,
	}
}

func (s *EmailService) Send(to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.User, s.Pass, s.Host)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(addr, auth, s.User, []string{to}, []byte(msg.String()))
}
