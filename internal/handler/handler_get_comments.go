package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerGetComments(w http.ResponseWriter, r *http.Request) {
	// auth check
	// check user's token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate token
	_, err = auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// get chirp ID from URL
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("Failed to parse chirp ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// check chirp existence
	_, err = apicfg.DB.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		ResponseWithError(w, 404, "Chirp not found")
		return
	}
	
	// get comments from db
	commentList, err := apicfg.DB.GetCommentsByChirpID(r.Context(), chirpID)
	if err != nil {
		fmt.Printf("Error getting comments from db: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong. Try again")
		return
	}

	// writing comment to response format
	list := []Comment{}
	for _, c := range commentList {
		list = append(list, Comment{
			ID:        c.ID,
			ChirpID:   c.ChirpID,
			UserID:    c.UserID,
			Username:  c.Username,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	ResponseWithJSON(w, http.StatusOK, list)
}
