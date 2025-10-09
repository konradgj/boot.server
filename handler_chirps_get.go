package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/konradgj/boot.server/internal/database"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, req *http.Request) {
	authorId := req.URL.Query().Get("author_id")
	sorting := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error
	if authorId != "" {
		id, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse author_id", err)
			return
		}

		chirps, err = cfg.database.ListChirpsByAuthor(req.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	} else {
		chirps, err = cfg.database.ListChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
	}

	if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
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

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse id", err)
		return
	}

	chirp, err := cfg.database.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
