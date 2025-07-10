package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/chirpy-server/internal/auth"
	"github.com/levon-dalakyan/chirpy-server/internal/helpers"
)

func (cfg *ApiConfig) HandlerChirpsDeleteOne(w http.ResponseWriter, r *http.Request) {
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

	id := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(id)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DBQueries.GetChirp(context.Background(), chirpID)
	if err != nil {
		helpers.RespondWithError(w, 404, "Chirp does not exist")
		return
	}

	if chirp.UserID != userID {
		helpers.RespondWithError(w, 403, "This chirp is not allowed to be deleted by this user")
		return
	}

	err = cfg.DBQueries.DeleteChirp(context.Background(), chirp.ID)
	if err != nil {
		helpers.RespondWithError(w, 500, "Unable to delete chirp")
		return
	}

	w.WriteHeader(204)
}
