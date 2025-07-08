package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		errResp := errorResp{
			Error: "Invalid request payload",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	if len(params.Body) > 140 {
		errorResp := errorResp{
			Error: "Chirp is too long",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorResp)
	} else {
		validResp := struct {
			Valid bool `json:"valid"`
		}{
			Valid: true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(validResp)
	}
}
