package handler

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

func (apicfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// check if there is query parameter in URL
	authorID := r.URL.Query().Get("author_id")
	sortPara := r.URL.Query().Get("sort")
	beforeStr := r.URL.Query().Get("bf")
	before := time.Now()
	if beforeStr != "" {
		parsed, err := time.Parse(time.RFC3339, beforeStr)
		if err == nil {
			before = parsed
		}
	}

	asc := false
	if sortPara == "asc" {
		asc = true
	}

	list := []CreatedChirp{}

	if authorID == "" {
		rows, err := apicfg.DB.GetAllChirps(r.Context(), before)
		if err != nil {
			fmt.Printf("Error getting all chirps: %s\n", err)
			ResponseWithError(w, http.StatusInternalServerError, "Can't get chirps. Try again")
			return
		}
		for _, c := range rows {
			list = append(list, CreatedChirp{
				ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
				Body: c.Body, UserID: c.UserID, Username: c.Username,
			})
		}
	} else {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		rows, err := apicfg.DB.GetAllChirpsFromUserID(r.Context(), database.GetAllChirpsFromUserIDParams{
			UserID:    userID,
			CreatedAt: before,
			Limit:     20,
		})
		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Can't get all chirps. Try again")
			return
		}
		for _, c := range rows {
			list = append(list, CreatedChirp{
				ID: c.ID, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
				Body: c.Body, UserID: c.UserID, Username: c.Username,
			})
		}
	}

	// asc order
	if asc {
		sort.Slice(list, func(i, j int) bool {
			return list[i].CreatedAt.Before(list[j].CreatedAt)
		})
	}

	ResponseWithJSON(w, http.StatusOK, list)
}
