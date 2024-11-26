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
			Email    string `json:"email"`
			Password string `json:"password"`
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

		jwt, err := auth.GenerateJWT(user.ID, apiCfg.secret, time.Hour)
		refresh, err := auth.GenerateRefreshToken()

		if _, err = apiCfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
			Token:  refresh,
			UserID: user.ID,
		}); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		respondWithJSON(w, http.StatusOK, response{
			User:         user,
			Token:        jwt,
			RefreshToken: refresh,
		})
	})
}

// POST /api/refresh
func handleRefresh(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		refreshToken, err := apiCfg.db.GetRefreshToken(r.Context(), bearer)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		if refreshToken.ExpiresAt.Time.Before(time.Now().UTC()) {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		jwt, err := auth.GenerateJWT(refreshToken.UserID, apiCfg.secret, time.Hour)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		respondWithJSON(w, http.StatusOK, response{Token: jwt})
	})
}

// POST /api/revoke
func handleRevoke(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		if err := apiCfg.db.RevokeRefreshToken(r.Context(), bearer); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
