package main

import (
	"log"
	"net/smtp"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	from     string
	password string
}

func NewEmailService(host, port, from, password string) *EmailService {
	return &EmailService{
		smtpHost: host,
		smtpPort: port,
		from:     from,
		password: password,
	}
}

func (s *EmailService) SendEmail(email Email) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	to := []string{email.To}
	msg := []byte("To: " + email.To + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" +
		email.Body + "\r\n")

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, s.from, to, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent to %s", email.To)
	return nil
}
