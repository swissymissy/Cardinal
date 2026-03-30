package handler

import (
	"fmt"
	"net/http"
)

// reset table users in db
func (apicfg *ApiConfig) HandlerResetUsers(w http.ResponseWriter, r *http.Request) {
	if apicfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := apicfg.DB.ResetUsers(r.Context())
	if err != nil {
		fmt.Printf("Error resetting users table: %s", err)
		ResponseWithError(w, 500, "Can't reset table. Something went wrong")
		return
	}
	w.WriteHeader(200)
}
