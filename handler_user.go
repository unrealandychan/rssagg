package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/unrealandychan/rssagg/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `name`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %v", err))
		return
	}

	user, userCreateError := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if userCreateError != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", userCreateError))
		return
	}
	responseWithJson(w, http.StatusOK, dbUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {

	responseWithJson(w, http.StatusOK, dbUserToUser(user))
}
