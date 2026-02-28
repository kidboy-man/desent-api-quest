package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCreateBook_Success(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	book := createBook(t, engine, token, "The Go Programming Language", "Alan Donovan", 2015)

	if book["id"] == nil || book["id"] == "" {
		t.Error("expected non-empty id")
	}
	if book["title"] != "The Go Programming Language" {
		t.Errorf("expected title 'The Go Programming Language', got %v", book["title"])
	}
	if book["author"] != "Alan Donovan" {
		t.Errorf("expected author 'Alan Donovan', got %v", book["author"])
	}
	if book["year"] != 2015.0 {
		t.Errorf("expected year 2015, got %v", book["year"])
	}
}

func TestGetAllBooks(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	createBook(t, engine, token, "Book A", "Author A", 2020)
	createBook(t, engine, token, "Book B", "Author B", 2021)

	w := doRequest(engine, http.MethodGet, "/books", nil, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	data, ok := body["data"].([]interface{})
	if !ok {
		t.Fatal("expected data to be an array")
	}
	if len(data) != 2 {
		t.Errorf("expected 2 books, got %d", len(data))
	}
	if body["total"] != 2.0 {
		t.Errorf("expected total 2, got %v", body["total"])
	}
}

func TestGetBookByID(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	created := createBook(t, engine, token, "Test Book", "Test Author", 2023)
	bookID := created["id"].(string)

	w := doRequest(engine, http.MethodGet, "/books/"+bookID, nil, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["id"] != bookID {
		t.Errorf("expected id %s, got %v", bookID, body["id"])
	}
	if body["title"] != "Test Book" {
		t.Errorf("expected title 'Test Book', got %v", body["title"])
	}
}

func TestUpdateBook_Success(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	created := createBook(t, engine, token, "Original Title", "Original Author", 2020)
	bookID := created["id"].(string)

	w := doRequest(engine, http.MethodPut, "/books/"+bookID, map[string]interface{}{
		"title":  "Updated Title",
		"author": "Updated Author",
		"year":   2025,
	}, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	body := parseJSON(t, w)
	if body["id"] != bookID {
		t.Errorf("expected id %s, got %v", bookID, body["id"])
	}
	if body["title"] != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %v", body["title"])
	}
	if body["author"] != "Updated Author" {
		t.Errorf("expected author 'Updated Author', got %v", body["author"])
	}
	if body["year"] != 2025.0 {
		t.Errorf("expected year 2025, got %v", body["year"])
	}
}

func TestDeleteBook_Success(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	created := createBook(t, engine, token, "To Delete", "Author", 2020)
	bookID := created["id"].(string)

	w := doRequest(engine, http.MethodDelete, "/books/"+bookID, nil, token)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}

	w = doRequest(engine, http.MethodGet, "/books/"+bookID, nil, token)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 after deletion, got %d", w.Code)
	}
}

func TestSearchBooksByAuthor(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	createBook(t, engine, token, "Go Book", "Alan Donovan", 2015)
	createBook(t, engine, token, "Clean Code", "Robert Martin", 2008)
	createBook(t, engine, token, "Go in Action", "Alan Donovan", 2016)

	w := doRequest(engine, http.MethodGet, "/books?author=Alan+Donovan", nil, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	data := body["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 books by Alan Donovan, got %d", len(data))
	}
	if body["total"] != 2.0 {
		t.Errorf("expected total 2, got %v", body["total"])
	}
}

func TestPaginateBooks(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	createBook(t, engine, token, "Book 1", "Author", 2020)
	createBook(t, engine, token, "Book 2", "Author", 2021)
	createBook(t, engine, token, "Book 3", "Author", 2022)

	w := doRequest(engine, http.MethodGet, "/books?page=1&limit=2", nil, token)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := parseJSON(t, w)
	data := body["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("expected 2 books on page 1, got %d", len(data))
	}
	if body["total"] != 3.0 {
		t.Errorf("expected total 3, got %v", body["total"])
	}
	if body["page"] != 1.0 {
		t.Errorf("expected page 1, got %v", body["page"])
	}
	if body["limit"] != 2.0 {
		t.Errorf("expected limit 2, got %v", body["limit"])
	}

	w = doRequest(engine, http.MethodGet, "/books?page=2&limit=2", nil, token)
	body = parseJSON(t, w)
	data = body["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("expected 1 book on page 2, got %d", len(data))
	}
}

func TestCreateBook_InvalidFields(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	w := doRequest(engine, http.MethodPost, "/books", map[string]interface{}{
		"title": "",
	}, token)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] != "BAD_REQUEST" {
		t.Errorf("expected error 'BAD_REQUEST', got %v", body["error"])
	}
	if body["message"] == nil || body["message"] == "" {
		t.Error("expected non-empty message")
	}
}

func TestGetBook_NotFound(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	w := doRequest(engine, http.MethodGet, "/books/nonexistent-id", nil, token)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] != "NOT_FOUND" {
		t.Errorf("expected error 'NOT_FOUND', got %v", body["error"])
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	w := doRequest(engine, http.MethodPut, "/books/nonexistent-id", map[string]interface{}{
		"title":  "Title",
		"author": "Author",
		"year":   2020,
	}, token)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}

	body := parseJSON(t, w)
	if body["error"] != "NOT_FOUND" {
		t.Errorf("expected error 'NOT_FOUND', got %v", body["error"])
	}
}

func TestDeleteBook_NotFound(t *testing.T) {
	engine := setupRouter()
	token := getToken(t, engine)

	w := doRequest(engine, http.MethodDelete, fmt.Sprintf("/books/%s", "nonexistent-id"), nil, token)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	body := parseJSON(t, w)
	if body["error"] != "NOT_FOUND" {
		t.Errorf("expected error 'NOT_FOUND', got %v", body["error"])
	}
}
