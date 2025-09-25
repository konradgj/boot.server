package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/konradgj/boot.server/internal/auth"
	"github.com/konradgj/boot.server/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	reqBody := requestBody{}
	err = json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(reqBody.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.database.CreatChirp(req.Context(), database.CreatChirpParams{
		Body:   cleaned,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    userId,
	})
}

func validateChirp(body string) (string, error) {

	if len(body) > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}

	profaneWords := makeStringSet([]string{"kerfuffle", "sharbert", "fornax"})
	return getCleanedBody(body, profaneWords), nil

}

func getCleanedBody(body string, profaneWords map[string]struct{}) string {
	parts := strings.Split(body, " ")
	for i, word := range parts {
		_, exists := profaneWords[strings.ToLower(word)]
		if exists {
			parts[i] = "****"
		}
	}
	return strings.Join(parts, " ")
}

func makeStringSet(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set
}
