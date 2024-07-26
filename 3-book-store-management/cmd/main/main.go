package main

import (
	"fmt"
	"go-proj/3-book-store-management/pkg/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	routes.RegisterBookstoreRoutes(r)

	http.Handle("/", r)

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe((":8080"), r); err != nil {
		log.Fatal(err)
	}
}
