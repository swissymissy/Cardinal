package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/swissymissy/Cardinal/internal/auth"
	"github.com/swissymissy/Cardinal/internal/database"
)

// struct to hold all query results
type commentQueryResults struct {
	comments    []database.GetCommentsByChirpIDRow
	counts      int64
	chirpErr    error
	commentsErr error
	countsErr   error
}

func (apicfg *ApiConfig) HandlerGetComments(w http.ResponseWriter, r *http.Request) {
	// auth check
	// check user's token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate token
	_, err = auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// get chirp ID from URL
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("Failed to parse chirp ID: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// run queries concurrently
	var wg sync.WaitGroup
	results := commentQueryResults{}

	// check chirp existence
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, results.chirpErr = apicfg.DB.GetOneChirp(r.Context(), chirpID)

	}()

	// get comments from db
	wg.Add(1)
	go func() {
		defer wg.Done()
		results.comments, results.commentsErr = apicfg.DB.GetCommentsByChirpID(r.Context(), chirpID)

	}()

	// get total number of comments
	wg.Add(1)
	go func() {
		defer wg.Done()
		results.counts, results.countsErr = apicfg.DB.GetCommentCount(r.Context(), chirpID)

	}()

	// wait for goroutines to be done
	wg.Wait()

	// check errors
	if results.chirpErr != nil {
		if errors.Is(results.chirpErr, sql.ErrNoRows) {
			ResponseWithError(w, 404, "Chirp not found")
			return
		}
		fmt.Printf("Error fetching chirp: %s\n", results.chirpErr)
		ResponseWithError(w, 500, "Something went wrong")
		return
	}
	if results.commentsErr != nil {
		fmt.Printf("Error getting comments from db: %s\n", results.commentsErr)
		ResponseWithError(w, 500, "Something went wrong. Try again")
		return
	}
	if results.countsErr != nil {
		fmt.Printf("Error getting total comments: %s\n", results.countsErr)
		ResponseWithError(w, 500, "Something went wrong. Try again")
		return
	}

	// writing comment to response format
	list := []Comment{}
	for _, c := range results.comments {
		list = append(list, Comment{
			ID:        c.ID,
			ChirpID:   c.ChirpID,
			UserID:    c.UserID,
			Username:  c.Username,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	ResponseWithJSON(w, http.StatusOK, CommentSummary{
		Total:    results.counts,
		Comments: list,
	})
}
