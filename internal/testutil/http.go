// Package testutil contains offline adapters used by bpi domain tests.
package testutil

import (
	"io"
	"net/http"
	"strings"
)

// RoundTripFunc adapts a function to http.RoundTripper.
type RoundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements http.RoundTripper.
func (f RoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

// Response returns an in-memory HTTP response.
func Response(status int, contentType, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{contentType}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// JSONResponse returns an in-memory JSON HTTP response.
func JSONResponse(status int, body string) *http.Response {
	return Response(status, "application/json", body)
}
