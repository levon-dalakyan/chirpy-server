package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}

	if apiKey != cfg.PolkaKey {
		helpers.RespondWithError(w, 401, "wrong api key")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}
	var params parameters
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid payload")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userId, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid user_id")
		return
	}

	user, err := cfg.DBQueries.GetUser(context.Background(), userId)
	if err != nil {
		helpers.RespondWithError(w, 404, "User does not exist")
		return
	}

	err = cfg.DBQueries.UpgradeUserToChirpyRed(context.Background(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, 500, "Failed to update user to chirpy red")
		return
	}

	w.WriteHeader(204)
}
