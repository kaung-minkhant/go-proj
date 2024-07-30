package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ErrorLog(text string, a ...any) {
	fmt.Printf(text+"\n", a...)
}

func PermissionDenied(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(&ApiError{
		Error: "permission denied",
	})
}
