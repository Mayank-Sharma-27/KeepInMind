package main

import (
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Route for sending reminder to SNS
	r.HandleFunc("/sendReminder", sendReminder).Methods("POST")

	return r
}
