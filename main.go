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
	serverMux := http.NewServeMux()

	// server
	address := fmt.Sprintf(":%s", port)
	cardinalServer := http.Server{
		Addr: address,
		Handler: serverMux,
	}
	fmt.Printf("Serving on: http://localhost:%s\n", port)

	// create handler
	handler := http.FileServer(http.Dir("."))
	serverMux.Handle("/", handler)

	// start server
	err := cardinalServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error listening and serve")
		return
	}
	return

}