package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/konradgj/boot.server/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	reqBody := UserLogin{}

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

	expiresIn := time.Hour
	if reqBody.ExpiresIn != nil && time.Duration(*reqBody.ExpiresIn)*time.Second < time.Hour {
		expiresIn = time.Duration(*reqBody.ExpiresIn) * time.Second
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
