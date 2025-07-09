package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/methods"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		os.Exit(1)
	}
	dbQueries := database.New(db)

	filepathRoot := "."
	port := "8080"

	apiCfg := methods.ApiConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	))

	mux.HandleFunc("GET /api/healthz", methods.HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", methods.HandlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
