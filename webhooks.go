package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/dUPYeYE/go-http/internal/auth"
	"github.com/dUPYeYE/go-http/internal/database"
)

// POST /api/polka/webhooks
func handlePolka(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		if apiKey, err := auth.GetAPIKey(r.Header); err != nil || apiKey != apiCfg.apiKey {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		type body struct {
			Event string `json:"event"`
			Data  struct {
				User uuid.UUID `json:"user_id"`
			} `json:"data"`
		}

		decoder := json.NewDecoder(r.Body)
		var b body
		if err := decoder.Decode(&b); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request", err)
			return
		}

		if b.Event == "user.upgraded" {
			if err := apiCfg.db.UpdateChirpyRed(r.Context(), database.UpdateChirpyRedParams{
				ID:          b.Data.User,
				IsChirpyRed: true,
			}); err != nil {
				respondWithError(w, http.StatusNotFound, "User not found", err)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	})
}
