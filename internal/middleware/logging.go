package middleware

import (
	"net/http"
	"time"

	"bank-service/internal/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logger.Log.Infof(
			"%s %s %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
