package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.queries.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all the chirps", err)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, c := range dbChirps {
		chirps[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, chirps)
}
