package client

import (
	"errors"
	"fmt"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

var (
	// ErrMissingData means the remote interface returned success without a
	// required payload.
	ErrMissingData = errors.New("bpi: response is missing required data")

	// ErrAuthenticationRequired means an operation needs an authenticated
	// account or CSRF token.
	ErrAuthenticationRequired = errors.New("bpi: authentication required")
)

// ParameterError describes an invalid caller-supplied value.
type ParameterError = bpierr.ParameterError

// TransportError wraps an error returned while performing an HTTP request.
type TransportError struct {
	Operation string
	Err       error
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("bpi: %s transport request failed", e.Operation)
}

func (e *TransportError) Unwrap() error { return e.Err }

// HTTPError reports a non-successful HTTP status.
type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("bpi: HTTP request failed with status %d", e.StatusCode)
}

// APIError reports a non-zero Bilibili response code.
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	if e.Message == "" || e.Message == "0" {
		return fmt.Sprintf("bpi: remote interface returned code %d", e.Code)
	}
	return fmt.Sprintf("bpi: remote interface returned code %d: %s", e.Code, e.Message)
}

// ResponseDecodeError reports a response-model mismatch and retains a private
// copy of the original response for explicit recovery. Formatting and JSON
// serialization never include that body.
type ResponseDecodeError struct {
	err  error
	body []byte
}

// NewResponseDecodeError records a response-model mismatch while keeping the
// body out of ordinary error formatting and logs.
func NewResponseDecodeError(err error, body []byte) *ResponseDecodeError {
	return &ResponseDecodeError{err: err, body: append([]byte(nil), body...)}
}

func (e *ResponseDecodeError) Error() string {
	return fmt.Sprintf("bpi: failed to decode response: %v", e.err)
}

func (e *ResponseDecodeError) Unwrap() error { return e.err }

// Body returns a copy of the original response body. It may contain sensitive
// data and should not be logged.
func (e *ResponseDecodeError) Body() []byte {
	return append([]byte(nil), e.body...)
}

// ResponseTooLargeError reports that a response exceeded the configured
// in-memory body limit.
type ResponseTooLargeError struct {
	Limit int64
}

func (e *ResponseTooLargeError) Error() string {
	return fmt.Sprintf("bpi: response exceeds %d-byte limit", e.Limit)
}

// RequiresLogin reports whether err represents an unauthenticated response.
func RequiresLogin(err error) bool {
	if errors.Is(err, ErrAuthenticationRequired) {
		return true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == -101 || apiErr.Code == -401 || apiErr.Code == 4_100_000 || apiErr.Code == 4_511_003 || apiErr.Code == 800501007
	}
	var httpErr *HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == 401
}

// RequiresVIP reports whether err represents a VIP-only response.
func RequiresVIP(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && (apiErr.Code == -106 || apiErr.Code == -650)
}

// IsPermissionError reports whether err represents an authorization failure.
func IsPermissionError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) && (apiErr.Code == -403 || apiErr.Code == -4 || apiErr.Code == 79_511 || apiErr.Code == 100_004 || apiErr.Code == 100_007) {
		return true
	}
	var httpErr *HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == 403
}

// IsRiskControl reports whether err represents Bilibili risk control.
func IsRiskControl(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) && (apiErr.Code == -352 || apiErr.Code == -412) {
		return true
	}
	var httpErr *HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == 412
}
