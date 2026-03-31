package main

import (
	"context"
	"fmt"

	"github.com/swissymissy/Cardinal/internal/pubsub"
	"github.com/wneessen/go-mail"
	"github.com/google/uuid"
)

// handler to send email notification to the followers' email
func (wkrcfg *WorkerConfig) HandlerEmailNotification(event pubsub.ChirpEvent) pubsub.AckType {
	// get followers email list
	emails, err := wkrcfg.getFollowerEmails(context.Background(), event.Triggerer)
	if err != nil {
		fmt.Printf("Failed to get follower emails: %s\n", err)
		return pubsub.NackRequeue
	}
	if len(emails) == 0 {
		fmt.Println("No follower to send email")
		return pubsub.Ack 
	}

	for _, email := range emails {
		// send email to follower
		if err := wkrcfg.sendChirpEmail(email, event); err != nil {
			fmt.Printf("Failed to send email to %s: %s\n", email, err)
		}
	}
	return pubsub.Ack
}


// helper 1: function geting followers email from DB
func (wkrcfg *WorkerConfig) getFollowerEmails(ctx context.Context, triggererID uuid.UUID) ([]string, error) {
	followersEmail, err := wkrcfg.DB.GetFollowersEmail(ctx, triggererID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get follower emails: %w", err)
	}

	// extract emails
	emails := make([]string, len(followersEmail))
	for i, e := range followersEmail{
		emails[i] = e.Email
	}
	return emails, nil
}

// helper 2: function to send email to user
func (wkrcfg *WorkerConfig) sendChirpEmail(email string, event pubsub.ChirpEvent) error {
	// create email message
	message := mail.NewMsg()
	if err := message.From(wkrcfg.SMTPUsername); err != nil {
		return fmt.Errorf("Failed to set From address: %w", err)
	}
	if err := message.To(email); err != nil {
		return fmt.Errorf("Failed to set To address: %w", err)
	}
	message.Subject("New chirp from someone you follow!")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(
		"Someone you follow posted a new chirp:\n\n%s", event.Body,
	))

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