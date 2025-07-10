package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/internal/handlers"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		os.Exit(1)
	}
	dbQueries := database.New(db)

	filepathRoot := "."
	port := "8080"

	apiCfg := handlers.ApiConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      dbQueries,
		Platform:       platform,
		JWTSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))

	mux.HandleFunc("GET /api/healthz", handlers.HandlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.HandlerUsers)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUsersUpdate)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandlerRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlerPolkaWebhook)

	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerChirpsGetOne)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerChirpsCreate)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandlerChirpsDeleteOne)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
