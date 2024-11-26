package main

import (
	"encoding/json"
	"net/http"
	"sort"
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
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		id := r.PathValue("id")
		chirp, err := apiCfg.db.GetChirp(r.Context(), uuid.MustParse(id))
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Not Found", err)
			return
		}

		respondWithJSON(w, http.StatusOK, chirp)
	})
}

// GET /api/chirps
func handleGetChirps(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		var chirps []database.Chirp
		var err error
		sortDir := r.URL.Query().Get("sort")
		if sortDir == "" || (sortDir != "asc" && sortDir != "desc") {
			sortDir = "asc"
		}
		authorId := r.URL.Query().Get("author_id")
		if authorId != "" {
			chirps, err = apiCfg.db.GetChirpsFromUser(r.Context(), uuid.MustParse(authorId))
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
				return
			}
		} else {
			chirps, err = apiCfg.db.GetAllChirps(r.Context())
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
				return
			}
		}

		if sortDir == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		}

		respondWithJSON(w, http.StatusOK, chirps)
	})
}

// POST /api/chirps
func handleNewChirp(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		userID, err := auth.ValidateJWT(token, apiCfg.secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
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

		respondWithJSON(w, http.StatusCreated, chirp)
	})
}

// DELETE /api/chirps/{id}
func handleDeleteChirp(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid auth header", err)
			return
		}

		userID, err := auth.ValidateJWT(token, apiCfg.secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		id := r.PathValue("id")
		chirp, err := apiCfg.db.GetChirp(r.Context(), uuid.MustParse(id))
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Not Found", err)
			return
		}

		if chirp.UserID != userID {
			respondWithError(w, http.StatusForbidden, "Forbidden", nil)
			return
		}

		if err := apiCfg.db.RemoveChirp(r.Context(), uuid.MustParse(id)); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
