package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

func (apicfg *ApiConfig) HandlerMarkAllRead(w http.ResponseWriter, r *http.Request) {
	// get user token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate user token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// mark all notification read
	err = apicfg.DB.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		fmt.Printf("Can't mark all notif as read: %s\n", err)
		ResponseWithError(w, 500, "Can't mark all as read. Try again")
		return
	}
	w.WriteHeader(200)
}

func (apicfg *ApiConfig) HandlerMarkOneRead(w http.ResponseWriter, r *http.Request) {
	// get user token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate user token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	notifIDStr := r.PathValue("notifID")
	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		fmt.Printf("Can't parse string to uuid")
		ResponseWithError(w, http.StatusBadRequest, "Invalid notif ID")
		return
	}
	// mark one notification as read
	err = apicfg.DB.MarkOneAsRead(r.Context(), database.MarkOneAsReadParams{
		ID:       notifID,
		Receiver: userID,
	})
	if err != nil {
		fmt.Printf("Can't mark notification as read: %s\n", err)
		ResponseWithError(w, 500, "Can't mark as read. Try again.")
		return
	}
	w.WriteHeader(200)
}
