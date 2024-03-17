package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestLogging(t *testing.T) {
	mockLogger, hook := test.NewNullLogger()
	fakeLogger := &logger.Logger{mockLogger}

	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	recoveryHandler := Logging(handler, fakeLogger)

	rr := httptest.NewRecorder()

	recoveryHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, rr.Code)
	}

	if len(hook.Entries) != 1 {
		t.Errorf("Expected one log entry, got %d", len(hook.Entries))
	}

	if !strings.Contains(hook.LastEntry().Message, "GET") {
		t.Errorf("Unexpected log message: %s", hook.LastEntry().Message)
	}
}
