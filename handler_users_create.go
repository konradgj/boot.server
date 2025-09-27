package main

import (
	"encoding/json"
	"net/http"

	"github.com/konradgj/boot.server/internal/auth"
	"github.com/konradgj/boot.server/internal/database"
)

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, req *http.Request) {
	reqBody := UserLogin{}
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPass, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	user, err := cfg.database.CreateUser(req.Context(), database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
