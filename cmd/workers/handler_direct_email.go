package main

import (
	"context"
	"fmt"

	"github.com/swissymissy/Cardinal/internal/pubsub"
)

func (wkrcfg *WorkerConfig) HandlerDirectEmail(event pubsub.DirectEvent) pubsub.AckType {
	// get receiver email
	user, err := wkrcfg.DB.GetUserByID(context.Background(), event.Receiver)
	if err != nil {
		fmt.Printf("Failed to get receiver email: %s\n", err)
		return pubsub.NackRequeue
	}

	// send email
	var subject string
	switch event.Type {
	case "comment":
		subject = fmt.Sprintf("%s commented on your chirp", event.Username)
	case "reaction":
		subject = fmt.Sprintf("%s reacted to your chirp", event.Username)
	case "follow":
		subject = "New Follower"
	}

	if err := wkrcfg.SendEmail(user.Email, subject, event.Body); err != nil {
		fmt.Printf("Failed to send email to %s: %s\n", user.Email, err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}
