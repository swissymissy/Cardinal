package handler

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerUnfollow(w http.ResponseWriter, r *http.Request) {
	// get user's access token
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
	// decode request
	var unfollow NewFollow
	err = DecodeRequest(r, &unfollow)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err)
		msg := "Something went wrong"
		ResponseWithError(w, 500, msg)
		return
	}
	followerID := userID
	followeeID := unfollow.FolloweeID
	// delete following connection in db
	err = apicfg.DB.UnfollowUser(r.Context(), database.UnfollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	if err != nil {
		fmt.Printf("Error unfollowing user: %s\n", err)
		ResponseWithError(w, 500, "Failed to unfollow. Try again.")
		return
	}

	ResponseWithJSON(w, 200, struct {
		Message string `json:"message"`
	}{
		Message: "Successfully unfollowed",
	})
}
