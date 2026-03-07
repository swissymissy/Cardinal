package main

import (
	"fmt"
	"os"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load() 				// get values from .env 
	port := os.Getenv("PORT")

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

	// start server
	err := cardinalServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error listening and serve")
		return
	}
	return

}