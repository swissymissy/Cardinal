package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// extract token from header
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error extracting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Error validating token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	chirpIDStr := r.PathValue("chirpsID")  // extract ID string from URL
	chirpID, err := uuid.Parse(chirpIDStr) // convert string to uuid
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	// retrieve chirp form table
	chirpInfo, err := apicfg.DB.GetOneChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		ResponseWithError(w, 404, "Chirp does not exist or already deleted")
		return
	} else if err != nil {
		ResponseWithError(w, 500, "Can't get chirp")
		return
	}

	// check if chirp belong to user
	if userID != chirpInfo.UserID {
		ResponseWithError(w, 403, "Unauthorized")
		return
	}

	// delete chirp
	err = apicfg.DB.DeleteOneChirp(r.Context(), chirpID)
	if err != nil {
		ResponseWithError(w, 500, "Unable to delete chirp")
		return
	}
	w.WriteHeader(204)
}
