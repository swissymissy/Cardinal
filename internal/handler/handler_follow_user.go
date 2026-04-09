package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerFollowUser(w http.ResponseWriter, r *http.Request) {
	// get user token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate user token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// decode request to get follower and followee ID
	var newFollow NewFollow
	err = DecodeRequest(r, &newFollow)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err)
		msg := "Something went wrong"
		ResponseWithError(w, 500, msg)
		return
	}
	followerID := userID
	followeeID := newFollow.FolloweeID
	// create follow connection in db
	newFollowing, err := apicfg.DB.FollowUser(r.Context(), database.FollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 409, "Already following this user")
			return
		}
		fmt.Printf("Error creating new follower: %s\n", err)
		ResponseWithError(w, 500, "Failed to follow")
		return
	}

	ResponseWithJSON(w, 201, Follower{
		FollowerID: newFollowing.FollowerID,
		FolloweeID: newFollowing.FolloweeID,
		CreatedAt:  newFollowing.CreatedAt,
		UpdatedAt:  newFollowing.UpdatedAt,
	})
}
