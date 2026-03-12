package handler

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/database"
)


// create new chirp
func (apicfg *ApiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// decode request
	var newChirp Chirp 
	err := DecodeRequest(r, &newChirp)
	if err != nil {
		fmt.Printf("Error decoding request: %s")
		msg := "Something went wrong"
		ResponseWithError(w, 500, msg )
		return
	}

	chirpBody := newChirp.Body
	author := newChirp.UserID

	err = CheckChirp(&newChirp) 
	if err != nil {
		fmt.Printf("%s", err)
		ResponseWithError(w, 400, "Chirp is too long")
		return
	}
	// save new chirp to db
	createdChirp , err := apicfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: chirpBody,
		UserID: author,
	})
	if err != nil {
		fmt.Printf("Error adding new chirp to db: %s\n", err)
		msg := "Can't create chirp"
		ResponseWithError(w, 500, msg)
		return 
	}
	ResponseWithJSON(w, 200, CreatedChirp{
		ID: createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body: createdChirp.Body,
		UserID: createdChirp. UserID,
	})

}