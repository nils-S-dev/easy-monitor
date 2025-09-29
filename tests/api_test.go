package tests

import (
	"easy-monitor/internal/api"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestHealthCheck(t *testing.T) {
	// Start test server with your handler
	ts := initApiTestServer()
	defer ts.Close()

	// Send real HTTP request
	resp, err := http.Get(fmt.Sprintf("%v/health", ts.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: got %v", resp.StatusCode)
	}

	var data HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if data.Status != "success" {
		t.Errorf("unexpected status: got %q, want %q", data.Status, "success")
	}

	if data.Message != "OK" {
		t.Errorf("unexpected message: got %q, want %q", data.Message, "OK")
	}
}

func initApiTestServer() *httptest.Server {
	return httptest.NewServer(api.Routes())
}
