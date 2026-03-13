package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

// create new access token for user
func MakeJWT(userID uuid.UUID, serverSecretToken string, expiresIn time.Duration ) (string, error) {

	//create a new registered claim
	claim := jwt.RegisteredClaims{
		Issuer: "Cardinal-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}
	//create new token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	//sign new token with server secret key
	signKey := []byte(serverSecretToken)
	signedToken, err := newToken.SignedString(signKey)
	if err != nil {
		return "", fmt.Errorf("Can't sign token with key : %w", err)
	}
	return signedToken, nil
}

// check user's token
func ValidateJWT(tokenString, serverSecretToken string) (uuid.UUID, error) {
	// create new empty claim struct to be filled
	claim := &jwt.RegisteredClaims{}
	// pass a pointer to that struct so the library can modify it
	_, err := jwt.ParseWithClaims(
		tokenString,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(serverSecretToken), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Token is expired or bad signature: %w", err)
	}

	// retrieve userID from claim's Subject field
	userIDStr := claim.Subject 
	// convert userID to uuid type
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error converting string to uuid: %w", err)
	}
	return userID, nil
}