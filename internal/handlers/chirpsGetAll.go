package handlers

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/chirpy-server/internal/database"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerChirpsGetAll(w http.ResponseWriter, req *http.Request) {
	idQuery := req.URL.Query().Get("author_id")
	sortQuery := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if idQuery != "" {
		authorId, err := uuid.Parse(idQuery)
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

	if sortQuery == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	var respSlice []ChirpResp
	for _, c := range chirps {
		respSlice = append(respSlice, formatChirpRespToJSON(c))
	}

	helpers.RespondWithJSON(w, 200, respSlice)
}
