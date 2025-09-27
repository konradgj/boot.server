package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type PolkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}
