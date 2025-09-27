package main

import (
	"net/http"
	"time"

	"github.com/konradgj/boot.server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing bearer token", err)
		return
	}

	rToken, err := cfg.database.GetRefreshToken(req.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not get refresh token", err)
		return
	}
	if rToken.ExpiresAt.Before(time.Now()) || rToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token expired", nil)
		return
	}

	token, err := auth.MakeJWT(rToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Token{
		Token: token,
	})
}
