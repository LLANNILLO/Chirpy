package main

import (
	"net/http"
	"time"

	"github.com/llannillo/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
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

	err = cfg.queries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
