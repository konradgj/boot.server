package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.database.ListChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	response := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		response[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}
