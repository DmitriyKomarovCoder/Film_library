package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthAdminRole(t *testing.T) {
	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "role", Value: "admin"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	Auth(handler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthUserRoleGET(t *testing.T) {
	req, err := http.NewRequest("GET", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "role", Value: "user"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	authHandler := Auth(handler)

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthUserRolePOST(t *testing.T) {
	req, err := http.NewRequest("POST", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "role", Value: "user"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	authHandler := Auth(handler)

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, rr.Code)
	}
}

func TestAuthUserWithoutCookies(t *testing.T) {
	req, err := http.NewRequest("POST", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	authHandler := Auth(handler)

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, rr.Code)
	}
}

func TestAuthUserNotCorrectCookies(t *testing.T) {
	req, err := http.NewRequest("POST", "/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{Name: "role", Value: "not_correct_user"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exampleHandler(w, r)
	})

	authHandler := Auth(handler)

	authHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, rr.Code)
	}
}
