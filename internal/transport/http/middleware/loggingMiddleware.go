package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Logs out which endpoint has been hit
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method":          r.Method,
				"Endpoint":        r.URL.Path,
				"Route Variables": mux.Vars(r),
			},
		).Info("endpoint hit")
		next.ServeHTTP(w, r)
	})
}
