package main

import (
	"net/http"

	"github.com/google/uuid"
)

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func shortDirFromRequest(r *http.Request) string {
	sortDir := r.URL.Query().Get("sort")
	if sortDir != "" {
		return sortDir
	}
	return ""
}

func sortChirps(chirps *[]Chirp, sortDir string) {
	if sortDir == "desc" {
		for i := 0; i < len(*chirps)/2; i++ {
			(*chirps)[i], (*chirps)[len(*chirps)-1-i] = (*chirps)[len(*chirps)-1-i], (*chirps)[i]
		}
	}
}

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	var chirps []Chirp

	if authorID != uuid.Nil {
		dbChirps, err := cfg.queries.ListChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
		// Convertir manualmente
		for _, c := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				UserID:    c.UserID,
				Body:      c.Body,
			})
		}
	} else {
		dbChirps, err := cfg.queries.ListChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
		// Convertir manualmente
		for _, c := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				UserID:    c.UserID,
				Body:      c.Body,
			})
		}
	}

	sortDir := shortDirFromRequest(r)
	sortChirps(&chirps, sortDir)
	respondWithJSON(w, http.StatusOK, chirps)
}
