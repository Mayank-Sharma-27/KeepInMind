package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	log.Println("Server started")

	router := SetupRoutes()
	handler := cors.Default().Handler(router)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", handler))
}
