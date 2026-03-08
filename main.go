package main

import (
	"fmt"
	"os"
	"net/http"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/swissymissy/Cardinal/internal/handler"
	"github.com/swissymissy/Cardinal/internal/database"
)


func main() {

	// get values from .env 
	godotenv.Load() 				
	port := os.Getenv("PORT")					// load port 
	dbURL := os.Getenv("DB_URL")				// load db url

	// open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err)
		return
	}
	dbQuery := database.New(db)

	// create apiConfig
	apicfg := &handler.ApiConfig{
		DB: dbQuery,
		Port: port,
	}

	// server mux
	mux := http.NewServeMux()
	// server
	address := fmt.Sprintf(":%s", port)
	cardinalServer := http.Server{
		Addr: address,
		Handler: mux,
	}
	fmt.Printf("Serving on: http://localhost:%s\n", port)

	// create handler
	handler := http.FileServer(http.Dir("."))
	mux.Handle("/", handler)

	// handle request
	mux.HandleFunc("POST /api/newuser", apicfg.HandlerCreateUser)

	// start server
	err = cardinalServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error listening and serve")
		return
	}
	return
}