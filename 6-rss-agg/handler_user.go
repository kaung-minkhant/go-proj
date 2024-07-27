package main

import (
	"encoding/json"
	"fmt"
	"go-proj/6-rss-agg/internal/database"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HanderCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing json %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create user %v", err))
		return
	}
	RespondWithJson(w, http.StatusCreated, databaseUserToUser(user))
}

func HandlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJson(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *ApiConfig) HandlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limitString := r.URL.Query().Get("limit")
	var limit32 int32 = 10
	if limitString != "" {
		limit, err := strconv.ParseInt(limitString, 0, 32)
		if err != nil {
			limit32 = 10
		} else {
			limit32 = int32(limit)
		}
	}
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit32,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Cannot get latest post %v", err))
		return
	}
	RespondWithJson(w, http.StatusOK, databasePostsToPosts(posts))
}
