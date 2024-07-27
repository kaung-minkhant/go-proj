package main

import (
	"fmt"
	"go-proj/6-rss-agg/internal/auth"
	"go-proj/6-rss-agg/internal/database"
	"net/http"
)

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error %v", err))
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get user %v", err))
			return
		}
		handler(w, r, user)
	}
}
