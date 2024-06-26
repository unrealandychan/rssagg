package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	loadEnvErr := godotenv.Load(".env")
	if loadEnvErr != nil {
		log.Fatalln("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT environment variable not set")
	}
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

	router.Mount("/v1", v1Router)

	// Start the server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Println("PORT:", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error on ListenAndServe: ", err)
	}

}
