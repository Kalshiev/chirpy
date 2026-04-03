package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kalshiev/chirpy/internal/auth"
	"github.com/kalshiev/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	respBody := params{}
	err := decoder.Decode(&respBody)
	if err != nil {
		respondWithError(w, 400, "JSON decoding failed")
		return
	}

	if respBody.Password == "" {
		respondWithError(w, 400, "Please provide a password")
		return
	}

	hash, err := auth.HashPassword(respBody.Password)
	if err != nil {
		respondWithError(w, 400, "Hashing failed")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          respBody.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email          string `json:"email"`
		Password       string `json:"password"`
		ExpiresSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	respBody := params{}
	err := decoder.Decode(&respBody)
	if err != nil {
		respondWithError(w, 400, "JSON decoding failed")
		return
	}

	if respBody.Password == "" || respBody.Email == "" {
		respondWithError(w, 400, "Please provide an email and password")
		return
	}

	if respBody.ExpiresSeconds == 0 || respBody.ExpiresSeconds > 3600 {
		respBody.ExpiresSeconds = 3600
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), respBody.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	valid, err := auth.CheckPasswordHash(respBody.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(respBody.ExpiresSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
