package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerChirpsGetOne(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")

	chirpId, err := uuid.Parse(id)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid chirp id")
		return
	}

	chirp, err := cfg.DBQueries.GetChirp(context.Background(), chirpId)
	if err != nil {
		helpers.RespondWithError(w, 404, "Chirp with provided id is not exist")
		return
	}

	chirpResp := formatChirpRespToJSON(chirp)

	helpers.RespondWithJSON(w, 200, chirpResp)
}
