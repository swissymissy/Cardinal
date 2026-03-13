package handler 

import (
	"fmt"
	"net/http"
	"errors"
	"database/sql"

	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	// decode request
	var user NewUser
	err := DecodeRequest(r, &user)
	if err != nil {
		fmt.Printf("Error decoding request: %s", err)
		ResponseWithError(w, 500 , "Something went wrong. Try again")
		return
	}

	email := user.Email
	passwd := user.Password

	// find user in database
	userInfo, err := apicfg.DB.GetUserByEmail(r.Context(), email)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Can't find user in db: %s", err)
		ResponseWithError(w, 401, "Incorrect Email or Password" )
		return
	} else if err != nil {
		fmt.Printf("Error getting user from db: %s", err)
		ResponseWithError(w, 401, "Incorrect Email or Password")
		return
	}

	// check password
	match, err := auth.CheckPasswordHash(passwd, userInfo.HashedPassword)
	if err != nil {
		fmt.Printf("%s", err)
		ResponseWithError(w, 401, "Incorrect Email or Password")
		return
	}
	if !match {
		ResponseWithError(w, 401 , "Incorrect Email or Password")
		return
	}

	// create new access token for user
	accessToken, err := auth.MakeJWT(userInfo.ID, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Error making new access token: %s", err)
		ResponseWithError(w, 500, "Something went wrong! Try again")
		return
	}
	// create a new refresh token for user
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Printf("Error making new refresh token: %s\n", err)
		respondWithError(w, 400, "Something went wrong")
		return
	}
	// store refresh token in database
	_, err = apicfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:	refreshToken,
		UserID: userInfo.ID,
	})
	if err != nil {
		fmt.Printf("Error storing refresh token in db: %s", err)
		ResponseWithError(w, 500, "Something went wrong. Try again")
		return
	}
	
	fmt.Printf("User %s has logged in", email)
	ResponseWithJSON(w, 200, LoginUser{
		ID: userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email: userInfo.Email,
		AccessToken: accessToken,
	})

} 