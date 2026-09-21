package main

import (
	"cloud-notes/internal/endpoints"
	"cloud-notes/internal/routes"
	"log"
	"net/http"
)

func main() {

	notesList := endpoints.NewList()

	mux := routes.New(notesList)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
