package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerGetFollowers(w http.ResponseWriter, r *http.Request) {
	// extract user ID from URL
	targetIDStr := r.PathValue("userID")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		fmt.Printf("Invalid user ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid user ID")
		return
	}
	// get the followers list of the target user ID
	var followerList []database.GetFollowersRow
	followerList, err = apicfg.DB.GetFollowers(r.Context(), targetID)
	if err != nil {
		fmt.Printf("Error getting follower list: %s\n", err)
		ResponseWithError(w, 500, "Failed to get followers. Try again.")
		return
	}
	// writing each follower to the response format
	list := []FollowList{}
	for _, follower := range followerList {
		list = append(list, FollowList{
			UserID:    follower.FollowerID,
			Username:  follower.Username,
			CreatedAt: follower.CreatedAt,
		})
	}
	// response with the list
	ResponseWithJSON(w, http.StatusOK, list)
}
