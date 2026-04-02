package handler

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (apicfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// check if there is query parameter in URL
	authorID := r.URL.Query().Get("author_id")
	sortPara := r.URL.Query().Get("sort")
	desc := false
	if sortPara == "desc" {
		desc = true
	}

	var list []CreatedChirp

	if authorID == "" {
		rows, err := apicfg.DB.GetAllChirps(r.Context())
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
		rows, err := apicfg.DB.GetAllChirpsFromUserID(r.Context(), userID)
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

	// desc order
	if desc {
		sort.Slice(list, func(i, j int) bool {
			return list[i].CreatedAt.After(list[j].CreatedAt)
		})
	}

	ResponseWithJSON(w, http.StatusOK, list)
}
