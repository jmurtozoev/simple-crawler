package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValid(t *testing.T) {
	cases := []struct{
		testCase string
		value string
		expected bool
	}{
		{
			testCase: "success",
			value: "https://golang.org",
			expected: true,
		},
		{
			testCase: "fail",
			value: "golang.org",
			expected: false,
		},
	}

	for _, tc :=  range cases {
		assert.Equal(t, tc.expected, valid(tc.value))
	}
}

func TestGetTitles(t *testing.T) {
	req, err := http.NewRequest("GET", "/crawl", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("urls", "https://google.com")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()

	// initialize crawler
	h := handler{newCrawler()}
	handler := http.HandlerFunc(h.GetTitles)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"invalid_urls":[],"urls":{"https://google.com":"Google"}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}