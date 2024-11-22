package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	serveMux := http.NewServeMux()
	apiCfg := &apiConfig{}
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
	serveMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})
	serveMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`
          <html>
            <body>
              <h1>Welcome, Chirpy Admin</h1>
              <p>Chirpy has been visited %d times!</p>
            </body>
          </html>
        `, apiCfg.fileserverHits.Load())))
	})
	serveMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		apiCfg.fileserverHits.Store(0)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	log.Fatal(server.ListenAndServe())
}
