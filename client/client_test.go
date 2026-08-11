package client_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestNewClientIsOfflineAndIsolated(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	httpClient := &http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		requests.Add(1)
		return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
	})}
	first, err := bpi.NewClient(
		bpi.WithHTTPClient(httpClient),
		bpi.WithAccount(completeAccount()),
	)
	if err != nil {
		t.Fatalf("NewClient(first) error = %v", err)
	}
	second, err := bpi.NewClient(bpi.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient(second) error = %v", err)
	}
	if requests.Load() != 0 {
		t.Fatalf("NewClient performed %d requests, want 0", requests.Load())
	}
	if !first.HasLoginCookies() {
		t.Fatal("first.HasLoginCookies() = false, want true")
	}
	if second.HasLoginCookies() {
		t.Fatal("second.HasLoginCookies() = true, want false")
	}
	first.ClearAccount()
	if first.HasLoginCookies() {
		t.Fatal("first.HasLoginCookies() remained true after ClearAccount")
	}
}

func TestDoAppliesBilibiliHeadersAndScopesCredentials(t *testing.T) {
	t.Parallel()

	var requests []*http.Request
	client, err := bpi.NewClient(
		bpi.WithCookie("DedeUserID=42; SESSDATA=session; bili_jct=csrf; buvid3=buvid"),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			requests = append(requests, request.Clone(context.Background()))
			return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	bilibiliRequest, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	if _, err := client.Do(context.Background(), bilibiliRequest, "test.bilibili"); err != nil {
		t.Fatalf("Do(Bilibili) error = %v", err)
	}
	externalRequest, _ := http.NewRequest(http.MethodGet, "https://example.com/x/test", nil)
	if _, err := client.Do(context.Background(), externalRequest, "test.external"); err != nil {
		t.Fatalf("Do(external) error = %v", err)
	}

	if len(requests) != 2 {
		t.Fatalf("request count = %d, want 2", len(requests))
	}
	if requests[0].Header.Get("User-Agent") == "" || requests[0].Header.Get("Referer") == "" || requests[0].Header.Get("Origin") == "" {
		t.Fatalf("Bilibili headers = %v, want user-agent, referer, and origin", requests[0].Header)
	}
	if cookie := requests[0].Header.Get("Cookie"); !strings.Contains(cookie, "SESSDATA=session") {
		t.Fatalf("Bilibili Cookie = %q, want session", cookie)
	}
	if cookie := requests[1].Header.Get("Cookie"); cookie != "" {
		t.Fatalf("external Cookie = %q, want empty", cookie)
	}
	if requests[1].Header.Get("Referer") != "" || requests[1].Header.Get("Origin") != "" {
		t.Fatalf("external headers leaked Bilibili origin: %v", requests[1].Header)
	}
}

func TestDoPropagatesContextCancellation(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		<-request.Context().Done()
		return nil, request.Context().Err()
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.Do(ctx, request, "test.cancel")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do() error = %v, want context.Canceled", err)
	}
}

func TestDoSanitizesLoggedURL(t *testing.T) {
	t.Parallel()

	var logOutput bytes.Buffer
	client, err := bpi.NewClient(
		bpi.WithLogger(slog.New(slog.NewJSONHandler(&logOutput, nil))),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test?mid=42&SESSDATA=session-secret&csrf=csrf-secret&w_rid=signature-secret", nil)
	if _, err := client.Do(context.Background(), request, "test.log"); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	logged := logOutput.String()
	for _, secret := range []string{"session-secret", "csrf-secret", "signature-secret"} {
		if strings.Contains(logged, secret) {
			t.Fatalf("log leaked %q: %s", secret, logged)
		}
	}
	if !strings.Contains(logged, "mid=42") {
		t.Fatalf("log removed safe query value: %s", logged)
	}
}

func TestDoRedactsURLFromTransportErrorAndLog(t *testing.T) {
	t.Parallel()

	var logOutput bytes.Buffer
	client, err := bpi.NewClient(
		bpi.WithLogger(slog.New(slog.NewJSONHandler(&logOutput, nil))),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			return nil, &url.Error{Op: request.Method, URL: request.URL.String(), Err: errors.New("network unavailable")}
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test?SESSDATA=transport-secret", nil)
	_, err = client.Do(context.Background(), request, "test.transport")
	if err == nil {
		t.Fatal("Do() error = nil, want transport error")
	}
	if strings.Contains(err.Error(), "transport-secret") || strings.Contains(logOutput.String(), "transport-secret") {
		t.Fatalf("transport failure leaked URL: error=%v log=%s", err, logOutput.String())
	}
	var transportErr *bpi.TransportError
	if !errors.As(err, &transportErr) || transportErr.Unwrap() == nil {
		t.Fatalf("error = %v, want unwrap-capable TransportError", err)
	}
}

func TestDoEnforcesResponseLimit(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(
		bpi.WithMaxResponseBody(4),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("12345"))}, nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	_, err = client.Do(context.Background(), request, "test.limit")
	var tooLarge *bpi.ResponseTooLargeError
	if !errors.As(err, &tooLarge) || tooLarge.Limit != 4 {
		t.Fatalf("Do() error = %v, want 4-byte ResponseTooLargeError", err)
	}
}

func TestDoCancellationInterruptsResponseBodyRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	body := &cancelBody{ctx: ctx, started: make(chan struct{})}
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: body}, nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	result := make(chan error, 1)
	go func() {
		_, err := client.Do(ctx, request, "test.body_cancel")
		result <- err
	}()
	<-body.started
	cancel()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatalf("Do() error = %v, want context.Canceled", err)
	}
}

func TestSendPayloadUsesPublicClientSurface(t *testing.T) {
	t.Parallel()

	var logOutput bytes.Buffer
	client, err := bpi.NewClient(
		bpi.WithLogger(slog.New(slog.NewJSONHandler(&logOutput, nil))),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusOK, string(fixture(t, "success.json"))), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	payload, err := bpi.SendPayload[fixturePayload](context.Background(), client, request, "test.payload")
	if err != nil {
		t.Fatalf("SendPayload() error = %v", err)
	}
	if payload.AID != 170001 {
		t.Fatalf("payload.AID = %d, want 170001", payload.AID)
	}
	if !strings.Contains(logOutput.String(), `"api_code":0`) {
		t.Fatalf("SendPayload log omitted API code: %s", logOutput.String())
	}
}

func TestClientSessionIsSafeDuringConcurrentUse(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	var group sync.WaitGroup
	for index := range 32 {
		group.Add(1)
		go func() {
			defer group.Done()
			if index%3 == 0 {
				if err := client.SetAccount(completeAccount()); err != nil {
					t.Errorf("SetAccount() error = %v", err)
				}
			} else if index%3 == 1 {
				client.ClearAccount()
			} else {
				_, _ = client.Account()
				request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
				if _, err := client.Do(context.Background(), request, "test.concurrent"); err != nil {
					t.Errorf("Do() error = %v", err)
				}
			}
		}()
	}
	group.Wait()
}

type cancelBody struct {
	ctx     context.Context
	started chan struct{}
	once    sync.Once
}

func (b *cancelBody) Read([]byte) (int, error) {
	b.once.Do(func() { close(b.started) })
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (b *cancelBody) Close() error { return nil }
