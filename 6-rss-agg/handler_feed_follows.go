package main

import (
	"encoding/json"
	"fmt"
	"go-proj/6-rss-agg/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing json %v", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create feed follow %v", err))
		return
	}
	RespondWithJson(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *ApiConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get feed follows %v", err))
		return
	}
	RespondWithJson(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *ApiConfig) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdString := chi.URLParam(r, "id")
	if feedFollowIdString == "" {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Delete ID required"))
		return
	}
	feedFollowID, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid UUID"))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID:     feedFollowID,
	})
	if err != nil {
		RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't delete feed follow %v", err))
		return
	}
	RespondWithJson(w, http.StatusOK, struct{}{})
}
