package handler

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
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
	userPassword := newUser.Password

	// hash password
	hashed, err := auth.HashPassword(userPassword)
	if err != nil {
		fmt.Printf("Error hashing password: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}

	// save user to database
	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: userEmail,
		HashedPassword: hashed,
	})
	if err != nil {
		fmt.Printf("Error inserting new user: %s\n", err)
		ResponseWithJSON(w, 500 , "Can't create new user! Try again")
		return
	}

	// response back to client
	ResponseWithJSON(w, 201, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
	return

}