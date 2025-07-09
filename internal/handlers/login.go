package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, req *http.Request) {
	type loginDataParams struct {
		loginData
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	var params loginDataParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid payload data")
	}

	user, err := cfg.DBQueries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		responseValidationError(w)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		responseValidationError(w)
	}

	jwtExpiresIn := time.Hour
	if params.ExpiresInSeconds > 0 {
		if time.Duration(params.ExpiresInSeconds)*time.Second < time.Duration(time.Hour.Seconds()) {
			jwtExpiresIn = time.Duration(params.ExpiresInSeconds)
		}
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, jwtExpiresIn)
	if err != nil {
		helpers.RespondWithError(w, 400, "Failed to login")
	}

	userResp := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token     string `json:"token"`
	}{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
		Token:     jwtToken,
	}
	helpers.RespondWithJSON(w, 200, userResp)
}

func responseValidationError(w http.ResponseWriter) {
	helpers.RespondWithError(w, 401, "Incorrect email or password")
}
