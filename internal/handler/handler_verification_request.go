package handler

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerRequestVerification(w http.ResponseWriter, r *http.Request) {
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

	// generate new token
	token, err := apicfg.DB.CreateVerificationToken(r.Context(), userID)
	if err != nil {
		fmt.Printf("Failed to create new token: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong.Try again")
		return
	}

}
