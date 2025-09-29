package tests

import (
	"easy-monitor/internal/config"
	"easy-monitor/internal/monitor"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestMessage struct {
	Message string `json:"message"`
}

func TestMonitorForStatus(t *testing.T) {
	// Start test server with your handler
	mock := initMonitorMockServer()
	defer mock.Close()

	expects200 := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects 200",
		Endpoint: url(mock.URL, "/status/200"),
		Method:   "GET",
		Expect: config.MonitorExpected{
			Status: 200,
		},
	})

	if expects200.Status != monitor.StatusSuccess {
		fmt.Println(expects200.Error)
		t.Errorf("unexpected status: got %q, expected %q", expects200.Status, monitor.StatusSuccess)
	}

	if expects200.Received.Status != 200 {
		fmt.Println(expects200.Error)
		t.Errorf("unexpected status: got %v, expected %v", expects200.Received.Status, 200)
	}

	expects403ButGets404 := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects 403 but gets 404",
		Endpoint: url(mock.URL, "/status/404"),
		Method:   "GET",
		Expect: config.MonitorExpected{
			Status: 403,
		},
	})

	if expects403ButGets404.Status != monitor.StatusFail {
		fmt.Println(expects403ButGets404.Error)
		t.Errorf("unexpected status: got %v, expected %v", expects403ButGets404.Status, monitor.StatusFail)
	}

	if expects403ButGets404.Received.Status != 404 {
		fmt.Println(expects403ButGets404.Error)
		t.Errorf("unexpected status: got %v, expected %v", expects403ButGets404.Status, 404)
	}

	expects400ButGets500 := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects 400 but gets 500",
		Endpoint: url(mock.URL, "/status/500"),
		Method:   "DELETE",
		Expect: config.MonitorExpected{
			Status: 400,
		},
	})

	if expects400ButGets500.Status != monitor.StatusFail {
		fmt.Println(expects400ButGets500.Error)
		t.Errorf("unexpected status: got %v, expected %v", expects400ButGets500.Status, monitor.StatusFail)
	}

	if expects400ButGets500.Received.Status != 500 {
		fmt.Println(expects400ButGets500.Error)
		t.Errorf("unexpected status: got %v, expected %v", expects400ButGets500.Received.Status, 500)
	}
}

func TestMonitorForBody(t *testing.T) {
	// Start test server with your handler
	mock := initMonitorMockServer()
	defer mock.Close()

	expectedResponseBody := "{ \"message\": \"Echo!\" }"
	var decodedExpectedResponseBody TestMessage
	err := json.Unmarshal([]byte(expectedResponseBody), &decodedExpectedResponseBody)

	if err != nil {
		t.Fatal(err)
	}

	expectsBodyFromRequestAndShouldSucceed := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects Request Body to be echoed and succeeds",
		Endpoint: url(mock.URL, "/body/echo"),
		Method:   "POST",
		Body:     expectedResponseBody,
		Expect: config.MonitorExpected{
			Status: http.StatusOK,
			Body:   expectedResponseBody,
		},
	})

	if expectsBodyFromRequestAndShouldSucceed.Status != monitor.StatusSuccess {
		t.Errorf("unexpected status: got %q, expected %q", expectsBodyFromRequestAndShouldSucceed.Status, monitor.StatusSuccess)
	}

	if expectsBodyFromRequestAndShouldSucceed.Received.Status != 200 {
		t.Errorf("unexpected status: got %v, expected %v", expectsBodyFromRequestAndShouldSucceed.Received.Status, 200)
	}

	var decodedResponseBody TestMessage
	err = json.Unmarshal([]byte(expectsBodyFromRequestAndShouldSucceed.Received.Body), &decodedResponseBody)

	if err != nil {
		t.Fatal(err)
	}

	if decodedResponseBody.Message != decodedExpectedResponseBody.Message {
		t.Errorf("unexpected body: got %v, expected %v", expectsBodyFromRequestAndShouldSucceed.Received.Body, expectedResponseBody)
	}

	expectsBodyFromRequestAndShouldFail := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects Request Body to be echoed and fails",
		Endpoint: url(mock.URL, "/body/fail"),
		Method:   "PUT",
		Body:     expectedResponseBody,
		Expect: config.MonitorExpected{
			Status: http.StatusOK,
			Body:   expectedResponseBody,
		},
	})

	if expectsBodyFromRequestAndShouldFail.Status != monitor.StatusFail {
		fmt.Println(expectsBodyFromRequestAndShouldFail.Error)
		t.Errorf("unexpected status: got %v, expected %v", expectsBodyFromRequestAndShouldFail.Status, monitor.StatusFail)
	}

	if expectsBodyFromRequestAndShouldFail.Received.Status != 200 {
		fmt.Println(expectsBodyFromRequestAndShouldFail.Error)
		t.Errorf("unexpected status: got %v, expected %v", expectsBodyFromRequestAndShouldFail.Received.Status, 200)
	}
}

func TestMonitorForHeaders(t *testing.T) {
	// Start test server with your handler
	mock := initMonitorMockServer()
	defer mock.Close()

	matchingResponseHeaders := map[string]string{
		"X-Test-Header-1": "Hello!",
		"X-Test-Header-2": "World!",
	}

	expectsRespondedHeadersAndShouldSucceed := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects Headers that are responded",
		Endpoint: url(mock.URL, "/headers"),
		Method:   "GET",
		Expect: config.MonitorExpected{
			Headers: matchingResponseHeaders,
		},
	})

	if expectsRespondedHeadersAndShouldSucceed.Status != monitor.StatusSuccess {
		t.Errorf("unexpected status: got %q, expected %q", expectsRespondedHeadersAndShouldSucceed.Status, monitor.StatusSuccess)
	}

	mismatchingResponseHeaders := map[string]string{
		"X-Test-Header-1": "Hello!",
		"X-Test-Header-2": "dlrow!",
		"X-Test-Header-3": "Wow!",
	}

	expectsNotRespondedHeadersAndShouldFail := monitor.GetMonitorResult(config.Monitor{
		Name:     "Expects Headers that are not responded",
		Endpoint: url(mock.URL, "/headers"),
		Method:   "GET",
		Expect: config.MonitorExpected{
			Headers: mismatchingResponseHeaders,
		},
	})

	if expectsNotRespondedHeadersAndShouldFail.Status != monitor.StatusFail {
		t.Errorf("unexpected status: got %q, expected %q", expectsNotRespondedHeadersAndShouldFail.Status, monitor.StatusFail)
	}
}

func url(base string, path string) string {
	return fmt.Sprintf("%v%v", base, path)
}

func initMonitorMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/status/200" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/status/404" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method == http.MethodDelete && r.URL.Path == "/status/500" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Method == http.MethodPost && r.URL.Path == "/body/echo" {
			defer r.Body.Close()

			// decode request JSON
			var req map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}

			// extract the "message" property
			message, _ := req["message"].(string)

			// respond with JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": message,
			})
			return
		}
		if r.Method == http.MethodPut && r.URL.Path == "/body/fail" {
			// respond with JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "I am a failure",
			})
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/headers" {
			w.Header().Set("X-Test-Header-1", "Hello!")
			w.Header().Set("X-Test-Header-2", "World!")
			return
		}
		http.NotFound(w, r)
	}))
}
