package main

import (
	"log"
	"net/http"

	"planx/internal/handlers"
	"planx/internal/middleware"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	r := handlers.Routes()

	middleware.InitJWKS()

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}