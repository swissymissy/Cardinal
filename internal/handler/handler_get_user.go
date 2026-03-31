package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	// get user's token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate user's token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// get user's identifier
	target := r.PathValue("identifier")
	var targetID uuid.UUID
	var targetInfo database.User 

	if id, err := uuid.Parse(target); err == nil {
		targetID = id
		targetInfo, err = apicfg.DB.GetUserByID(r.Context(), targetID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ResponseWithError(w, 404, "User not found")
				return
			}
			fmt.Printf("Error fetching user: %s\n", err)
			ResponseWithError(w, 500, "Can't find user. Try again.")
			return
		}
	} else {
		targetInfo, err = apicfg.DB.GetUserByUsername(r.Context(), target)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ResponseWithError(w, 404, "User not found")
				return
			}
			fmt.Printf("Error fetching user: %s\n", err)
			ResponseWithError(w, 500, "Can't find user. Try again")
			return
		}
		targetID = targetInfo.ID
	}
	
	// retrieve user's number of followers
	numFollowers, err := apicfg.DB.GetCountFollowers(r.Context(), targetID)
	if err != nil {
		fmt.Printf("Error fetching follower count: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}
	// retrieve user's number of followings
	numFollowings, err := apicfg.DB.GetCountFollowings(r.Context(), targetID)
	if err != nil {
		fmt.Printf("Error fetching follower count: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}
	// repsonse with information
	if targetID == userID {
		ResponseWithJSON(w, 200, struct {
			ID             uuid.UUID `json:"id"`
			Username 	   string    `json:"username"`
			CreatedAt      time.Time `json:"created_at"`
			UpdatedAt      time.Time `json:"updated_at"`
			Email          string    `json:"email"`
			FollowerCount  int64     `json:"followers_count"`
			FollowingCount int64     `json:"followings_count"`
		}{
			ID:             userID,
			Username:		targetInfo.Username,
			CreatedAt:      targetInfo.CreatedAt,
			UpdatedAt:      targetInfo.UpdatedAt,
			Email:          targetInfo.Email,
			FollowerCount:  numFollowers,
			FollowingCount: numFollowings,
		})
	} else {
		ResponseWithJSON(w, 200, UserProfile{
			ID:             targetInfo.ID,
			Username:		targetInfo.Username,
			CreatedAt:      targetInfo.CreatedAt,
			FollowerCount:  numFollowers,
			FollowingCount: numFollowings,
		})
	}
}
