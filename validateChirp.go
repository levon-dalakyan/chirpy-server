package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		cleanedBody := replaceBadWords(params.Body)

		validResp := struct {
			CleanedBody string `json:"cleaned_body"`
		}{
			CleanedBody: cleanedBody,
		}
		respondWithJSON(w, 200, validResp)
	}
}

func replaceBadWords(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splited := strings.Fields(s)
	for i, w := range splited {
		for _, bw := range badWords {
			if strings.ToLower(w) == bw {
				splited[i] = "****"
			}
		}
	}
	return strings.Join(splited, " ")
}
