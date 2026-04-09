package main

import (
	"context"
	"fmt"

	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/swissymissy/Cardinal/internal/pubsub"
)

func (wkrcfg *WorkerConfig) HandlerDirectPush(event pubsub.DirectEvent) pubsub.AckType {
	// check follow or react-comment
	if event.ChirpID != nil {
		_, err := wkrcfg.DB.CreateNotifications(context.Background(), database.CreateNotificationsParams{
			Body:      event.Body,
			Receiver:  event.Receiver,
			Triggerer: event.Triggerer,
			ChirpID:   event.ChirpID,
		})
		if err != nil {
			fmt.Printf("Failed to save new react-comment notification to db: %s\n", err)
			return pubsub.NackRequeue
		}
	} else {
		_, err := wkrcfg.DB.CreateFollowNotification(context.Background(), database.CreateFollowNotificationParams{
			Body:      event.Body,
			Receiver:  event.Receiver,
			Triggerer: event.Triggerer,
		})
		if err != nil {
			fmt.Printf("Failed to save new follow notification to db: %s\n", err)
			return pubsub.NackRequeue
		}
	}
	return pubsub.Ack
}
