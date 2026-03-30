package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerGetFollowings(w http.ResponseWriter, r *http.Request) {
	// get user token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate user token
	_, err = auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// get target ID from URL
	targetIDStr := r.PathValue("userID")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		fmt.Printf("Invalid user ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid user ID")
		return
	}
	// retrieve followings list from DB
	var followingList []database.GetFollowingsRow
	followingList, err = apicfg.DB.GetFollowings(r.Context(), targetID)
	if err != nil {
		fmt.Printf("Error getting following list: %s\n", err)
		ResponseWithError(w, 500, "Failed to get followings. Try again.")
		return
	}
	// put each following into format response
	list := []FollowList{}
	for _, followee := range followingList {
		list = append(list, FollowList{
			UserID:    followee.FolloweeID,
			CreatedAt: followee.CreatedAt,
		})
	}
	// response the list
	ResponseWithJSON(w, http.StatusOK, list)
}
