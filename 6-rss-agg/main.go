package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", HandlerReadiness)
	v1Router.Get("/error", HandlerErr)

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on Port %v", portString)
	// this way probably allow more configuration
	// serv := &http.Server{
	// 	Handler: router,
	// 	Addr:    ":" + portString,
	// }

	// err := serv.ListenAndServe()
	// or
	err := http.ListenAndServe(fmt.Sprintf(":%v", portString), router)
	if err != nil {
		log.Fatal("Failed to server", err)
	}
}
