package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/dUPYeYE/go-http/internal/auth"
	"github.com/dUPYeYE/go-http/internal/database"
)

func replaceCussWords(s string) string {
	cussWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitString := strings.Split(s, " ")
	for i, word := range splitString {
		for _, cussWord := range cussWords {
			if strings.ToLower(word) == cussWord {
				splitString[i] = "****"
			}
		}
	}
	s = strings.Join(splitString, " ")

	return s
}

// GET /api/chirps/{id}
func handleGetChirp(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")
		chirp, err := apiCfg.db.GetChirp(r.Context(), uuid.MustParse(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		respBody, err := json.Marshal(chirp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(respBody)
	})
}

// GET /api/chirps
func handleGetChirps(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		chirps, err := apiCfg.db.GetAllChirps(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		respBody, err := json.Marshal(chirps)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(respBody)
	})
}

// POST /api/chirps
func handleNewChirp(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(token, apiCfg.secret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type body struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		var b body
		if err := decoder.Decode(&b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(b.Body) > 140 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b.Body = replaceCussWords(b.Body)

		chirp, err := apiCfg.db.AddChirp(r.Context(), database.AddChirpParams{
			ID:     uuid.New(),
			Body:   b.Body,
			UserID: userID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		respBody, err := json.Marshal(chirp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(respBody)
	})
}
