package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kalshiev/chirpy/internal/auth"
)

func (cfg *apiConfig) HandlerUpgradeCR(w http.ResponseWriter, r *http.Request) {
	ApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ApiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	upgrade := parameters{}
	err = decoder.Decode(&upgrade)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := uuid.Parse(upgrade.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if upgrade.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Event other than user.upgraded")
		return
	}

	err = cfg.db.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
