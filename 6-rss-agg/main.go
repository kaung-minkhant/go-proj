package main

import (
	"database/sql"
	"fmt"
	"go-proj/6-rss-agg/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is required")
	}
	dbString := os.Getenv("DB_URL")
	if dbString == "" {
		log.Fatal("DB connection string is required")
	}

	conn, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	queries := database.New(conn)

	apiCfg := &ApiConfig{
		DB: queries,
	}

	go StartScrapping(queries, 10, time.Minute)

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
	v1Router.Post("/users", apiCfg.HanderCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(HandlerGetUserByApiKey))

	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandlerGetFeeds)

	v1Router.Post("/feedfollows", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeedFollow))
	v1Router.Get("/feedfollows", apiCfg.MiddlewareAuth(apiCfg.HandlerGetFeedFollows))
	v1Router.Delete("/feedfollows/{id}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteFeedFollow))

	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	log.Printf("Server starting on Port %v", portString)
	// this way probably allow more configuration
	// serv := &http.Server{
	// 	Handler: router,
	// 	Addr:    ":" + portString,
	// }

	// err := serv.ListenAndServe()
	// or
	err = http.ListenAndServe(fmt.Sprintf(":%v", portString), router)
	if err != nil {
		log.Fatal("Failed to server", err)
	}
}
