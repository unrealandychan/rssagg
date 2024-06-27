package main

import (
	"github.com/unrealandychan/rssagg/internal/auth"
	"github.com/unrealandychan/rssagg/internal/database"
	"net/http"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			responseWithError(w, http.StatusForbidden, "Error decoding json")
			return
		}

		user, userGetError := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if userGetError != nil {
			responseWithError(w, http.StatusInternalServerError, "Error getting user")
			return
		}

		handler(w, r, user)
	}
}
