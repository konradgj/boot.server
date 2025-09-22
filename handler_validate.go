package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirpy(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	reqBody := requestBody{}
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	body := reqBody.Body

	if len(body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := makeStringSet([]string{"kerfuffle", "sharbert", "fornax"})
	body = getCleanedBody(body, profaneWords)

	respondWithJSON(w, http.StatusOK, responseBody{CleanedBody: body})
}

func getCleanedBody(body string, profaneWords map[string]struct{}) string {
	parts := strings.Split(body, " ")
	for i, word := range parts {
		_, exists := profaneWords[strings.ToLower(word)]
		if exists {
			parts[i] = "****"
		}
	}
	return strings.Join(parts, " ")
}

func makeStringSet(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set
}
