package main

import (
	"net/http"
	"time"

	"github.com/llannillo/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	refreshTokenDB, err := cfg.queries.GetRefresh(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	if refreshTokenDB.RevokedAt.Valid && refreshTokenDB.RevokedAt.Time.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	userID, err := cfg.queries.GetUserFromRefreshToken(r.Context(), refreshTokenDB.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(userID.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
