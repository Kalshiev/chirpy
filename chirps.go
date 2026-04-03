package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kalshiev/chirpy/internal/auth"
	"github.com/kalshiev/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
		Token  string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedUser, err := auth.ValidateJWT(reqToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: validatedUser,
		Body:   censorProfanity(params.Body),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func censorProfanity(body string) (cleaned string) {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	lower := strings.ToLower(body)
	lwords := strings.Split(lower, " ")
	owords := strings.Split(body, " ")

	var new []string

	for idx, word := range lwords {
		if badWords[word] {
			new = append(new, "****")
		} else {
			new = append(new, owords[idx])
		}
	}

	return strings.Join(new, " ")
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
	}

	response := make([]Chirp, len(chirps))

	for idx, chirp := range chirps {
		response[idx] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println(uuid)
	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp by ID")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedUser, err := auth.ValidateJWT(reqToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if chirp.UserID != validatedUser {
		respondWithError(w, http.StatusForbidden, http.StatusText(403))
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		ID:     uuid,
		UserID: validatedUser,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
