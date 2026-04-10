package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/swissymissy/Cardinal/internal/handler"
)

func main() {

	godotenv.Load()
	port := os.Getenv("PORT")            // load port
	platform := os.Getenv("PLATFORM")    // check if is dev
	dbURL := os.Getenv("DB_URL")         // load db url
	jwtSecret := os.Getenv("JWT_SECRET") // load jwt secret
	baseUrl := os.Getenv("BASE_URL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Printf("Invalid SMTP_PORT: %s\n", err)
		return
	}

	// open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err)
		return
	}
	dbQuery := database.New(db)

	// connect to rabbitmq
	rabbitConnectionStr := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitConnectionStr)
	if err != nil {
		fmt.Printf("Failed to establish connection to Rabbit server: %s\n", err)
		return
	}
	// create a channel from the connection
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Can't create new channel: %s\n", err)
		return
	}
	defer conn.Close()
	defer ch.Close()

	// create apiConfig
	apicfg := &handler.ApiConfig{
		DB:           dbQuery,
		Port:         port,
		Platform:     platform,
		JWTSecret:    jwtSecret,
		MQConn:       conn,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		BaseURL:      baseUrl,
	}

	// server mux
	mux := http.NewServeMux()
	// create http server
	address := fmt.Sprintf(":%s", port)
	cardinalServer := http.Server{
		Addr:    address,
		Handler: mux,
	}

	// create handler
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer)

	// handle request
	mux.HandleFunc("GET /api/health", handler.HandlerHealthCheck)
	mux.HandleFunc("POST /api/newuser", apicfg.HandlerCreateUser)
	mux.HandleFunc("POST /admin/reset", apicfg.HandlerResetUsers)
	mux.HandleFunc("POST /api/userlogin", apicfg.HandlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apicfg.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apicfg.HandlerRevokeRefreshToken)
	mux.HandleFunc("POST /api/newchirp", apicfg.HandlerCreateChirp)
	mux.HandleFunc("GET /api/getallchirps", apicfg.HandlerGetAllChirps)
	mux.HandleFunc("POST /api/feed", apicfg.HandlerGetFeed)
	mux.HandleFunc("DELETE /api/chirps/{chirpsID}", apicfg.HandlerDeleteChirp)
	mux.HandleFunc("GET /api/chirps/{chirpsID}", apicfg.HandlerGetOneChirp)
	mux.HandleFunc("POST /api/newfollow", apicfg.HandlerFollowUser)
	mux.HandleFunc("DELETE /api/unfollow", apicfg.HandlerUnfollow)
	mux.HandleFunc("GET /api/users/{identifier}", apicfg.HandlerGetUser)
	mux.HandleFunc("GET /api/users/{userID}/followers", apicfg.HandlerGetFollowers)
	mux.HandleFunc("GET /api/users/{userID}/followings", apicfg.HandlerGetFollowings)
	mux.HandleFunc("GET /api/notifications", apicfg.HandlerGetNotifications)
	mux.HandleFunc("PUT /api/notifications", apicfg.HandlerMarkAllRead)
	mux.HandleFunc("PUT /api/notifications/{notifID}", apicfg.HandlerMarkOneRead)
	mux.HandleFunc("POST /api/verify/request", apicfg.HandlerRequestVerification)
	mux.HandleFunc("GET /api/verify", apicfg.HandlerVerifyEmail)
	mux.HandleFunc("POST /api/chirps/{chirpID}/react", apicfg.HandlerAddReaction)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}/react", apicfg.HandlerRemoveReaction)
	mux.HandleFunc("GET /api/chirps/{chirpID}/react", apicfg.HandlerGetReactions)
	mux.HandleFunc("POST /api/chirps/{chirpID}/comments", apicfg.HandlerCreateComment)
	mux.HandleFunc("GET /api/chirps/{chirpID}/comments", apicfg.HandlerGetComments)
	mux.HandleFunc("DELETE /api/comments/{commentID}", apicfg.HandlerDeleteComment)
	mux.HandleFunc("PUT /api/comments/{commentID}", apicfg.HandlerEditComment)

	// run server in background
	go func() {
		fmt.Printf("Serving on: http://localhost:%s\n", port)
		if err := cardinalServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	// this blocks until OS sends SIGTERM or SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	// give in-flight requests up to 10s to finish
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := cardinalServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error. Forced shutdown: %s\n", err)
	}
	log.Println("Graceful shutdown complete")
}
