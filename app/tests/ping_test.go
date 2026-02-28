package tests

import (
	"net/http"
	"testing"
)

func TestPing_ReturnsOK(t *testing.T) {
	engine := setupRouter()

	w := doRequest(engine, http.MethodGet, "/ping", nil, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["message"] != "pong" {
		t.Errorf("expected message 'pong', got %v", body["message"])
	}
}
