package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/konradgj/boot.server/internal/auth"
	"github.com/konradgj/boot.server/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	reqBody := UserLogin{}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.database.GetUser(req.Context(), reqBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	rToken := auth.MakeRefreshToken()

	refreshToken, err := cfg.database.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     rToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token: %w", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}
