package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFileServerOnlyServesPublicDir(t *testing.T) {
	cfg := &apiConfig{}
	mux := newServeMux(cfg)

	tests := []struct {
		name     string
		path     string
		wantCode int
	}{
		{name: "index page", path: "/app/", wantCode: http.StatusOK},
		{name: "logo asset", path: "/app/assets/logo.png", wantCode: http.StatusOK},
		{name: "env file", path: "/app/.env", wantCode: http.StatusNotFound},
		{name: "source file", path: "/app/main.go", wantCode: http.StatusNotFound},
		{name: "go module file", path: "/app/go.mod", wantCode: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if rec.Code != tt.wantCode {
				t.Errorf("GET %s: got status %d, want %d", tt.path, rec.Code, tt.wantCode)
			}
		})
	}
}

func TestFileServerStillCountsHits(t *testing.T) {
	cfg := &apiConfig{}
	mux := newServeMux(cfg)

	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/app/", nil))

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/metrics", nil))
	if !strings.Contains(rec.Body.String(), "visited 1 times") {
		t.Errorf("metrics body = %q, want it to report 1 visit", rec.Body.String())
	}
}

func TestValidateChirp(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{name: "short chirp", body: "Hello, Chirpy!", want: "Hello, Chirpy!"},
		{name: "140 multibyte characters are allowed", body: strings.Repeat("ñ", 140), want: strings.Repeat("ñ", 140)},
		{name: "141 characters are rejected", body: strings.Repeat("a", 141), wantErr: true},
		{name: "bad words are masked", body: "what a Kerfuffle today", want: "what a **** today"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateChirp(tt.body)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateChirp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("validateChirp() = %q, want %q", got, tt.want)
			}
		})
	}
}
