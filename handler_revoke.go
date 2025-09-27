package main

import (
	"net/http"

	"github.com/konradgj/boot.server/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing bearer token", err)
		return
	}

	err = cfg.database.RevokeRefreshToken(req.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke session", err)
		return
	}

	respondWithCode(w, http.StatusNoContent)
}
