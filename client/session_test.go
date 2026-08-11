package client

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestClientAbsorbsBilibiliResponseCookies(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	var calls atomic.Int32
	client, err := NewClient(
		WithClock(func() time.Time { return fixed }),
		WithAccount(Account{
			DedeUserID: "42",
			SESSDATA:   "old-session",
			BiliJCT:    "old-csrf",
			Buvid3:     "old-buvid",
		}),
		WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			switch calls.Add(1) {
			case 1:
				response := testutil.JSONResponse(http.StatusOK, `{"code":0}`)
				response.Header.Add("Set-Cookie", "SESSDATA=new-session; Path=/; HttpOnly")
				response.Header.Add("Set-Cookie", "bili_jct=new-csrf; Path=/")
				return response, nil
			default:
				cookie := request.Header.Get("Cookie")
				if !strings.Contains(cookie, "SESSDATA=new-session") || !strings.Contains(cookie, "bili_jct=new-csrf") {
					t.Errorf("second request Cookie = %q, want absorbed response cookies", cookie)
				}
				return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
			}
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	for range 2 {
		request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
		if _, err := client.Do(context.Background(), request, "test.cookie_refresh"); err != nil {
			t.Fatalf("Do() error = %v", err)
		}
	}
	account, ok := client.Account()
	if !ok || account.SESSDATA != "new-session" || account.BiliJCT != "new-csrf" {
		t.Fatalf("Account() = %v, %v; want refreshed complete account", account, ok)
	}
}

func TestSessionRemovesExpiredResponseCookies(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	session := newSession()
	if err := session.replaceCookieHeader("DedeUserID=42; SESSDATA=session; bili_jct=csrf; buvid3=buvid"); err != nil {
		t.Fatalf("replaceCookieHeader() error = %v", err)
	}
	session.absorbResponse("api.bilibili.com", []*http.Cookie{{Name: "SESSDATA", Value: "stale", Expires: now}}, now)
	if session.hasLoginCookies() {
		t.Fatal("expired SESSDATA remained in session")
	}
	if _, ok := session.getAccount(); ok {
		t.Fatal("expired SESSDATA left a complete account")
	}
}

func TestRawCookieHeaderSupportsPartialAndUnknownPairs(t *testing.T) {
	t.Parallel()

	var sentCookie string
	client, err := NewClient(
		WithCookie("SESSDATA=session; bili_jct=csrf; DedeUserID__ckMd5=device-check"),
		WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			sentCookie = request.Header.Get("Cookie")
			return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if csrf, err := client.CSRF(); err != nil || csrf != "csrf" {
		t.Fatalf("CSRF() = %q, %v; want csrf, nil", csrf, err)
	}
	if _, ok := client.Account(); ok {
		t.Fatal("Account() returned a complete account for a partial Cookie header")
	}
	request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	if _, err := client.Do(context.Background(), request, "test.raw_cookie"); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	for _, pair := range []string{"SESSDATA=session", "bili_jct=csrf", "DedeUserID__ckMd5=device-check"} {
		if !strings.Contains(sentCookie, pair) {
			t.Fatalf("Cookie = %q, want pair %q", sentCookie, pair)
		}
	}
}

func TestClientIgnoresExternalResponseCookies(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	client, err := NewClient(
		WithCookie("DedeUserID=42; SESSDATA=original-session; bili_jct=csrf; buvid3=buvid"),
		WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			if calls.Add(1) == 1 {
				if got := request.Header.Get("Cookie"); got != "caller-cookie=allowed" {
					t.Errorf("external explicit Cookie = %q, want caller Cookie only", got)
				}
				response := testutil.JSONResponse(http.StatusOK, `{"code":0}`)
				response.Header.Add("Set-Cookie", "SESSDATA=external-session; Path=/")
				return response, nil
			}
			cookie := request.Header.Get("Cookie")
			if !strings.Contains(cookie, "SESSDATA=original-session") || strings.Contains(cookie, "external-session") {
				t.Errorf("Bilibili Cookie = %q, external response changed session", cookie)
			}
			return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	external, _ := http.NewRequest(http.MethodGet, "https://example.com/x/test", nil)
	external.Header.Set("Cookie", "caller-cookie=allowed")
	if _, err := client.Do(context.Background(), external, "test.external_cookie"); err != nil {
		t.Fatalf("Do(external) error = %v", err)
	}
	bilibili, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test", nil)
	if _, err := client.Do(context.Background(), bilibili, "test.bilibili_cookie"); err != nil {
		t.Fatalf("Do(Bilibili) error = %v", err)
	}
}

func TestResponseCookieUpdatesAreSafeDuringConcurrentRequests(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithAccount(Account{DedeUserID: "42", SESSDATA: "session", BiliJCT: "csrf", Buvid3: "buvid"}),
		WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			response := testutil.JSONResponse(http.StatusOK, `{"code":0}`)
			response.Header.Add("Set-Cookie", "buvid3="+request.URL.Query().Get("value")+"; Path=/")
			return response, nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	var group sync.WaitGroup
	for index := range 64 {
		group.Add(1)
		go func() {
			defer group.Done()
			request, _ := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/test?value="+string(rune('a'+index%26)), nil)
			if _, err := client.Do(context.Background(), request, "test.concurrent_cookie"); err != nil {
				t.Errorf("Do() error = %v", err)
			}
		}()
	}
	group.Wait()
	if account, ok := client.Account(); !ok || account.Buvid3 == "" {
		t.Fatalf("Account() = %v, %v; want complete account after concurrent updates", account, ok)
	}
}

func FuzzParseCookieHeader(f *testing.F) {
	for _, seed := range []string{
		"DedeUserID=42; SESSDATA=session; bili_jct=csrf; buvid3=buvid",
		"name=value",
		"name=",
		"",
		"missing-equals",
		"=missing-name",
		"bad name=value",
		"name=value\r\nInjected: true",
		";;;",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, header string) {
		cookies, err := parseCookieHeader(header)
		if err != nil {
			var parameterErr *ParameterError
			if !errors.As(err, &parameterErr) || parameterErr.Field != "cookie" {
				t.Fatalf("parseCookieHeader(%q) error = %T %v, want cookie ParameterError", header, err, err)
			}
			return
		}
		if len(cookies) == 0 {
			t.Fatalf("parseCookieHeader(%q) succeeded with no cookies", header)
		}
		for name, value := range cookies {
			if err := (&http.Cookie{Name: name, Value: value}).Valid(); err != nil {
				t.Fatalf("parseCookieHeader(%q) returned invalid cookie %q: %v", header, name, err)
			}
		}
	})
}
