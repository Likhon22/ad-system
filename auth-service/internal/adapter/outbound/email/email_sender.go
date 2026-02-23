package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/likhon22/ad-system/auth-service/internal/port/outbound"
)

type smtpSender struct {
	host string
	port string
	from string
}

func NewSmtpSender(host, port, from string) outbound.EmailSender {
	return &smtpSender{host: host, port: port, from: from}
}
func (s *smtpSender) SendPasswordResetEmail(ctx context.Context, toEmail, resetLink string) error {
	addr := s.host + ":" + s.port

	subject := "Reset your password"
	body := fmt.Sprintf("Click the link to reset your password: %s\n\nThis link expires in 15 minutes.", resetLink)

	msg := []byte("From: " + s.from + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, nil, s.from, []string{toEmail}, msg)
}
