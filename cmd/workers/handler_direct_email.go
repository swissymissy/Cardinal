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
	var body string
	switch event.Type {
	case "comment":
		subject = fmt.Sprintf("Cardinal: %s commented on your chirp", event.Username)
		body = fmt.Sprintf("User %s commented on your chirp: \n\n%s", event.Username, event.Body)
	case "reaction":
		subject = fmt.Sprintf("Cardinal: %s reacted to your chirp", event.Username)
		body = event.Body
	case "follow":
		subject = "Cardinal: You have new Follower"
		body = fmt.Sprintf("User %s has started following you! Yay!", event.Username)
	}

	
	if err := wkrcfg.SendEmail(user.Email, subject, body); err != nil {
		fmt.Printf("Failed to send email to %s: %s\n", user.Email, err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}
