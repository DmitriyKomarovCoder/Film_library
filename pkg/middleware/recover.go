package middleware

import (
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"net/http"
	"runtime/debug"
)

func PanicRecovery(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}
