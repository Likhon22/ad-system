package config

import (
	"fmt"
	"os"
)

type EmailConfig struct {
	SmtpHost  string
	SmtpPort  string
	EmailFrom string
}

func loadEmailConfig() (*EmailConfig, error) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	emailFrom := os.Getenv("EMAIL_FROM")

	if smtpHost == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}
	if smtpPort == "" {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}

	if emailFrom == "" {
		return nil, fmt.Errorf("EMAIL_FROM is required")
	}

	return &EmailConfig{

		SmtpPort:  smtpPort,
		SmtpHost:  smtpHost,
		EmailFrom: emailFrom,
	}, nil
}
