package main

import (
	"net/http"

	database "github.com/PraneGIT/rssagg/internal/database"
)

func (apiConfig *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiConfig.DB.GetPosts(r.Context(), database.GetPostsParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, posts)
}
