package handler 

import (
	"net/http"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// check if there is query parameter in URL
	authorID := r.URL.Query().Get("author_id")
	sortPara := r.URL.Query().Get("sort")
	desc := false
	if sortPara == "desc" {
		desc = true
	} 

	// fetch data according to query parameter
	var chirpList []database.Chirp
	var err error

	if authorID == "" {
		// author id is not given
		chirpList, err = apicfg.DB.GetAllChirps( r.Context())
		if err != nil {
			fmt.Printf("Error getting all chirps : %s", err)
			ResponseWithError(w, http.StatusInternalServerError, "Can't get chirps. Try again")
			return
		}
	} else {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			ResponseWithError(w, http.StatusBadRequest , "Invalid ID")
			return
		}

		chirpList, err = apicfg.DB.GetAllChirpsFromUserID(r.Context(), userID)
		if err != nil {
			ResponseWithError(w, http.StatusInternalServerError, "Can't get all chirps. Try again")
			return
		}
	}
	
	// desc order
	if desc {
		sort.Slice(chirpList, func(i ,j int) bool{
			return chirpList[i].CreatedAt.After(chirpList[j].CreatedAt)
		})
	}

	// writing each chirp to response
	var list []CreatedChirp
	for _, chirp := range chirpList{
		list = append(list, CreatedChirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	ResponseWithJSON(w, http.StatusOK, list)
}