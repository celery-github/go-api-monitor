package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/config"
)

func TestCheckAll_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := NewChecker(NewHTTPClient(2*time.Second), 2)

	endpoints := []config.Endpoint{
		{Name: "ok", Method: "GET", URL: srv.URL, ExpectedStatus: 200},
	}

	results := c.CheckAll(context.Background(), endpoints)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].OK {
		t.Fatalf("expected OK=true, got OK=false, err=%s", results[0].Error)
	}
	if results[0].Status != 200 {
		t.Fatalf("expected status 200, got %d", results[0].Status)
	}
}

func TestCheckAll_UnexpectedStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	c := NewChecker(NewHTTPClient(2*time.Second), 1)

	endpoints := []config.Endpoint{
		{Name: "bad", Method: "GET", URL: srv.URL, ExpectedStatus: 200},
	}

	results := c.CheckAll(context.Background(), endpoints)
	if results[0].OK {
		t.Fatalf("expected OK=false")
	}
	if results[0].Status != 500 {
		t.Fatalf("expected status 500, got %d", results[0].Status)
	}
}
