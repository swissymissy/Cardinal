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
	userName := newUser.Username 

	// hash password
	hashed, err := auth.HashPassword(userPassword)
	if err != nil {
		fmt.Printf("Error hashing password: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}

	// save user to database
	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          userEmail,
		HashedPassword: hashed,
		Username:		userName,
	})
	if err != nil {
		fmt.Printf("Error inserting new user: %s\n", err)
		ResponseWithError(w, 400, "Can't create new user! Email or Username already taken.")
		return
	}

	// response back to client
	ResponseWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Username:  user.Username,
	})
	return

}
