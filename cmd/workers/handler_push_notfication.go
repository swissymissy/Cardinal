package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/swissymissy/Cardinal/internal/pubsub"
)

func (wkrcfg *WorkerConfig) HandlerPushNotification(event pubsub.ChirpEvent) pubsub.AckType {
	// get followers ID list
	list, err := wkrcfg.getFollowersID(context.Background(), event.Triggerer)
	if err != nil {
		fmt.Printf("Failed to get followers' ID: %s\n", err)
		return pubsub.NackRequeue
	}
	if len(list) == 0 {
		fmt.Println("No followers to send notfication")
		return pubsub.Ack
	}

	// update notfications table
	err = wkrcfg.saveNotifications(context.Background(), event, list)
	if err != nil {
		fmt.Printf("Failed to update notfications table: %s\n", err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}

// helper 1: get Followers ID List
func (wkrcfg *WorkerConfig) getFollowersID(ctx context.Context, triggererID uuid.UUID) ([]uuid.UUID, error) {
	followerList, err := wkrcfg.DB.GetFollowers(ctx, triggererID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get followers ID: %w", err)
	}
	IDList := make([]uuid.UUID, len(followerList))
	for i, id := range followerList {
		IDList[i] = id.FollowerID
	}
	return IDList, nil
}

// helper 2: write new notfication info to Notification table in bulk
func (wkrcfg *WorkerConfig) saveNotifications(ctx context.Context, event pubsub.ChirpEvent, idList []uuid.UUID) error {
	msg := fmt.Sprintf("%s just posted new chirp: %s", event.Username, truncate(event.Body, 50))
	err := wkrcfg.DB.CreateNotificationsBulk(ctx, database.CreateNotificationsBulkParams{
		Body:      msg,
		Column2:   idList,
		Triggerer: event.Triggerer,
		ChirpID:   event.ChirpID,
	})
	if err != nil {
		return fmt.Errorf("Errors writing new notification to notifications table: %w", err)
	}
	return nil
}

// helper 3: truncate the body of new posted chirp to include it in the notification
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max] + "...")
}
