package client

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Option configures a Client before it is constructed.
type Option interface {
	apply(*clientConfig) error
}

type optionFunc func(*clientConfig) error

func (f optionFunc) apply(config *clientConfig) error { return f(config) }

type clientConfig struct {
	httpClient      *http.Client
	logger          *slog.Logger
	timeout         time.Duration
	timeoutSet      bool
	userAgent       string
	referer         string
	origin          string
	maxResponseBody int64
	cookie          string
	account         *Account
	now             func() time.Time
}

// WithClock injects the clock used by time-dependent request signatures. Most
// callers should use the default system clock; this option exists for
// deterministic adapters and tests.
func WithClock(now func() time.Time) Option {
	return optionFunc(func(config *clientConfig) error {
		if now == nil {
			return &ParameterError{Field: "clock", Message: "clock cannot be nil"}
		}
		config.now = now
		return nil
	})
}

// WithHTTPClient uses a shallow copy of client for transport, redirects, and
// timeout policy. Cookie state is always managed independently by bpi.
func WithHTTPClient(client *http.Client) Option {
	return optionFunc(func(config *clientConfig) error {
		if client == nil {
			return &ParameterError{Field: "http_client", Message: "client cannot be nil"}
		}
		clone := *client
		clone.Jar = nil
		config.httpClient = &clone
		return nil
	})
}

// WithLogger enables sanitized structured request logging. A nil logger is
// rejected; clients are quiet by default.
func WithLogger(logger *slog.Logger) Option {
	return optionFunc(func(config *clientConfig) error {
		if logger == nil {
			return &ParameterError{Field: "logger", Message: "logger cannot be nil"}
		}
		config.logger = logger
		return nil
	})
}

// WithTimeout sets the total HTTP request timeout.
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(config *clientConfig) error {
		if timeout <= 0 {
			return &ParameterError{Field: "timeout", Message: "timeout must be positive"}
		}
		config.timeout = timeout
		config.timeoutSet = true
		return nil
	})
}

// WithUserAgent changes the default User-Agent header.
func WithUserAgent(userAgent string) Option {
	return headerOption("user_agent", userAgent, func(config *clientConfig, value string) {
		config.userAgent = value
	})
}

// WithReferer changes the default Referer header for Bilibili requests.
func WithReferer(referer string) Option {
	return headerOption("referer", referer, func(config *clientConfig, value string) {
		config.referer = value
	})
}

// WithOrigin changes the default Origin header for Bilibili requests.
func WithOrigin(origin string) Option {
	return headerOption("origin", origin, func(config *clientConfig, value string) {
		config.origin = value
	})
}

// WithMaxResponseBody sets the maximum response body buffered in memory.
func WithMaxResponseBody(limit int64) Option {
	return optionFunc(func(config *clientConfig) error {
		if limit <= 0 {
			return &ParameterError{Field: "max_response_body", Message: "limit must be positive"}
		}
		config.maxResponseBody = limit
		return nil
	})
}

// WithCookie initializes the client from a raw HTTP Cookie request-header
// value. All valid Cookie pairs are preserved, including pairs unknown to bpi.
func WithCookie(cookieHeader string) Option {
	return optionFunc(func(config *clientConfig) error {
		if config.account != nil {
			return &ParameterError{Field: "credentials", Message: "cookie and account options are mutually exclusive"}
		}
		if strings.TrimSpace(cookieHeader) == "" {
			return &ParameterError{Field: "cookie", Message: "cookie cannot be empty"}
		}
		config.cookie = cookieHeader
		return nil
	})
}

// WithAccount initializes the client from a complete structured account.
func WithAccount(account Account) Option {
	return optionFunc(func(config *clientConfig) error {
		if config.cookie != "" {
			return &ParameterError{Field: "credentials", Message: "cookie and account options are mutually exclusive"}
		}
		if err := account.Validate(); err != nil {
			return err
		}
		clone := account
		config.account = &clone
		return nil
	})
}

func headerOption(field, value string, set func(*clientConfig, string)) Option {
	return optionFunc(func(config *clientConfig) error {
		if strings.TrimSpace(value) == "" {
			return &ParameterError{Field: field, Message: "header value cannot be empty"}
		}
		if strings.ContainsAny(value, "\r\n") {
			return &ParameterError{Field: field, Message: "header value cannot contain a newline"}
		}
		set(config, value)
		return nil
	})
}
