package main

import (
	"fmt"
	"net/http"
)

// POST /admin/reset
func handeReset(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiCfg.env != "development" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err := apiCfg.db.Reset(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		apiCfg.fileserverHits.Store(0)
	})
}

// GET /admin/metrics
func handleMetrics(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
}
