package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeRequest[T any](r *http.Request, val *T) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(val); err != nil {
		return fmt.Errorf("Error decoding request: %w", err)
	}
	return nil
}
