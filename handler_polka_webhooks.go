package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/konradgj/boot.server/internal/auth"
	"github.com/konradgj/boot.server/internal/database"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	reqBody := PolkaWebhook{}

	err = json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if reqBody.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(reqBody.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse id", err)
		return
	}

	_, err = cfg.database.UpdateUserChirpyRed(req.Context(), database.UpdateUserChirpyRedParams{
		ID:          userId,
		IsChirpyRed: true,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
