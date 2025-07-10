package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerChirpsGetAll(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if id != "" {
		authorId, err := uuid.Parse(id)
		if err != nil {
			helpers.RespondWithError(w, 400, "author_id is not valid")
		}
		chirps, err = cfg.DBQueries.GetChirpsByUserId(context.Background(), authorId)
	} else {
		chirps, err = cfg.DBQueries.GetChirps(context.Background())
	}

	if err != nil {
		helpers.RespondWithError(w, 500, "Failed to get all chirps")
	}

	var respSlice []ChirpResp
	for _, c := range chirps {
		respSlice = append(respSlice, formatChirpRespToJSON(c))
	}

	helpers.RespondWithJSON(w, 200, respSlice)
}
