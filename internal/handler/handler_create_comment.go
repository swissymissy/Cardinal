package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerCreateComment(w http.ResponseWriter, r *http.Request) {
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

	// get chirp ID from url
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("Failed to parse chirp ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// decode body from request
	var commentReq Body
	err = DecodeRequest(r, &commentReq)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err)
		ResponseWithError(w, 400, "Invalid request")
		return
	}

	// check comment body
	err = CheckChirp(&commentReq)
	if err != nil {
		fmt.Printf("%s\n", err)
		ResponseWithError(w, 400, err.Error())
		return
	}

	// check chirp existence
	_, err = apicfg.DB.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 404, "Chirp not found")
			return
		}
		fmt.Printf("Error fetching chirp: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}

	// insert comment in db
	comment, err := apicfg.DB.CreateComment(r.Context(), database.CreateCommentParams{
		ChirpID: chirpID,
		UserID:  userID,
		Body:    commentReq.Body,
	})
	if err != nil {
		fmt.Printf("Error adding new comment to db: %s\n", err)
		ResponseWithError(w, 500, "Failed to create comment.")
		return
	}

	// get username
	user, err := apicfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error fetching user: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong. Try again.")
		return
	}
	ResponseWithJSON(w, 201, Comment{
		ID:        comment.ID,
		ChirpID:   comment.ChirpID,
		UserID:    comment.UserID,
		Username:  user.Username,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}
