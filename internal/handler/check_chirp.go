package handler

import (
	"fmt"
	"strings"
)

// check for length and language in chirp 
func CheckChirp(msg *Chirp) error {

	// check length
	if len(msg.Body) > 500 {
		return fmt.Errorf("Chirp is too long")
	}

	msg.Body = cleanChirp(msg.Body)
	return nil
}

// clean bad words
func cleanChirp(msg string) string {
	split_msg := strings.Fields(msg)

	// bad words
	bad := map[string]struct{}{
		"fuck": {},
		"asshole": {},
		"dick": {},
		"dickhead": {},
		"shithead": {},
		"bitch": {},
	}

	for i := range split_msg{
		word := strings.ToLower(split_msg[i])
		if _, ok := bad[word]; ok {
			split_msg[i] = "meowmeow"
		}
	}
	msg = strings.Join(split_msg, " ")
	return msg
}