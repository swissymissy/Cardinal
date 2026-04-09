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

func (apicfg *ApiConfig) HandlerDeleteComment(w http.ResponseWriter, r *http.Request) {
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
	cmtIDStr := r.PathValue("commentID")
	cmtID, err := uuid.Parse(cmtIDStr)
	if err != nil {
		fmt.Printf("Failed to parse comment ID from URL: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}
	
	// delete comment
	// authorization check DB level
	_, err = apicfg.DB.DeleteComment(r.Context(), database.DeleteCommentParams{
		ID:     cmtID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 404, "Comment not found")
			return
		}
		fmt.Printf("Error deleting comment: %s\n", err)
		ResponseWithError(w, 500, "Failed to remove comment")
		return
	}

	ResponseWithJSON(w, 200, struct {
		Message string `json:"message"`
	}{
		Message: "Comment removed",
	})
}
