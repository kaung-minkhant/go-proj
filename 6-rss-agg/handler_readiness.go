package main

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJson(w, http.StatusOK, struct{}{})
}
