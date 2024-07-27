package main

import "net/http"

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}
