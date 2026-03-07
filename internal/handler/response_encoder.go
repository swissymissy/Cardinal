package handler

import (
	"fmt"
	"net/http"
	"encoding/json"	
)

// successful response with json
func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header.Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// encode to json bytes
	bytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error encoding payload to json: %s", err)
		return
	}
	w.Write(data)
}

// write response with error message
func ResponseWithError( w http.ResponseWriter, code int, msg string) {
	type errorMsg struct {
		Error string  `json:"error"`
	}

	response := errorMsg{
		Error: msg,
	}
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error encoding msg to json: %s", err)
		return
	}
	w.Header.Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
