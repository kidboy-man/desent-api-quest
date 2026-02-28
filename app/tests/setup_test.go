package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	v1 "github.com/kidboy-man/8-level-desent/app/controllers/http/v1"
	"github.com/kidboy-man/8-level-desent/app/repositories/inmemory"
	"github.com/kidboy-man/8-level-desent/app/services"
)

const testJWTSecret = "test-secret"

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	bookRepo := inmemory.NewBookRepository()
	bookService := services.NewBookService(bookRepo)
	authService := services.NewAuthService(testJWTSecret)

	pingCtrl := v1.NewPingController()
	echoCtrl := v1.NewEchoController()
	authCtrl := v1.NewAuthController(authService)
	bookCtrl := v1.NewBookController(bookService)

	router := v1.NewRouter(engine, pingCtrl, echoCtrl, authCtrl, bookCtrl, authService)
	router.Setup()

	return engine
}

func doRequest(engine *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func doRawRequest(engine *gin.Engine, method, path, rawBody string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func parseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response JSON: %v, body: %s", err, w.Body.String())
	}
	return result
}

func getToken(t *testing.T, engine *gin.Engine) string {
	t.Helper()
	w := doRequest(engine, http.MethodPost, "/auth/token", map[string]string{
		"username": "admin",
		"password": "password",
	}, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from /auth/token, got %d: %s", w.Code, w.Body.String())
	}

	body := parseJSON(t, w)
	token, ok := body["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected non-empty token in response")
	}
	return token
}

func createBook(t *testing.T, engine *gin.Engine, token, title, author string, year int) map[string]interface{} {
	t.Helper()
	w := doRequest(engine, http.MethodPost, "/books", map[string]interface{}{
		"title":  title,
		"author": author,
		"year":   year,
	}, token)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 from POST /books, got %d: %s", w.Code, w.Body.String())
	}

	return parseJSON(t, w)
}
