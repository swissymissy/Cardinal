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
type reactionQueryResults struct {
	reactors     []database.GetReactionsByChirpIDRow
	counts       []database.GetReactionCountsRow
	total        int64
	userReaction string
	reactorsErr  error
	countsErr    error
	totalErr     error
}

func (apicfg *ApiConfig) HandlerGetReactions(w http.ResponseWriter, r *http.Request) {
	// auth check
	// check user's token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}
	// validate token
	userID, err := auth.ValidateJWT(accessToken, apicfg.JWTSecret)
	if err != nil {
		fmt.Printf("Invalid token: %s\n", err)
		ResponseWithError(w, 401, "Invalid Token")
		return
	}

	// get chirp ID from URL
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		fmt.Printf("Error parsing ID from URL: %s\n", err)
		ResponseWithError(w, 400, "Invalid ID")
		return
	}

	// run all queries concurrently
	var wg sync.WaitGroup
	results := reactionQueryResults{}

	// query 1: get list of people who reacts to the chirp
	wg.Add(1)
	go func() {
		defer wg.Done()
		results.reactors, results.reactorsErr = apicfg.DB.GetReactionsByChirpID(r.Context(), chirpID)

	}()

	// query 2: get number of each reaction on the chirp
	wg.Add(1)
	go func() {
		defer wg.Done()
		results.counts, results.countsErr = apicfg.DB.GetReactionCounts(r.Context(), chirpID)
	}()

	// query 3: get current user's reaction on the chirp
	wg.Add(1)
	go func() {
		defer wg.Done()
		reaction, err := apicfg.DB.GetUserReactions(r.Context(), database.GetUserReactionsParams{
			ChirpID: chirpID,
			UserID:  userID,
		})
		if err == nil {
			results.userReaction = reaction
		} else if !errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Error getting current user's reaction on chirp: %s\n", err)
		}
	}()

	// query 4: get total amount of reactions on the chirp
	wg.Add(1)
	go func() {
		defer wg.Done()
		results.total, results.totalErr = apicfg.DB.GetReactionTotal(r.Context(), chirpID)
	}()

	// wait for all queries to complete
	wg.Wait()

	// check errors after all goroutines finishes
	if results.reactorsErr != nil {
		fmt.Printf("Error getting reactors: %s\n", results.reactorsErr)
		ResponseWithError(w, 500, "Failed to get reactions")
		return
	}
	if results.countsErr != nil {
		fmt.Printf("Error getting counts: %s\n", results.countsErr)
		ResponseWithError(w, 500, "Failed to get reaction counts")
		return
	}
	if results.totalErr != nil {
		fmt.Printf("Error getting total reactions: %s\n", results.totalErr)
		ResponseWithError(w, 500, "Failed to get total reactions")
		return
	}

	// writing each reactor to response format
	reactorsList := []Reaction{}
	for _, r := range results.reactors {
		reactorsList = append(reactorsList, Reaction{
			ChirpID:   r.ChirpID,
			UserID:    r.UserID,
			Type:      r.Type,
			CreatedAt: r.CreatedAt,
			Username:  r.Username,
		})
	}

	// writing each count type to response fmt
	countList := []ReactionCount{}
	for _, c := range results.counts {
		countList = append(countList, ReactionCount{
			Type:  c.Type,
			Count: c.Count,
		})
	}

	ResponseWithJSON(w, 200, ReactionSummary{
		Counts:       countList,
		Reactors:     reactorsList,
		UserReaction: results.userReaction,
		Total:        results.total,
	})
}
