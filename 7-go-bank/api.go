package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (apiServer *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHttpHandlerFunc(apiServer.handleLogin))
	router.HandleFunc("/accounts", makeHttpHandlerFunc(apiServer.handleAccountAdmin))
	router.HandleFunc("/accounts/{id}", withJWTAuth(makeHttpHandlerFunc(apiServer.handlerAccountWithID), apiServer.store))
	router.HandleFunc("/transfer", makeHttpHandlerFunc(apiServer.handleTransfer))

	log.Println("Go bank api listening on port", apiServer.listenAddr)
	if err := http.ListenAndServe(apiServer.listenAddr, router); err != nil {
		log.Fatal("Cannot listen with error", err)
	}
}
