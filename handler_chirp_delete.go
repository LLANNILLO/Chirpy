package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/llannillo/Chirpy/internal/auth"
	"github.com/llannillo/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	chirp, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not the owner of this chirp", nil)
		return
	}

	err = cfg.queries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userID,
		ID:     chirpID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
