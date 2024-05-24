package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.WithFields(logrus.Fields{
				"method": r.Method,
				"url":    r.URL.Path,
				"time":   time.Since(start),
			}).Info("handled request")
		})
	}
}
