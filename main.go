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
	"github.com/swissymissy/Cardinal/internal/pubsub"
)


func main() {

	// get values from .env 
	godotenv.Load() 				
	port := os.Getenv("PORT")					// load port 
	platform := os.Getenv("PLATFORM")			// check if is dev
	dbURL := os.Getenv("DB_URL")				// load db url
	jwtSecret := os.Getenv("JWT_SECRET")		// load jwt secret

	// open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err)
		return
	}
	dbQuery := database.New(db)

	// connect to rabbitmq
	rabbitConnectionStr := "amqp://guest:guest@localhost:5672/"
	conn, ch, err := pubsub.Connect(rabbitConnectionStr)
	if err != nil {
		fmt.Printf("Failed to establish connection to Rabbit server: %s\n", err)
		return 
	}
	defer conn.Close()
	defer ch.Close()

	// create apiConfig
	apicfg := &handler.ApiConfig{
		DB: dbQuery,
		Port: port,
		Platform: platform,
		JWTSecret: jwtSecret,
		MQConn: conn,
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
	mux.HandleFunc("POST /admin/reset", apicfg.HandlerResetUsers)
	mux.HandleFunc("POST /api/userlogin", apicfg.HandlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apicfg.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apicfg.HandlerRevokeRefreshToken)
	mux.HandleFunc("POST /api/newchirp", apicfg.HandlerCreateChirp)
	mux.HandleFunc("GET /api/getallchirps", apicfg.HandlerGetAllChirps)
	mux.HandleFunc("DELETE /api/chirps/{chirpsID}", apicfg.HandlerDeleteChirp)
	mux.HandleFunc("GET /api/chirps/{chirpsID}", apicfg.HandlerGetOneChirp)
	mux.HandleFunc("POST /api/newfollow", apicfg.HandlerFollowUser)
	mux.HandleFunc("DELETE /api/unfollow", apicfg.HandlerUnfollow)
	mux.HandleFunc("GET /api/users/{userID}", apicfg.HandlerGetUser)
	mux.HandleFunc("GET /api/users/{userID}/followers", apicfg.HandlerGetFollowers)
	mux.HandleFunc("GET /api/users/{userID}/followings", apicfg.HandlerGetFollowings)

	// start server
	err = cardinalServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error listening and serve")
		return
	}
	return
}