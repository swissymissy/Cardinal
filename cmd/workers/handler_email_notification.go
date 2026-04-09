package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/pubsub"
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

	subject := fmt.Sprintf("New chirp from %s!", event.Username)
	body := fmt.Sprintf("%s has posted a new chirp: \n\n%s", event.Username, event.Body)

	for _, email := range emails {
		// send email to follower
		if err := wkrcfg.SendEmail(email, subject, body); err != nil {
			fmt.Printf("Failed to send email to %s: %s\n", email, err)
		}
	}
	return pubsub.Ack
}

// helper: function geting followers email from DB
func (wkrcfg *WorkerConfig) getFollowerEmails(ctx context.Context, triggererID uuid.UUID) ([]string, error) {
	followersEmail, err := wkrcfg.DB.GetFollowersEmail(ctx, triggererID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get follower emails: %w", err)
	}

	// extract emails
	emails := make([]string, len(followersEmail))
	for i, e := range followersEmail {
		emails[i] = e.Email
	}
	return emails, nil
}
