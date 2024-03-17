package middleware

import (
	"net/http"
)

// здесь типо взяли куку, пошли в бд, взяли юзера и дальше проверки
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("role")
		if err != nil || cookie == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role := cookie.Value
		if role != "admin" && role != "user" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if role == "user" && (r.Method != http.MethodGet) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
