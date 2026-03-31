package handler

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/swissymissy/Cardinal/internal/pubsub"
)

// create new chirp
func (apicfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// decode request
	var newChirp Chirp
	err := DecodeRequest(r, &newChirp)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err)
		msg := "Something went wrong"
		ResponseWithError(w, 500, msg)
		return
	}

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

	chirpBody := newChirp.Body
	author := userID

	err = CheckChirp(&newChirp)
	if err != nil {
		fmt.Printf("%s\n", err)
		ResponseWithError(w, 400, "Chirp is too long")
		return
	}
	// save new chirp to db
	createdChirp, err := apicfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   chirpBody,
		UserID: author,
	})
	if err != nil {
		fmt.Printf("Error adding new chirp to db: %s\n", err)
		msg := "Can't create chirp"
		ResponseWithError(w, 500, msg)
		return
	}
	// get username 
	user, err := apicfg.DB.GetUserByID(r.Context(), createdChirp.UserID)
	if err != nil {
		fmt.Printf("Error fetching user: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong. Try again.")
		return
	}

	ResponseWithJSON(w, 201, CreatedChirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
		Username:  user.Username,
	})

	// publish notification to rabbit
	// open channel from connection
	ch, err := apicfg.MQConn.Channel()
	if err != nil {
		fmt.Printf("Failed to open MQ channel: %s\n", err)
		return
	}
	defer ch.Close()

	// publish to exchange, using fanout exchange type
	err = pubsub.PublishJSON(r.Context(), ch, "notifications", "", pubsub.ChirpEvent{
		Body:      createdChirp.Body,
		Triggerer: createdChirp.UserID,
		Username:  user.Username,
		ChirpID:   createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
	})
	if err != nil {
		fmt.Printf("Failed to publish notification to exchange: %s\n", err)
		return
	}
	fmt.Printf("Notification is published for new chirp: %s\n", createdChirp.ID)
}
