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
	var params loginData
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

	userResp := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	}
	helpers.RespondWithJSON(w, 200, userResp)
}

func responseValidationError(w http.ResponseWriter) {
	helpers.RespondWithError(w, 401, "Incorrect email or password")
}
