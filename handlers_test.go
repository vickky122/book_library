package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetStore() {
	mu.Lock()
	defer mu.Unlock()
	books = make(map[int]Book)
	nextID = 1
}

func TestCreateBook(t *testing.T) {
	resetStore()

	body := `{"title":"Go in Action","author":"William"}`
	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	createBook(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}
}

func TestListBooks(t *testing.T) {
	resetStore()

	books[1] = Book{ID: 1, Title: "Go", Author: "Alan"}

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rr := httptest.NewRecorder()

	listBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestJSONFilterMiddlewareRejectsInvalidContentType(t *testing.T) {
	handler := jsonFilterMiddleware(http.HandlerFunc(createBook))

	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "text/plain")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected status 415, got %d", rr.Code)
	}
}
