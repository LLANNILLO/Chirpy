package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/llannillo/Chirpy/internal/database"
)

// Only files inside this directory are served under /app/, so secrets
// like .env and the source code are never exposed.
const filepathRoot = "./public"

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	dbQueries := database.New(db)
	// Si hay problemas porque el puerto esta bloqueado: sudo lsof -i :8080
	// Y despues realizar el: sudo kill -9 <PID>
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	mux := newServeMux(&apiCfg)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func newServeMux(cfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsList)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGet)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerChirpDelete)

	mux.HandleFunc("POST /api/login", cfg.handlerUsersLogin)
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersUpdate)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeToken)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerWebhook)

	return mux
}
