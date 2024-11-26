package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dUPYeYE/go-http/internal/auth"
	"github.com/dUPYeYE/go-http/internal/database"
)

// POST /api/login
func handleLogin(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		type body struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			ExpiresIn int64  `json:"expires_at"`
		}
		type response struct {
			database.User
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		var b body
		if err := decoder.Decode(&b); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request", err)
			return
		}

		user, err := apiCfg.db.GetUserByEmail(r.Context(), b.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		if err := auth.CheckPasswordHash(b.Password, user.Password); err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		expirationTime := time.Hour
		if b.ExpiresIn > 0 && b.ExpiresIn < 3600 {
			expirationTime = time.Duration(b.ExpiresIn) * time.Second
		}
		token, err := auth.GenerateJWT(user.ID, apiCfg.secret, expirationTime)
		respondWithJSON(w, http.StatusOK, response{
			User:  user,
			Token: token,
		})
	})
}
