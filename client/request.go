package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Response contains a fully buffered successful HTTP response. Body is an
// explicit raw-data surface and may contain sensitive information.
type Response struct {
	StatusCode int
	Body       []byte
	// Cookies contains a detached copy of Set-Cookie values returned by the
	// endpoint. Client has already incorporated trusted Bilibili cookies into
	// its isolated session before returning this response.
	Cookies  []http.Cookie
	Duration time.Duration
}

// Do executes request with the client's headers, session, logging, response
// limit, and HTTP error handling. The caller retains ownership of request;
// Client executes a clone.
func (c *Client) Do(ctx context.Context, request *http.Request, operation string) (*Response, error) {
	if ctx == nil {
		return nil, &ParameterError{Field: "context", Message: "context cannot be nil"}
	}
	if request == nil {
		return nil, &ParameterError{Field: "request", Message: "request cannot be nil"}
	}
	if request.URL == nil {
		return nil, &ParameterError{Field: "request_url", Message: "request URL cannot be nil"}
	}
	if strings.TrimSpace(operation) == "" {
		return nil, &ParameterError{Field: "operation", Message: "operation cannot be empty"}
	}

	clone := request.Clone(ctx)
	clone.Header = request.Header.Clone()
	if clone.Header.Get("User-Agent") == "" {
		clone.Header.Set("User-Agent", c.userAgent)
	}
	if isBilibiliHost(clone.URL.Hostname()) {
		if clone.Header.Get("Referer") == "" {
			clone.Header.Set("Referer", c.referer)
		}
		if clone.Header.Get("Origin") == "" {
			clone.Header.Set("Origin", c.origin)
		}
		c.session.addToRequest(clone)
	}

	safeURL := sanitizeURL(clone.URL)
	c.logger.InfoContext(ctx, "sending Bilibili request",
		"operation", operation,
		"method", clone.Method,
		"url", safeURL,
	)

	started := time.Now()
	httpResponse, err := c.httpClient.Do(clone)
	if err != nil {
		transportErr := &TransportError{Operation: operation, Err: err}
		c.logger.ErrorContext(ctx, "Bilibili request failed",
			"operation", operation,
			"method", clone.Method,
			"url", safeURL,
			"error", transportErr,
		)
		return nil, transportErr
	}
	defer httpResponse.Body.Close()
	duration := time.Since(started)
	responseCookiePointers := httpResponse.Cookies()
	c.session.absorbResponse(clone.URL.Hostname(), responseCookiePointers, c.now())
	responseCookies := make([]http.Cookie, 0, len(responseCookiePointers))
	for _, cookie := range responseCookiePointers {
		if cookie != nil {
			responseCookies = append(responseCookies, *cookie)
		}
	}

	body, err := readLimitedBody(httpResponse.Body, c.maxResponseBody)
	if err != nil {
		c.logger.ErrorContext(ctx, "Bilibili response body failed",
			"operation", operation,
			"status", httpResponse.StatusCode,
			"duration_ms", duration.Milliseconds(),
			"error", err,
		)
		return nil, err
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		result := &HTTPError{StatusCode: httpResponse.StatusCode}
		c.logger.ErrorContext(ctx, "Bilibili request returned HTTP error",
			"operation", operation,
			"status", httpResponse.StatusCode,
			"duration_ms", duration.Milliseconds(),
		)
		return nil, result
	}

	c.logger.InfoContext(ctx, "Bilibili request completed",
		"operation", operation,
		"status", httpResponse.StatusCode,
		"duration_ms", duration.Milliseconds(),
	)
	return &Response{StatusCode: httpResponse.StatusCode, Body: body, Cookies: responseCookies, Duration: duration}, nil
}

// SendPayload executes request and returns a required Bilibili business
// payload. It is the generic custom-request counterpart to domain methods.
func SendPayload[T any](ctx context.Context, client *Client, request *http.Request, operation string) (T, error) {
	var zero T
	if client == nil {
		return zero, &ParameterError{Field: "client", Message: "client cannot be nil"}
	}
	response, err := client.Do(ctx, request, operation)
	if err != nil {
		return zero, err
	}
	envelope, err := DecodeEnvelope[T](response.Body)
	if err != nil {
		client.logger.ErrorContext(ctx, "Bilibili response decode failed",
			"operation", operation,
			"error", err,
		)
		return zero, err
	}
	client.logger.InfoContext(ctx, "Bilibili response decoded",
		"operation", operation,
		"api_code", envelope.Code,
	)
	payload, err := envelope.IntoPayload()
	if err != nil {
		client.logger.ErrorContext(ctx, "Bilibili interface returned an error",
			"operation", operation,
			"api_code", envelope.Code,
		)
	}
	return payload, err
}

// SendOptionalPayload executes request and returns an optional Bilibili
// business payload.
func SendOptionalPayload[T any](ctx context.Context, client *Client, request *http.Request, operation string) (*T, error) {
	if client == nil {
		return nil, &ParameterError{Field: "client", Message: "client cannot be nil"}
	}
	response, err := client.Do(ctx, request, operation)
	if err != nil {
		return nil, err
	}
	envelope, err := DecodeEnvelope[T](response.Body)
	if err != nil {
		client.logger.ErrorContext(ctx, "Bilibili response decode failed",
			"operation", operation,
			"error", err,
		)
		return nil, err
	}
	client.logger.InfoContext(ctx, "Bilibili response decoded",
		"operation", operation,
		"api_code", envelope.Code,
	)
	payload, err := envelope.IntoOptionalPayload()
	if err != nil {
		client.logger.ErrorContext(ctx, "Bilibili interface returned an error",
			"operation", operation,
			"api_code", envelope.Code,
		)
	}
	return payload, err
}

func readLimitedBody(reader io.Reader, limit int64) ([]byte, error) {
	limited := io.LimitReader(reader, limit+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, &TransportError{Operation: "read response body", Err: err}
	}
	if int64(len(body)) > limit {
		return nil, &ResponseTooLargeError{Limit: limit}
	}
	return body, nil
}

func sanitizeURL(input *url.URL) string {
	if input == nil {
		return "<nil-url>"
	}
	clone := *input
	query := clone.Query()
	for key := range query {
		if isSensitiveQueryKey(key) {
			query.Del(key)
		}
	}
	clone.RawQuery = query.Encode()
	clone.User = nil
	return clone.String()
}

func isSensitiveQueryKey(key string) bool {
	switch strings.ToLower(key) {
	case "sessdata", "dedeuserid", "dedeuserid__ckmd5", "bili_jct", "csrf", "csrf_token", "w_rid", "access_key", "token", "cookie":
		return true
	default:
		return false
	}
}
