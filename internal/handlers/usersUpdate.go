package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(jwtToken, cfg.JWTSecret)
	if err != nil {
		helpers.RespondWithError(w, 401, fmt.Sprintf("Unauthorized: %s", err))
		return
	}

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var params parameters
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Email or password is not valid")
		return
	}

	newHashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		helpers.RespondWithError(w, 401, "Failed to update user data")
		return
	}

	userEmail, err := cfg.DBQueries.UpdateUser(context.Background(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: newHashedPass,
	})
	if err != nil {
		helpers.RespondWithError(w, 401, "Failed to update user data")
		return
	}

	userResp := struct {
		Email string `json:"email"`
	}{
		Email: userEmail,
	}

	helpers.RespondWithJSON(w, 200, userResp)
}
