package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/wneessen/go-mail"
)

func (apicfg *ApiConfig) HandlerRequestVerification(w http.ResponseWriter, r *http.Request) {
	// get user's access token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate user's token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// generate new token
	token, err := apicfg.DB.CreateVerificationToken(r.Context(), userID)
	if err != nil {
		fmt.Printf("Failed to create new token: %s\n", err)
		ResponseWithError(w, 500, "Something went wrong.Try again")
		return
	}

	// get user's email
	user, err := apicfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, 404, "User not found")
			return
		}
		fmt.Printf("Error fetching user: %s\n", err)
		ResponseWithError(w, 500, "Can't find user. Try again")
		return
	}

	// verification link
	link := fmt.Sprintf("%s/api/verify?t=%s", apicfg.BaseURL, token.Token)

	// send email to user
	msg := mail.NewMsg()
	if err := msg.From(apicfg.SMTPUsername); err != nil {
		fmt.Printf("Failed to set From: %s\n", err)
		ResponseWithError(w, 500, "Failed to send verification email. Try again")
		return
	}
	if err := msg.To(user.Email); err != nil {
		fmt.Printf("Failed to set To: %s\n", err)
		ResponseWithError(w, 500, "Failed to send verification email. Try again")
		return
	}
	msg.Subject("Cardinal: Email verification")
	msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(
		"Hello %s,\n\nPlease click the link below to verify your email.\nThis link will expire in 30 days.\n\n%s\n\nIf you did not request this, please ignore this email.\nIf the link above has expired, please request for a new link.\n", user.Username, link,
	))

	client, err := mail.NewClient(apicfg.SMTPHost,
		mail.WithPort(apicfg.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(apicfg.SMTPUsername),
		mail.WithPassword(apicfg.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		fmt.Printf("Failed to create mail client: %s\n", err)
		ResponseWithError(w, 500, "Failed to send verification email. Try again")
		return
	}
	if err = client.DialAndSend(msg); err != nil {
		fmt.Printf("Failed to send verification email: %s\n", err)
		ResponseWithError(w, 500, "Failed to send email. Try again")
		return
	}
	ResponseWithJSON(w, 200, struct {
		Message string `json:"message"`
	}{
		Message: "Successfully sent verification email. Please check email",
	})
}
