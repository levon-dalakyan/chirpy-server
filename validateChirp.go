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

	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		badWords := map[string]struct{}{
			"kerfuffle": {},
			"sharbert":  {},
			"fornax":    {},
		}
		cleanedBody := replaceBadWords(params.Body, badWords)

		validResp := struct {
			CleanedBody string `json:"cleaned_body"`
		}{
			CleanedBody: cleanedBody,
		}
		respondWithJSON(w, 200, validResp)
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
