package middleware

import (
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"net/http"
)

func Logging(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
		logger.Infof("%s %s", req.Method, req.RequestURI)
	})
}
