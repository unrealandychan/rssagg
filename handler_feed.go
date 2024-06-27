package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/unrealandychan/rssagg/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `name`
		Url  string `url`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %v", err))
		return
	}

	feed, userFeedError := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.Url,
		UserID:    user.ID,
	})
	if userFeedError != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed: %v", userFeedError))
		return
	}
	responseWithJson(w, http.StatusOK, dbFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error on getting feeds: %v", err))
		return
	}
	responseWithJson(w, http.StatusOK, dbFeedsToFeeds(feeds))
}
