package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/levon-dalakyan/chirpy-server/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DBQueries      *database.Queries
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileserverHits.Load())
}
