package monitor

import (
	"bytes"
	"easy-monitor/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

type MonitorReceived struct {
	Status  int               `json:"status,omitempty"`
	Body    string            `json:"body,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type MonitorResult struct {
	Name     string                 `json:"name"`
	Endpoint string                 `json:"endpoint"`
	Method   string                 `json:"method"`
	Body     string                 `json:"body,omitempty"`
	Status   string                 `json:"status"`
	Expected config.MonitorExpected `json:"expected"`
	Received MonitorReceived        `json:"received,omitempty"`
	Error    string                 `json:"error,omitempty"`
}

func GetMonitorResults(monitors []config.Monitor) []MonitorResult {

	// wait group syncs the concurent requests
	var wg sync.WaitGroup

	// create results channel which the goroutines will populate
	results := make(chan MonitorResult, len(monitors))

	for _, monitor := range monitors {
		wg.Add(1)
		go collectMonitorResult(&wg, monitor, results)
	}

	// await wait group and close channel
	wg.Wait()
	close(results)

	var output []MonitorResult
	for res := range results {
		output = append(output, res)
	}

	return output
}

// pass client as param to allow unit testing with httptest
func GetMonitorResult(monitor config.Monitor) MonitorResult {

	client := &http.Client{}

	req, err := initRequest(
		monitor.Method,
		monitor.Endpoint,
		monitor.Body,
	)

	if err != nil {
		// Handle the error properly (log, return, etc.)
		return MonitorResult{
			Name:     monitor.Name,
			Endpoint: monitor.Endpoint,
			Method:   monitor.Method,
			Status:   StatusFail,
			Expected: monitor.Expect,
			Error:    fmt.Errorf("request failed: %w", err).Error(),
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		// Handle the error properly (log, return, etc.)
		return MonitorResult{
			Name:     monitor.Name,
			Endpoint: monitor.Endpoint,
			Method:   monitor.Method,
			Status:   StatusFail,
			Expected: monitor.Expect,
			Error:    fmt.Errorf("request failed: %w", err).Error(),
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return MonitorResult{
			Name:     monitor.Name,
			Endpoint: monitor.Endpoint,
			Method:   monitor.Method,
			Status:   StatusFail,
			Expected: monitor.Expect,
			Error:    fmt.Errorf("request failed: %w", err).Error(),
		}
	}

	var status string = StatusSuccess
	var failed bool = false
	var expectation config.MonitorExpected = monitor.Expect

	// assert status
	if expectation.Status != 0 {
		failed = expectation.Status != resp.StatusCode
	}

	// assert headers
	var receivedHeaders map[string]string = map[string]string{}
	if !failed && expectation.Headers != nil {
		for key, expectedValue := range monitor.Expect.Headers {
			actualValue := resp.Header.Get(key)
			if expectedValue != actualValue {
				failed = true
			}
			receivedHeaders[key] = actualValue
		}
	}

	// assert body
	var receivedBody string
	if expectation.Body != "" {
		receivedBody = string(body)
	}
	if !failed && receivedBody != "" {

		parsedExpectedBody, _ := readJSON([]byte(expectation.Body))
		parsedActualBody, _ := readJSON(body)
		if parsedExpectedBody == nil || parsedActualBody == nil {
			// no json --> string compare
			failed = bytes.Equal(body, []byte(expectation.Body))
		} else {
			// json --> deep compare
			result, err := compareJSON(parsedExpectedBody, parsedActualBody)
			if err != nil {
				return MonitorResult{
					Name:     monitor.Name,
					Endpoint: monitor.Endpoint,
					Method:   monitor.Method,
					Status:   StatusFail,
					Expected: expectation,
					Error:    fmt.Errorf("request failed: %w", err).Error(),
				}
			}
			failed = !result
		}

	}

	if failed {
		status = StatusFail
	}

	return MonitorResult{
		Name:     monitor.Name,
		Endpoint: monitor.Endpoint,
		Method:   monitor.Method,
		Status:   status,
		Expected: monitor.Expect,
		Received: MonitorReceived{
			Status:  resp.StatusCode,
			Headers: receivedHeaders,
			Body:    receivedBody,
		},
	}

}

func collectMonitorResult(wg *sync.WaitGroup, monitor config.Monitor, results chan<- MonitorResult) {

	defer wg.Done()

	results <- GetMonitorResult(monitor)

}

func initRequest(method string, endpoint string, body string) (*http.Request, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	return http.NewRequest(method, endpoint, reader)
}

func readJSON(body []byte) (json.RawMessage, error) {
	var js json.RawMessage
	err := json.Unmarshal(body, &js)
	return js, err
}

// deep compare json objects
func compareJSON(a, b json.RawMessage) (bool, error) {
	var objA, objB interface{}

	if err := json.Unmarshal(a, &objA); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &objB); err != nil {
		return false, err
	}
	return reflect.DeepEqual(objA, objB), nil
}
