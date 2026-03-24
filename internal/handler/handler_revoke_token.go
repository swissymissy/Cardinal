package handler 

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
)

// revoke the expired refresh token
func (apicfg *ApiConfig) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	// check token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s", err)
		ResponseWithError(w, 401, "Invalid token")
		return
	}

	// look for user with the given token in db
	user, err := apicfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Printf("Can't find user: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// revoke token
	refreshTokenDB := user.Token
	err = apicfg.DB.RevokedToken(r.Context(), refreshTokenDB)
	if err != nil {
		fmt.Printf("Error revoking refresh token: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}
	w.WriteHeader(204)
}