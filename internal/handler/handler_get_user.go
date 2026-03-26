package handler

import (
	"net/http"
	"fmt"
	"errors"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerGetUser (w http.ResponseWriter, r *http.Request) {
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
	// get user's ID 
	targetIDStr := r.PathValue("userID")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		fmt.Printf("Invalid user ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid user ID")
		return
	}
	// retrieve user's information by userID
	targetInfo, err := apicfg.DB.GetUserByID(r.Context(), targetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 404, "User not found")
			return
		}
		fmt.Printf("Error fetching user: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong. Try again.")
		return
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
		ResponseWithJSON(w, 200, struct{
			ID			uuid.UUID	`json:"id"`
			CreatedAt	time.Time	`json:"created_at"`
			UpdatedAt	time.Time	`json:"updated_at"`
			Email		string		`json:"email"`
			FollowerCount int64		`json:"followers_count"`
			FollowingCount int64	`json:"followings_count"`
		}{
			ID: userID,
			CreatedAt: targetInfo.CreatedAt,
			UpdatedAt: targetInfo.UpdatedAt,
			Email: targetInfo.Email,
			FollowerCount: numFollowers,
			FollowingCount: numFollowings,
		} )
	} else {
		ResponseWithJSON(w , 200, UserProfile{
			ID:	targetInfo.ID,
			CreatedAt: targetInfo.CreatedAt,
			FollowerCount: numFollowers,
			FollowingCount: numFollowings,
		})
	}
}