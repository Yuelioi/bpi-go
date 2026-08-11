package bpi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	core "github.com/Yuelioi/bpi-go/client"
)

// Account contains the four common Cookie values used for a complete Bilibili
// account projection. It retains the client module's redacted formatting
// semantics.
type Account = core.Account

// Option configures a Client before it is constructed.
type Option = core.Option

// Response contains a fully buffered successful HTTP response.
type Response = core.Response

// Envelope is the common Bilibili JSON response wrapper.
type Envelope[T any] = core.Envelope[T]

// ParameterError describes an invalid caller-supplied value.
type ParameterError = core.ParameterError

// TransportError wraps a failure from the configured HTTP transport.
type TransportError = core.TransportError

// HTTPError reports a non-successful HTTP status.
type HTTPError = core.HTTPError

// APIError reports a non-zero Bilibili response code.
type APIError = core.APIError

// ResponseDecodeError retains a recoverable, explicitly accessed response body.
type ResponseDecodeError = core.ResponseDecodeError

// ResponseTooLargeError reports that a response exceeded the configured limit.
type ResponseTooLargeError = core.ResponseTooLargeError

var (
	// ErrMissingData means a successful response omitted a required payload.
	ErrMissingData = core.ErrMissingData
	// ErrAuthenticationRequired means an operation needs account credentials.
	ErrAuthenticationRequired = core.ErrAuthenticationRequired
)

// Client is the stable root facade over the shared low-level client module.
// It is isolated and safe for concurrent use.
type Client struct {
	inner *core.Client
}

// NewClient constructs a client without reading files, mutating global state,
// or performing network I/O.
func NewClient(options ...Option) (*Client, error) {
	inner, err := core.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return &Client{inner: inner}, nil
}

func (c *Client) core() *core.Client {
	if c == nil {
		return nil
	}
	return c.inner
}

// SetAccount atomically replaces the client's authenticated session.
func (c *Client) SetAccount(account Account) error { return c.inner.SetAccount(account) }

// SetCookie atomically replaces the client's session from a raw HTTP Cookie
// request-header value.
func (c *Client) SetCookie(cookieHeader string) error { return c.inner.SetCookie(cookieHeader) }

// ClearAccount removes all client session values.
func (c *Client) ClearAccount() { c.inner.ClearAccount() }

// Account returns a copy of the current complete account, when available.
func (c *Client) Account() (Account, bool) { return c.inner.Account() }

// HasLoginCookies reports whether the session has a non-empty SESSDATA value.
func (c *Client) HasLoginCookies() bool { return c.inner.HasLoginCookies() }

// CSRF returns the current bili_jct Cookie value.
func (c *Client) CSRF() (string, error) { return c.inner.CSRF() }

// Do executes request through the shared bounded, credential-scoped transport.
func (c *Client) Do(ctx context.Context, request *http.Request, operation string) (*Response, error) {
	if c == nil {
		return nil, &ParameterError{Field: "client", Message: "client cannot be nil"}
	}
	return c.inner.Do(ctx, request, operation)
}

// WithHTTPClient supplies the HTTP adapter used by the client.
func WithHTTPClient(client *http.Client) Option { return core.WithHTTPClient(client) }

// WithLogger enables sanitized structured request logging.
func WithLogger(logger *slog.Logger) Option { return core.WithLogger(logger) }

// WithTimeout sets the total HTTP request timeout.
func WithTimeout(timeout time.Duration) Option { return core.WithTimeout(timeout) }

// WithUserAgent changes the default User-Agent header.
func WithUserAgent(userAgent string) Option { return core.WithUserAgent(userAgent) }

// WithReferer changes the default Referer header for Bilibili requests.
func WithReferer(referer string) Option { return core.WithReferer(referer) }

// WithOrigin changes the default Origin header for Bilibili requests.
func WithOrigin(origin string) Option { return core.WithOrigin(origin) }

// WithMaxResponseBody sets the maximum response body buffered in memory.
func WithMaxResponseBody(limit int64) Option { return core.WithMaxResponseBody(limit) }

// WithCookie initializes the client from a raw HTTP Cookie request-header
// value. The header may contain any valid Cookie pairs.
func WithCookie(cookieHeader string) Option { return core.WithCookie(cookieHeader) }

// WithAccount initializes the client from a complete structured account.
func WithAccount(account Account) Option { return core.WithAccount(account) }

// SendPayload executes request and returns a required business payload.
func SendPayload[T any](ctx context.Context, client *Client, request *http.Request, operation string) (T, error) {
	if client == nil {
		var zero T
		return zero, &ParameterError{Field: "client", Message: "client cannot be nil"}
	}
	return core.SendPayload[T](ctx, client.inner, request, operation)
}

// SendOptionalPayload executes request and returns an optional business payload.
func SendOptionalPayload[T any](ctx context.Context, client *Client, request *http.Request, operation string) (*T, error) {
	if client == nil {
		return nil, &ParameterError{Field: "client", Message: "client cannot be nil"}
	}
	return core.SendOptionalPayload[T](ctx, client.inner, request, operation)
}

// DecodeEnvelope decodes the common Bilibili response envelope.
func DecodeEnvelope[T any](body []byte) (Envelope[T], error) { return core.DecodeEnvelope[T](body) }

// RequiresLogin reports whether err represents an unauthenticated response.
func RequiresLogin(err error) bool { return core.RequiresLogin(err) }

// RequiresVIP reports whether err represents a VIP-only response.
func RequiresVIP(err error) bool { return core.RequiresVIP(err) }

// IsPermissionError reports whether err represents an authorization failure.
func IsPermissionError(err error) bool { return core.IsPermissionError(err) }

// IsRiskControl reports whether err represents Bilibili risk control.
func IsRiskControl(err error) bool { return core.IsRiskControl(err) }
