package main

import (
	"encoding/json"
	"fmt"
	"go-proj/6-rss-agg/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing json %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeeds(r.Context(), database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create feed %v", err))
		return
	}
	RespondWithJson(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiCfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get feeds %v", err))
		return
	}
	RespondWithJson(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
