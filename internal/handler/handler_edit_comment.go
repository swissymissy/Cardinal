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

func (apicfg *ApiConfig) HandlerEditComment(w http.ResponseWriter, r *http.Request) {
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

	// get comment ID from URL
	commentIDStr := r.PathValue("commentID")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		fmt.Printf("Failed to parse comment ID from URL: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// decode body request
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
		ResponseWithError(w, 400, "Comment is too long")
		return
	}

	// update comment body in db
	comment, err := apicfg.DB.EditComment(r.Context(), database.EditCommentParams{
		Body:   commentReq.Body,
		ID:     commentID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 404, "Comment not found")
			return
		}
		fmt.Printf("Error editing comment in db: %s\n", err)
		ResponseWithError(w, 500, "Failed to edit comment. Try again")
		return
	}
	// get username
	user, err := apicfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error fetching user: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong. Try again.")
		return
	}

	ResponseWithJSON(w, 200, Comment{
		ID:        comment.ID,
		ChirpID:   comment.ChirpID,
		UserID:    comment.UserID,
		Username:  user.Username,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}
