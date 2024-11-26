package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/dUPYeYE/go-http/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	env            string
	secret         string
	apiKey         string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if os.Getenv("ENV") == "" {
		log.Fatal("ENV is required")
	}
	if os.Getenv("DB_URL") == "" {
		log.Fatal("DB_URL is required")
	}

	// database
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		env:            os.Getenv("ENV"),
		secret:         os.Getenv("AUTH_SECRET"),
		apiKey:         os.Getenv("API_KEY"),
	}

	// app
	serveMux := http.NewServeMux()
	serveMux.Handle(
		"GET /app",
		http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("./app")))),
	)
	serveMux.Handle(
		"GET /app/assets/",
		http.StripPrefix(
			"/app/assets",
			apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("./assets"))),
		),
	)

	// api
	serveMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})
	serveMux.Handle("GET /api/chirps/{id}", handleGetChirp(apiCfg))
	serveMux.Handle("GET /api/chirps", handleGetChirps(apiCfg))
	serveMux.Handle("POST /api/chirps", handleNewChirp(apiCfg))
	serveMux.Handle("DELETE /api/chirps/{id}", handleDeleteChirp(apiCfg))

	serveMux.Handle("POST /api/login", handleLogin(apiCfg))
	serveMux.Handle("POST /api/refresh", handleRefresh(apiCfg))
	serveMux.Handle("POST /api/revoke", handleRevoke(apiCfg))

	serveMux.Handle("GET /api/users", handleGetUsers(apiCfg))
	serveMux.Handle("POST /api/users", handleNewUser(apiCfg))
	serveMux.Handle("PUT /api/users", handleUpdateUser(apiCfg))

	// webhooks
	serveMux.Handle("POST /api/polka/webhooks", handlePolka(apiCfg))

	// admin
	serveMux.Handle("GET /admin/metrics", handleMetrics(apiCfg))
	serveMux.Handle("POST /admin/reset", handeReset(apiCfg))

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
