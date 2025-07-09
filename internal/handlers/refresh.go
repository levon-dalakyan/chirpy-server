package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}

	tokenDB, err := cfg.DBQueries.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}
	if tokenDB.ExpiresAt.Before(time.Now()) {
		helpers.RespondWithError(w, 401, "Refresh token is expired")
		return
	}
	if tokenDB.RevokedAt.Valid && tokenDB.RevokedAt.Time.UTC().Before(time.Now().UTC()) {
		helpers.RespondWithError(w, 401, "Refresh token is expired")
		return
	}

	user, err := cfg.DBQueries.GetUserFromRefreshToken(context.Background(), refreshToken)
	if err != nil {
		helpers.RespondWithError(w, 401, "Failed to login")
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		helpers.RespondWithError(w, 401, "Failed to login")
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	}

	helpers.RespondWithJSON(w, 200, resp)
}
