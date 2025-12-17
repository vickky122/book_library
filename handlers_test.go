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

func TestCreateBook_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "valid book",
			body:       `{"title":"Go","author":"Alan"}`,
			statusCode: http.StatusCreated,
		},
		{
			name:       "invalid json",
			body:       `{invalid}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetStore()

			req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			createBook(rr, req)

			if rr.Code != tt.statusCode {
				t.Fatalf("expected %d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}

func TestListBooks_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		statusCode int
	}{
		{
			name: "empty store",
			setup: func() {
				resetStore()
			},
			statusCode: http.StatusOK,
		},
		{
			name: "with data",
			setup: func() {
				resetStore()
				books[1] = Book{ID: 1, Title: "Go", Author: "Alan"}
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			rr := httptest.NewRecorder()

			listBooks(rr, req)

			if rr.Code != tt.statusCode {
				t.Fatalf("expected %d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}

func TestJSONFilterMiddleware_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		statusCode  int
	}{
		{
			name:        "valid content type",
			contentType: "application/json",
			statusCode:  http.StatusCreated,
		},
		{
			name:        "invalid content type",
			contentType: "text/plain",
			statusCode:  http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetStore()

			handler := jsonFilterMiddleware(http.HandlerFunc(createBook))

			req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Go","author":"Alan"}`))
			req.Header.Set("Content-Type", tt.contentType)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.statusCode {
				t.Fatalf("expected %d, got %d", tt.statusCode, rr.Code)
			}
		})
	}
}
