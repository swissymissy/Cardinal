package handler 

import (
	"fmt"
	"net/http"
	"time"

	"github.com/swissymissy/Cardinal/internal/auth"
)

// checking if user's refresh token is expired/ revoked yet , then create a new access token for user
func (apicfg *ApiConfig) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// get token from request header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("error getting token from header: %s\n", err)
		ResponseWithError(w , 401, "Invalid token")
		return
	}

	// get refresh token from db
	refreshTokenDb, err := apicfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Printf("error getting user from db: %s\n", err)
		ResponseWithError(w , 401, "Invalid token")
		return
	}

	// check if token expires yet
	if refreshTokenDb.ExpiresAt.Before(time.Now()) {
		ResponseWithError(w, 401 , "Token has expired")
		return
	} 
	// check if token is revoked yet ( is not null)
	if refreshTokenDb.RevokedAt.Valid {
		ResponseWithError(w , 401, "Token has been revoked")
		return
	}

	// create new access token for user
	newAccessToken, err := auth.MakeJWT(refreshTokenDb.UserID, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Error making new access token: %s\n", err)
		ResponseWithError(w , 500, "Something went wrong. Try again")
		return
	}
	ResponseWithJSON(w, 200 , ResponseAccessToken{
		Token: newAccessToken,
	})
}
