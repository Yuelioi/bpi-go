package client

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultTimeout         = 10 * time.Second
	defaultMaxResponseBody = 16 << 20
	defaultUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	defaultReferer         = "https://www.bilibili.com/"
	defaultOrigin          = "https://www.bilibili.com"
)

// Client is an isolated, concurrency-safe Bilibili client.
type Client struct {
	httpClient      *http.Client
	logger          *slog.Logger
	session         *sessionState
	userAgent       string
	referer         string
	origin          string
	maxResponseBody int64
	now             func() time.Time
	wbiKeys         *wbiKeyCache
}

// NewClient constructs a client without reading files, mutating global state,
// or performing network I/O.
func NewClient(options ...Option) (*Client, error) {
	config := clientConfig{
		timeout:         defaultTimeout,
		userAgent:       defaultUserAgent,
		referer:         defaultReferer,
		origin:          defaultOrigin,
		maxResponseBody: defaultMaxResponseBody,
		now:             time.Now,
	}
	for index, option := range options {
		if option == nil {
			return nil, &ParameterError{Field: "option", Message: "option at index " + strconv.Itoa(index) + " is nil"}
		}
		if err := option.apply(&config); err != nil {
			return nil, err
		}
	}

	httpClient := config.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: defaultTransport(),
			Timeout:   config.timeout,
		}
	} else if config.timeoutSet {
		httpClient.Timeout = config.timeout
	}
	if httpClient.Transport == nil {
		httpClient.Transport = defaultTransport()
	}

	logger := config.logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	session := newSession()
	if config.cookie != "" {
		if err := session.replaceCookieHeader(config.cookie); err != nil {
			return nil, err
		}
	}
	if config.account != nil {
		if err := session.replaceAccount(*config.account); err != nil {
			return nil, err
		}
	}

	return &Client{
		httpClient:      httpClient,
		logger:          logger,
		session:         session,
		userAgent:       config.userAgent,
		referer:         config.referer,
		origin:          config.origin,
		maxResponseBody: config.maxResponseBody,
		now:             config.now,
		wbiKeys:         newWBIKeyCache(),
	}, nil
}

// SetAccount atomically replaces the client's authenticated session.
func (c *Client) SetAccount(account Account) error {
	return c.session.replaceAccount(account)
}

// SetCookie atomically replaces the client's session from a raw HTTP Cookie
// request-header value. All valid Cookie pairs are preserved.
func (c *Client) SetCookie(cookieHeader string) error {
	return c.session.replaceCookieHeader(cookieHeader)
}

// ClearAccount removes all client session values.
func (c *Client) ClearAccount() {
	c.session.clear()
}

// Account returns a copy of the current complete account, when available.
func (c *Client) Account() (Account, bool) {
	return c.session.getAccount()
}

// HasLoginCookies reports whether the client has a non-empty SESSDATA Cookie.
func (c *Client) HasLoginCookies() bool {
	return c.session.hasLoginCookies()
}

// CSRF returns the current bili_jct Cookie value. It does not require the
// session to contain every field needed for a complete Account projection.
func (c *Client) CSRF() (string, error) {
	return c.session.csrf()
}

// Now returns the client's configured clock value. Domain modules use it for
// request signatures whose timestamp must share the client's testable clock.
func (c *Client) Now() time.Time { return c.now() }

// Referer returns the default Referer header used for Bilibili requests.
func (c *Client) Referer() string { return c.referer }

// Origin returns the default Origin header used for Bilibili requests.
func (c *Client) Origin() string { return c.origin }

func defaultTransport() http.RoundTripper {
	if transport, ok := http.DefaultTransport.(*http.Transport); ok {
		return transport.Clone()
	}
	return http.DefaultTransport
}
