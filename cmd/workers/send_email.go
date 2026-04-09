package main

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

// send email to user
func (wkrcfg *WorkerConfig) SendEmail(to, subject, body string) error {
	// create email message
	message := mail.NewMsg()
	if err := message.From(wkrcfg.SMTPUsername); err != nil {
		return fmt.Errorf("Failed to set From address: %w", err)
	}
	if err := message.To(to); err != nil {
		return fmt.Errorf("Failed to set To address: %w", err)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)

	// create client and send
	client, err := mail.NewClient(
		wkrcfg.SMTPHost,
		mail.WithPort(wkrcfg.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(wkrcfg.SMTPUsername),
		mail.WithPassword(wkrcfg.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("Failed to create email client: %w", err)
	}
	if err = client.DialAndSend(message); err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}
	return nil
}
