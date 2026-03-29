package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

type error struct {
	Error string `json:"error"`
}

type valid struct {
	Valid     bool   `json:"valid"`
	CleanBody string `json:"cleaned_body"`
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	respBody := chirp{}
	err := decoder.Decode(&respBody)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if len(respBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		respondWithJSON(w, 200, valid{
			Valid:     true,
			CleanBody: censorProfanity(respBody.Body),
		})
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, error{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
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
