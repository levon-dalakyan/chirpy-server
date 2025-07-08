package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/levon-dalakyan/chirpy-server/methods"
)

func main() {
	filepathRoot := "."
	port := "8080"

	apiCfg := methods.ApiConfig{
		FileserverHits: atomic.Int32{},
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
