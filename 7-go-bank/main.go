package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	// env sanity checks
	token := os.Getenv("JWT_SECRET")
	if token == "" {
		log.Fatal("JWT_SECRET env variable required")
	}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal("Could not connect to db with error", err)
	}
	if err := store.Init(); err != nil {
		log.Fatal("Could not initialize db", err)
	}
	// fmt.Printf("%#v\n", store)
	server := NewApiServer(":8080", store)
	server.Run()
}
