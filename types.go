package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type UserLogin struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn *int   `json:"expires_in_seconds"`
}
