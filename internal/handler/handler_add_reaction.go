package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/swissymissy/Cardinal/internal/pubsub"
)

type ReactionRequest struct {
	Type string `json:"type"`
}

var validReactions = map[string]bool{
	"❤️": true,
	"😂":  true,
	"😮":  true,
	"😢":  true,
	"👍":  true,
}

func (apicfg *ApiConfig) HandlerAddReaction(w http.ResponseWriter, r *http.Request) {
	// auth check
	// check user's token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// get chirp ID
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("Failed to parse chirp ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// decode request for reaction type
	var req ReactionRequest
	err = DecodeRequest(r, &req)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err)
		ResponseWithError(w, 400, "Invalid request")
		return
	}

	// check reaction type
	if !validReactions[req.Type] {
		ResponseWithError(w, 400, "Invalid reaction type")
		return
	}

	// add - update reaction
	react, err := apicfg.DB.AddReaction(r.Context(), database.AddReactionParams{
		ChirpID: chirpID,
		UserID:  userID,
		Type:    req.Type,
	})
	if err != nil {
		fmt.Printf("Failed to add reactions: %s\n", err)
		ResponseWithError(w, 500, "Failed to add reaction. Try again")
		return
	}

	user, err := apicfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error fetching user: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}

	ResponseWithJSON(w, 201, Reaction{
		ChirpID:   react.ChirpID,
		UserID:    react.UserID,
		Type:      react.Type,
		CreatedAt: react.CreatedAt,
		Username:  user.Username,
	})

	// fetch for chirp's author ID
	chirp, err := apicfg.DB.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		fmt.Printf("Failed to fetch chirp for notification: %s\n", err)
		return
	}

	// publish notification to rabbit
	// only notify if reactor is not chirp author
	if userID != chirp.UserID {
		// open channel
		ch, err := apicfg.MQConn.Channel()
		if err != nil {
			fmt.Printf("Failed to open MQ channel: %s\n", err)
			return
		}
		defer ch.Close()

		// publish to exchange "direct_notifications"
		err = pubsub.PublishJSON(r.Context(), ch, "direct_notifications", "", pubsub.DirectEvent{
			Type:      "reaction",
			Body:      fmt.Sprintf("%s reacted %s to your chirp.", user.Username, react.Type),
			Triggerer: userID,
			Username:  user.Username,
			Receiver:  chirp.UserID,
			ChirpID:   &chirpID,
		})
		if err != nil {
			fmt.Printf("Failed to publish reaction notification to exchange: %s\n", err)
			return
		}
		fmt.Printf("Reaction notification is published: %s\n", react.ChirpID)
	}
}
