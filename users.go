package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	respBody := User{}
	err := decoder.Decode(&respBody)
	if err != nil {
		respondWithError(w, 400, "JSON decoding failed")
	}
	user, err := cfg.db.CreateUser(r.Context(), respBody.Email)

	respBody.ID = user.ID
	respBody.CreatedAt = user.CreatedAt
	respBody.UpdatedAt = user.UpdatedAt
	respondWithJSON(w, 201, respBody)
}
