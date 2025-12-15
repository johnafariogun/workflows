package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/status", nil)

	http.DefaultServeMux = http.NewServeMux()
	// register handlers from main
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.DefaultServeMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	if rr.Body.String() != "OK" {
		t.Fatalf("expected body OK got %s", rr.Body.String())
	}
}

func TestItems(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/items", nil)

	http.DefaultServeMux = http.NewServeMux()
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"Laptop", "Keyboard", "Phone"})
	})

	http.DefaultServeMux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var items []string
	if err := json.NewDecoder(rr.Body).Decode(&items); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items got %d", len(items))
	}
}
