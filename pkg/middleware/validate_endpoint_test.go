package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateEndpointValid(t *testing.T) {
	req, err := http.NewRequest("GET", "/actors", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	validate := ValidateEndpoint(handler)

	validate.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestValidateEndpointNotAllowed(t *testing.T) {
	req, err := http.NewRequest("GET", "/not_correct", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	validatedHandler := ValidateEndpoint(handler)

	validatedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
	}
}

func TestValidateEndpointMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("PUT", "/actors", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	validatedHandler := ValidateEndpoint(handler)

	validatedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}
