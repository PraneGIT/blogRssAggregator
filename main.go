package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"time"

	database "github.com/PraneGIT/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port not found in environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL not found in environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(conn)
	apiConfig := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiConfig.handlerCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUser))
	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.handlerGetFeeds)
	v1Router.Post("/feedFollows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollows))
	v1Router.Get("/feedFollows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))
	v1Router.Delete("/feedFollows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollows))
	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handlerGetPosts))
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
