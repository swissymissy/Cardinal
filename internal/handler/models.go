package handler

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`json:"id"`			
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

type NewUser struct {
	Email		string 		`json:"email"`
	Password 	string      `json:"password"`
}

type LoginUser struct {
	ID 			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
	AccessToken string  	`json:"access_token"`
}

type Chirp struct {
	Body	string		`json:"body"`
	UserID  uuid.UUID	`json:"user_id"`
}

type CreatedChirp struct {
	ID			uuid.UUID		`json:"id"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
	Body 		string			`json:"body"`
	UserID		uuid.UUID		`json:"user_id"`
}


type ResponseAccessToken struct {
	Token 	string 	`json:"token"`
}