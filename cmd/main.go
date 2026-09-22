package main

import (
	"cloud-notes/internal/endpoints"
	"cloud-notes/internal/routes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	notesList := endpoints.NewList()
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		log.Fatal("S3_BUCKET is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fileHandler := endpoints.NewFileHandler(s3.NewFromConfig(cfg), bucket)

	mux := routes.New(notesList, fileHandler)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
