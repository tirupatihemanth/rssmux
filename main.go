package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tirupatihemanth/rssmux/internal/database"
)

// Add API Configurations here
type apiConfig struct {
	DB *database.Queries
}

var apiCfg apiConfig

func init() {
	godotenv.Load(".env")

	// Setup DB Connection
	dbURL := os.Getenv("DB_CONN_URL")
	if dbURL == "" {
		log.Fatalln("DB_CONN URL environment variable not seet")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Cannot connect to the DB", err)
	}

	apiCfg = apiConfig{
		DB: database.New(db),
	}

	go scheduleScraping()
}

func main() {

	router := chi.NewRouter()
	configureMiddleware(router)

	// sub-router for v1 api namespace
	v1Router := chi.NewRouter()
	configureRoutes(v1Router)
	router.Mount("/v1", v1Router)

	// Start Server
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatalln("PORT environment variable not set")
	}

	log.Println("Starting Server")

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}

func configureRoutes(router *chi.Mux) {
	router.Get("/health", healthHandler)

	router.Post("/user", createUserHandler)
	router.Get("/user", middleware_auth(getUserHandler))

	router.Post("/feed", middleware_auth(createFeedHandler))
	router.Get("/feed", getAllFeedsHandler)

	router.Post("/feed_follow", middleware_auth(feedFollowHandler))
	router.Get("/feed_follow", middleware_auth(getUserFeedFollowsHandler))
	router.Delete("/feed_follow/{feedId}", middleware_auth(unfollowFeedHandler))

	router.Get("/post", middleware_auth(getPostsForUserHandler))
}

func configureMiddleware(router *chi.Mux) {

	// CORS Middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}
