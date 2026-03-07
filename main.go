package main

import (
	"fmt"
	"os"
	"net/http"
)

func main() {

	port := os.Getenv("PORT")

	// server mux
	serverMux := http.NewServerMux()

	// server
	cardinalServer := http.Server{
		Addr: ":" + port,
		Handler: serverMux,
	}
	fmt.Printf("Serving on: http://localhost:%s\n", port)

	// start server
	err := cardinalServer.ListenAndServe()
	if err != nil {
		fmt.Println("Error listening and serve")
		return
	}
	return

}