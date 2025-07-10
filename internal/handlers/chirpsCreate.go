package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

type ChirpResp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
}

func (cfg *ApiConfig) HandlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}
	userId, err := auth.ValidateJWT(jwtToken, cfg.JWTSecret)
	if err != nil {
		helpers.RespondWithError(w, 401, fmt.Sprintf("Unauthorized: %s", err))
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid request payload")
		return
	}

	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		helpers.RespondWithError(w, 400, "Chirp is too long")
	} else {
		badWords := map[string]struct{}{
			"kerfuffle": {},
			"sharbert":  {},
			"fornax":    {},
		}
		cleanedBody := replaceBadWords(params.Body, badWords)

		chirp, err := cfg.DBQueries.CreateChirp(context.Background(), database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: userId,
		})
		if err != nil {
			helpers.RespondWithError(w, 500, fmt.Sprintf("Failed to create chirp: %v", err))
		}

		chirpResp := formatChirpRespToJSON(chirp)

		helpers.RespondWithJSON(w, 201, chirpResp)
	}
}

func replaceBadWords(s string, badWords map[string]struct{}) string {
	splited := strings.Fields(s)
	for i, w := range splited {
		loweredWord := strings.ToLower(w)
		if _, ok := badWords[loweredWord]; ok {
			splited[i] = "****"
		}
	}
	return strings.Join(splited, " ")
}

func formatChirpRespToJSON(chirp database.Chirp) ChirpResp {
	chirpResp := ChirpResp{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: chirp.UpdatedAt.Format(time.RFC3339),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}

	return chirpResp
}
