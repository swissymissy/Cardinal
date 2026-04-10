package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

type FeedRequest struct {
	Before *time.Time `json:"before"`
}

func (apicfg *ApiConfig) HandlerGetFeed(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		ResponseWithError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		ResponseWithError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	// decode request
	var req FeedRequest
	DecodeRequest(r, &req)

	before := time.Now() // default is at Now
	if req.Before != nil {
		before = *req.Before
	}

	// get chirps
	feedChirps, err := apicfg.DB.GetFeedChirps(r.Context(), database.GetFeedChirpsParams{
		UserID:    userID,
		CreatedAt: before,
		Limit:     20,
	})
	if err != nil {
		fmt.Printf("Error getting feed: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "Can't get feed. Try again.")
		return
	}

	list := []CreatedChirp{}
	for _, c := range feedChirps {
		list = append(list, CreatedChirp{
			ID:            c.ID,
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
			Body:          c.Body,
			UserID:        c.UserID,
			Username:      c.Username,
			ReactionCount: c.ReactionCount,
			CommentCount:  c.CommentCount,
		})
	}
	ResponseWithJSON(w, http.StatusOK, list)
}
