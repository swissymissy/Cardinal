package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (apicfg *ApiConfig) HandlerGetOneChirp(w http.ResponseWriter, r *http.Request) {
	// get chirp ID
	chirpIDStr := r.PathValue("chirpsID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	// retrieve chirp from table
	chirpInfo, err := apicfg.DB.GetOneChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Error getting row from chirps table: %s\n", err)
		ResponseWithError(w, 404, "Chirp does not exist or deleted")
		return
	} else if err != nil {
		fmt.Printf("Error getting row from chirps table: %s\n", err)
		ResponseWithError(w, 500, "Unable to get chirp. Try again")
		return
	}

	ResponseWithJSON(w, 200, CreatedChirp{
		ID:        chirpInfo.ID,
		CreatedAt: chirpInfo.CreatedAt,
		UpdatedAt: chirpInfo.UpdatedAt,
		Body:      chirpInfo.Body,
		UserID:    chirpInfo.UserID,
	})
}
