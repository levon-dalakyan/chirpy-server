package handlers

import (
	"context"
	"net/http"

	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}

	err = cfg.DBQueries.RevokeRefreshToken(context.Background(), refreshToken)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}

	w.WriteHeader(204)
}
