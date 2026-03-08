package handler

import (
	"fmt"
	"net/http"
)

// handle create new user
func (apicfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	
	// decode the request 
	var newUser NewUser
	err := DecodeRequest(r, &newUser)
	if err != nil {
		fmt.Printf("Error decoding request: %s\n", err)
		ResponseWithError(w, 500, "Can't create new user")
		return
	}
	userEmail := newUser.Email

	// save user to database
	user, err := apicfg.DB.CreateUser(r.Context(), userEmail)
	if err != nil {
		fmt.Printf("Error inserting new user: %s", err)
		ResponseWithJSON(w, 500 , "Can't create new user! Try again")
		return
	}

	// response back to client
	ResponseWithJSON(w, 200, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
	return

}