package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func apiHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	hits := cfg.filesServerHits.Load()
	message := fmt.Sprintf("Hits: %d", hits)

	if _, err := w.Write([]byte(message)); err != nil {
		http.Error(w, "Error trying to create response", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.filesServerHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hits: 0"))
}

type apiConfig struct {
	filesServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsIns(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.filesServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		filesServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", apiHealthHandler)
	mux.HandleFunc("/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/reset", apiCfg.handleReset)

	fileServer := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("/app/", apiCfg.middlewareMetricsIns(http.StripPrefix("/app", fileServer)))
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Printf("API Health: http://localhost:%s/healthz", port)
	log.Printf("Metrics: http://localhost:%s/metrics", port)
	log.Printf("Reset: http://localhost:%s/reset", port)
	log.Printf("App: http://localhost:%s/app/", port)

	log.Fatal(srv.ListenAndServe())
}
