package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerUsers(w http.ResponseWriter, req *http.Request) {
	params := loginData{}
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request payload")
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request payload")
		return
	}

	user, err := cfg.DBQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		errStr := fmt.Sprintf("Failed to create user: %v", err)
		helpers.RespondWithError(w, 500, errStr)
	}

	userResp := struct {
		Id        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}{
		Id:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	}
	helpers.RespondWithJSON(w, 201, userResp)
}
