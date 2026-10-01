package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const readinessTimeout = 3 * time.Second

func registerHealthRoutes(router chi.Router, readiness []Pinger) {
	router.Get("/livez", func(w http.ResponseWriter, _ *http.Request) {
		writeHealth(w, http.StatusOK, "ok", "")
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
		defer cancel()

		for _, pinger := range readiness {
			if err := pinger.PingContext(ctx); err != nil {
				writeHealth(w, http.StatusServiceUnavailable, "unavailable", err.Error())
				return
			}
		}
		writeHealth(w, http.StatusOK, "ok", "")
	})
}

func writeHealth(w http.ResponseWriter, status int, statusText, errText string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}{Status: statusText, Error: errText})
}
