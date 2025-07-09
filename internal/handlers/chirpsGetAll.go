package handlers

import (
	"context"
	"net/http"

	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerChirpsGetAll(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.DBQueries.GetChirps(context.Background())
	if err != nil {
		helpers.RespondWithError(w, 500, "Failed to get all chirps")
	}

	var respSlice []ChirpResp
	for _, c := range chirps {
		respSlice = append(respSlice, formatChirpRespToJSON(c))
	}

	helpers.RespondWithJSON(w, 200, respSlice)
}
