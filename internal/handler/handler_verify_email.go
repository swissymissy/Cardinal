package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerVerifyEmail(w http.ResponseWriter, r *http.Request) {
	tokenIDStr := r.URL.Query().Get("t")
	if token == "" {
		http.Redirect(w, r, "/verified.html?error=invalid", http.StatusSeeOther)
		return
	}

	// convert token string to UUID
	tokenID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		http.Redirect(w, r, "/verified.html?error=invalid", http.StatusSeeOther)
		return
	}

	// look up token in DB
	user, err := apicfg.DB.GetUserByVerificationToken(r.Context(), tokenID)
	if err != nil {
		http.Redirect(w, r, "/verified.html?error=expired", http.StatusSeeOther)
		return
	}

	// mark user verified
	if err = apicfg.DB.MarkUserVerified(r.Context(), user.ID); err != nil {
		http.Redirect(w, r, "/verified.html?error=invalid", http.StatusSeeOther)
		return
	}

	apicfg.DB.DeleteVerificationToken(r.Context(), tokenID)    // delete token after used
	http.Redirect(w, r, "/verified.html", http.StatusSeeOther) // redirect to success page
}
