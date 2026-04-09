package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
)

func (apicfg *ApiConfig) HandlerGetNotifications(w http.ResponseWriter, r *http.Request) {
	// get user's token
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

	// get notifications
	notifications, err := apicfg.DB.GetNotificationByReceiver(r.Context(), userID)
	if err != nil {
		fmt.Printf("Error getting notifications: %s\n", err)
		ResponseWithError(w, 500, "Can't get notifications")
		return
	}
	notiList := []Notification{}
	for _, noti := range notifications {
		var chirpID *uuid.UUID
		if noti.chirpID.Valid { // check if chirpID is not null
			chirpID = &noti.ChirpID.UUID
		}
		notiList = append(notiList, Notification{
			ID:        noti.ID,
			CreatedAt: noti.CreatedAt,
			Body:      noti.Body,
			Receiver:  noti.Receiver,
			Username:  noti.Username,
			ChirpID:   chirpID,
			IsRead:    noti.IsRead,
		})
	}
	ResponseWithJSON(w, http.StatusOK, notiList)
}
