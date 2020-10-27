package interfaces

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

type MockEngine struct {
	BaseCurrency       string
	ExchangeRepository ExchangeRepository
}

func (m MockEngine) GetLatestExchange() ForexResponse {
	return ForexResponse{
		Base: "TEST",
		Rates: map[string]float64{
			"IDR": 10,
		},
		Date: time.Now().Format("2006-01-02"),
	}
}

func (m MockEngine) GetExchangeByDate(date string) ForexResponse {
	return ForexResponse{
		Base: "TEST",
		Date: date,
	}
}

func (m MockEngine) GetAnalyzedRates() ForexResponse {
	return ForexResponse{
		Base: "TEST",
		RatesAnalyzer: map[string]RatesAnalyze{
			"IDR": {
				Min: 10,
				Max: 10,
				Avg: 10,
			},
		},
		Date: time.Now().Format("2006-01-02"),
	}
}

func Test_CustomMiddleware(t *testing.T) {
	customHandler := &Middleware{}

	customHandler.HandleFunc(regexp.MustCompile(`/(\w+)$`), func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("match").(string) != "foo" {
			http.NotFound(w, r)
		}
		_, _ = io.WriteString(w, r.Context().Value("match").(string))
	})

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	customHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}

	// Check the response body is what we expect.
	if rr.Body.String() != "foo" {
		t.Error()
	}

	customHandler = &Middleware{}

	customHandler.HandleFunc(regexp.MustCompile(`/hello$`), func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "")
	})

	rr = httptest.NewRecorder()
	customHandler.ServeHTTP(rr, req) // Check the status code is what we expect.

	if status := rr.Code; status != http.StatusNotFound {
		t.Error(rr.Code)
	}
}

func Test_HelloWorld(t *testing.T) {
	customHandler := &Middleware{}

	webservices := WebserviceHandler{}
	webservices.ExchangeAgent = MockEngine{}

	customHandler.HandleFunc(regexp.MustCompile(`/hello-world`), webservices.HelloWorld)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hello-world", nil)
	if err != nil {
		t.Fatal(err)
	}
	customHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}

	// Check the response body is what we expect.
	var helloWorld HelloWorld
	if err := json.NewDecoder(rr.Body).Decode(&helloWorld); err != nil {
		t.Error(err)
	}

	if helloWorld.Hello != "world" {
		t.Error(err)
	}
}

func Test_GetLatestExchange(t *testing.T) {
	customHandler := &Middleware{}

	webservices := WebserviceHandler{}
	webservices.ExchangeAgent = MockEngine{}

	customHandler.HandleFunc(regexp.MustCompile(`/rates/latest$`), webservices.GetLatestExchange)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/rates/latest", nil)
	if err != nil {
		t.Fatal(err)
	}
	customHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}

	// Check the response body is what we expect.

	expected := ForexResponse{
		Base: "TEST",
		Rates: map[string]float64{
			"IDR": 10,
		},
		Date: time.Now().Format("2006-01-02"),
	}
	var response ForexResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Error(err)
	}

	if response.Base != expected.Base {
		t.Error()
	}

	if response.Date != expected.Date {
		t.Error()
	}

	if response.Rates["IDR"] != expected.Rates["IDR"] {
		t.Error()
	}
}

func Test_GetLatestExchangeByDate(t *testing.T) {
	customHandler := &Middleware{}

	webservices := WebserviceHandler{}
	webservices.ExchangeAgent = MockEngine{}

	customHandler.HandleFunc(regexp.MustCompile(`/rates/(\d{4}-\d{2}-\d{2})$`), webservices.GetLatestExchangeByDate)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/rates/2006-01-02", nil)
	if err != nil {
		t.Fatal(err)
	}
	customHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}

	// Check the response body is what we expect.
	expected := ForexResponse{
		Base: "TEST",
		Date: "2006-01-02",
	}
	var response ForexResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Error(err)
	}

	if response.Base != expected.Base {
		t.Error()
	}

	if response.Date != expected.Date {
		t.Error()
	}
}

func Test_RatesAnalyze(t *testing.T) {
	customHandler := &Middleware{}

	webservices := WebserviceHandler{}
	webservices.ExchangeAgent = MockEngine{}

	customHandler.HandleFunc(regexp.MustCompile(`/rates/analyze$`), webservices.RatesAnalyze)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/rates/analyze", nil)
	if err != nil {
		t.Fatal(err)
	}
	customHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}

	// Check the response body is what we expect.
	expected := ForexResponse{
		Base: "TEST",
		RatesAnalyzer: map[string]RatesAnalyze{
			"IDR": {
				Min: 10,
				Max: 10,
				Avg: 10,
			},
		},
		Date: time.Now().Format("2006-01-02"),
	}
	var response ForexResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Error(err)
	}

	if response.Base != expected.Base {
		t.Error()
	}

	if response.Date != expected.Date {
		t.Error()
	}

	if response.Rates["IDR"] != expected.Rates["IDR"] {
		t.Error()
	}
}
