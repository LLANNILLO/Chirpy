package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirps, err := cfg.queries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get the chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirps.ID,
		CreatedAt: dbChirps.CreatedAt,
		UpdatedAt: dbChirps.UpdatedAt,
		Body:      dbChirps.Body,
		UserID:    dbChirps.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
