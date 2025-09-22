package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirpy(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseValid struct {
		Valid bool `json:"valid"`
	}

	reqBody := requestBody{}
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(reqBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, responseValid{Valid: true})
}
