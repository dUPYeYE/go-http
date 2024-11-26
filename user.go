package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/dUPYeYE/go-http/internal/auth"
	"github.com/dUPYeYE/go-http/internal/database"
)

// GET /api/users
func handleGetUsers(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		users, err := apiCfg.db.GetAllUsers(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		respondWithJSON(w, http.StatusOK, users)
	})
}

// POST /api/users
func handleNewUser(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
			return
		}

		type body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		var b body
		if err := decoder.Decode(&b); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request", err)
			return
		}

		var err error
		b.Password, err = auth.HashPassword(b.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		user, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
			ID:       uuid.New(),
			Email:    b.Email,
			Password: b.Password,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	})
}

// PUT /api/users
func handleUpdateUser(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
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

		type body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		var b body
		if err := decoder.Decode(&b); err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request", err)
			return
		}

		hashedPassword, err := auth.HashPassword(b.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
			return
		}
		updatedUser, err := apiCfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
			ID:       userID,
			Email:    b.Email,
			Password: hashedPassword,
		})

		respondWithJSON(w, http.StatusOK, updatedUser)
	})
}
