package tests

import (
	"net/http"
	"testing"
)

func TestEcho_ReturnsBody(t *testing.T) {
	engine := setupRouter()

	payload := map[string]interface{}{
		"message": "hello",
		"number":  42.0,
	}
	w := doRequest(engine, http.MethodPost, "/echo", payload, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["message"] != "hello" {
		t.Errorf("expected message 'hello', got %v", body["message"])
	}
	if body["number"] != 42.0 {
		t.Errorf("expected number 42, got %v", body["number"])
	}
}

func TestEcho_InvalidJSON(t *testing.T) {
	engine := setupRouter()

	w := doRawRequest(engine, http.MethodPost, "/echo", "{bad json}", "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] == nil {
		t.Error("expected error field in response")
	}
	if body["message"] == nil {
		t.Error("expected message field in response")
	}
}
