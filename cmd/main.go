package main

import (
	"cloud-notes/internal/endpoints"
	"cloud-notes/internal/routes"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-sql-driver/mysql"
)

func main() {
	dbConfig := mysql.NewConfig()
	dbConfig.User = os.Getenv("DB_USER")
	dbConfig.Passwd = os.Getenv("DB_PASSWORD")
	dbConfig.Net = "tcp"
	dbConfig.Addr = os.Getenv("DB_HOST") + ":3306"
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.ClientFoundRows = true
	if dbConfig.User == "" || dbConfig.Passwd == "" || dbConfig.Addr == ":3306" || dbConfig.DBName == "" {
		log.Fatal("DB_USER DB_PASSWORD DB_HOST and DB_NAME are required")
	}

	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	notesList := endpoints.NewList(db)
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
