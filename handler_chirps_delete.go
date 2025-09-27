package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/konradgj/boot.server/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, req *http.Request) {
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing jwt", err)
		return
	}

	userId, err := auth.ValidateJWT(bearer, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirp, err := cfg.database.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirp", err)
		return
	}

	if userId != chirp.UserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.database.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
