package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/unrealandychan/rssagg/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load the environment variables
	loadEnvErr := godotenv.Load(".env")
	if loadEnvErr != nil {
		log.Fatalln("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	go startScraping(db, 5, time.Second*2)
	// Create a new router and set up the routes
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Create a new router for the v1 API
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/user", apiCfg.handlerCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))

	// Feed routes
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handlerCreateFeeds))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	// Feed follow routes
	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	// Start the server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Println("PORT:", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error on ListenAndServe: ", err)
	}

}
