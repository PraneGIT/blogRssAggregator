package main

import (
	"net/http"

	"github.com/PraneGIT/rssagg/internal/auth"
	"github.com/PraneGIT/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cgf *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 401, err.Error())
			return
		}
		user, err := cgf.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		handler(w, r, user)
	}
}
