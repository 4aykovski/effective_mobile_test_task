package middleware

import (
	"log/slog"
	"net/http"
)

func Swagger(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		log = log.With(
			slog.String("middleware", "swagger"),
		)

		log.Debug("Swagger middleware initialized")

		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "https://editor.swagger.io")

			handler.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
