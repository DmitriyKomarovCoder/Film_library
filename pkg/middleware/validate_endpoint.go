package middleware

import (
	"net/http"

	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
)

func ValidateEndpoint(next http.Handler, l *logger.Logger) http.Handler {
	allowedEndpoints := map[string]map[string]bool{
		"/actors": {
			"GET": true,
		},
		"/actors/add": {
			"POST": true,
		},
		"/actors/update": {
			"PUT": true,
		},
		"/actors/delete": {
			"DELETE": true,
		},
		"/films/add": {
			"POST": true,
		},
		"/films/update": {
			"PUT": true,
		},
		"/films": {
			"GET": true,
		},
		"films/delete": {
			"DELETE": true,
		},
		"/films/search": {
			"GET": true,
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedMethods, ok := allowedEndpoints[r.URL.Path]
		if !ok {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
			return
		}

		if !allowedMethods[r.Method] {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
