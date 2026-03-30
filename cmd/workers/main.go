package main

import (
	"fmt"
	"database/sql"
	"strconv"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/swissymissy/Cardinal/internal/pubsub"
	"github.com/swissymissy/Cardinal/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	
	// get values from .env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Printf("Invalid SMTP_PORT: %s\n", err)
		return
	}

	// open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %s\n", err)
		return
	}
	dbQuery := database.New(db)

	// create worker config
	wkrcfg := &WorkerConfig{
		DB: dbQuery,
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	}

	// connect to rabbitmq
	rabbitConnectionStr := "amqp://guest:guest@localhost:5672/"
	conn, err := pubsub.Dial(rabbitConnectionStr)
	if err != nil {
		fmt.Printf("Failed to establish connection to Rabbit server: %s\n", err)
		return
	}
	defer conn.Close()

	// subscribe to "email" queue
	err = pubsub.SubscribeJSON(
		conn,
		"notifications",
		"email",
		"",
		pubsub.Durable,
		pubsub.ExchangeFanout,
		wkrcfg.HandlerEmailNotification,
	)
	if err != nil {
		fmt.Printf("Failed to subscribe to email queue: %s\n", err)
		return
	}

	// wait for signal ctrl+c
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
	fmt.Println("Workers has stopped.")
}
