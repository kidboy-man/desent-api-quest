package tests

import (
	"net/http"
	"testing"
)

func TestAuth_GenerateToken(t *testing.T) {
	engine := setupRouter()

	w := doRequest(engine, http.MethodPost, "/auth/token", map[string]string{
		"username": "admin",
		"password": "password",
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	token, ok := body["token"].(string)
	if !ok || token == "" {
		t.Error("expected non-empty token in response")
	}
}

func TestAuth_MissingCredentials(t *testing.T) {
	engine := setupRouter()

	w := doRequest(engine, http.MethodPost, "/auth/token", map[string]string{
		"username": "",
		"password": "",
	}, "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] != "BAD_REQUEST" {
		t.Errorf("expected error 'BAD_REQUEST', got %v", body["error"])
	}
}

func TestAuth_ProtectedWithoutToken(t *testing.T) {
	engine := setupRouter()

	w := doRequest(engine, http.MethodGet, "/books", nil, "")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] != "UNAUTHORIZED" {
		t.Errorf("expected error 'UNAUTHORIZED', got %v", body["error"])
	}
}

func TestAuth_ProtectedWithInvalidToken(t *testing.T) {
	engine := setupRouter()

	w := doRequest(engine, http.MethodGet, "/books", nil, "invalid-token")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] != "UNAUTHORIZED" {
		t.Errorf("expected error 'UNAUTHORIZED', got %v", body["error"])
	}
}
