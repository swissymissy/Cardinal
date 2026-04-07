package handler

import (
	"fmt"
	"net/http"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerRemoveReaction(w http.ResponseWriter, r *http.Request) {
	//auth check
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

	// get chirpID from URL
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("Error parsing chirp ID from URL: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// remove reaction
	_, err = apicfg.DB.RemoveReaction(r.Context(), database.RemoveReactionParams{
		ChirpID: chirpID,
		UserID:  userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 404, "Reaction does not exist")
			return
		}
		fmt.Printf("Error deleting reaction: %s\n", err)
		ResponseWithError(w, 500, "Failed to remove reaction.")
		return
	}

	ResponseWithJSON(w, 200, struct {
		Message string `json:"message"`
	}{
		Message: "Reaction removed",
	})
}
