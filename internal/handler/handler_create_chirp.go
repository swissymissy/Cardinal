package handler

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/swissymissy/Cardinal/internal/auth"
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

	// check user's token 
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s", err)
		ResponseWithError(w, 401, "Invalid Token")
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
	ResponseWithJSON(w, 201, CreatedChirp{
		ID: createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body: createdChirp.Body,
		UserID: createdChirp. UserID,
	})

}