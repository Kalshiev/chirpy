package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kalshiev/chirpy/internal/auth"
	"github.com/kalshiev/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	if respBody.Password == "" || respBody.Email == "" {
		respondWithError(w, 400, "Please provide an email and password")
		return
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

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: user.ID,
		ExpiresAt: sql.NullTime{
			Time:  time.Now().Add((time.Duration(24) * time.Hour) * 60),
			Valid: true,
		},
	})

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func (cfg *apiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if params.Email == "" || params.Password == "" || params.Email == "" && params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide email and password")
		return
	}

	validUser, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)

	updatedUser, err := cfg.db.UpdateUserPasswordAndEmail(r.Context(), database.UpdateUserPasswordAndEmailParams{
		ID:             validUser,
		Email:          params.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}
