package middleware

import (
	"net/http"
)

func ValidateEndpoint(next http.Handler) http.Handler {
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
		"/movies/add": {
			"POST": true,
		},
		"/movies/update": {
			"PUT": true,
		},
		"/movies": {
			"GET": true,
		},
		"/movies/delete": {
			"DELETE": true,
		},
		"/movies/search": {
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
