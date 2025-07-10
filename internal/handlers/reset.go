package handlers

import (
	"context"
	"net/http"

	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.Platform != "dev" {
		helpers.RespondWithError(w, 403, "Forbidden")
		return
	}

	err := cfg.DBQueries.DeleteUsers(context.Background())
	if err != nil {
		helpers.RespondWithError(w, 500, "Failed to delete users")
		return
	}

	w.WriteHeader(200)
}
