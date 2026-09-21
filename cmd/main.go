package main

import (
	"cloud-notes/internal/routes"
	"log"
	"net/http"
)

func main() {
	mux := routes.New()
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}