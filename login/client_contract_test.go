package login_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestLoginNavMatchesAuthenticatedPromotedFixtures(t *testing.T) {
	t.Parallel()

	for _, profile := range []struct {
		name string
		mid  uint64
	}{
		{name: "normal", mid: 1_000_001},
		{name: "vip", mid: 1_000_002},
	} {
		profile := profile
		t.Run(profile.name, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "login", "nav", "responses", profile.name+".success.json")
			client, err := bpi.NewClient(
				bpi.WithCookie("DedeUserID=42; SESSDATA=session; bili_jct=csrf; buvid3=buvid"),
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					if request.Method != http.MethodGet || request.URL.String() != "https://api.bilibili.com/x/web-interface/nav" {
						t.Errorf("request = %s %s", request.Method, request.URL)
					}
					if !strings.Contains(request.Header.Get("Cookie"), "SESSDATA=session") {
						t.Errorf("Cookie = %q, want authenticated session", request.Header.Get("Cookie"))
					}
					return testutil.JSONResponse(http.StatusOK, string(body)), nil
				})}),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			nav, err := client.Login().Nav(context.Background())
			if err != nil {
				t.Fatalf("Login().Nav() error = %v", err)
			}
			if !nav.IsLogin || nav.MID == nil || nav.MID.Uint64() != profile.mid || nav.Username == nil {
				t.Fatalf("Nav() = %+v, want authenticated %s profile", nav, profile.name)
			}
		})
	}
}

func TestLoginNavAnonymousFixtureReturnsRequiresLogin(t *testing.T) {
	t.Parallel()

	body := contracttest.Fixture(t, "login", "nav", "responses", "anonymous.error.json")
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Cookie") != "" {
			t.Errorf("anonymous Cookie = %q, want empty", request.Header.Get("Cookie"))
		}
		return testutil.JSONResponse(http.StatusOK, string(body)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.Login().Nav(context.Background())
	if !bpi.RequiresLogin(err) {
		t.Fatalf("Login().Nav() error = %v, want RequiresLogin", err)
	}
}
