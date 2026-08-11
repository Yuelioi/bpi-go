package wallet_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/wallet"
)

func TestWalletInfoUsesAuthenticatedProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					if request.Method != http.MethodPost || request.URL.Path != "/paywallet/wallet/getUserWallet" || request.Header.Get("Cookie") == "" || request.Header.Get("Content-Type") != "application/json" {
						t.Fatalf("request = %s %s headers %v", request.Method, request.URL, request.Header)
					}
					body, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("read body: %v", err)
					}
					var got map[string]any
					if err := json.Unmarshal(body, &got); err != nil || got["csrf"] != "fixture-csrf" || got["platformType"] != float64(3) || got["timestamp"] != float64(1_700_000_000_000) || got["traceId"] != float64(1_700_000_000_000) || got["version"] != "1.0" {
						t.Fatalf("body = %v, %v", got, err)
					}
					fixture := contracttest.Fixture(t, "wallet", "read", "info", "responses", "authenticated.success.json")
					return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
				})}),
				contracttest.ProfileOption(profile),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := wallet.InfoParamsAtMilliseconds(1_700_000_000_000)
			info, err := client.Wallet().Info(context.Background(), params)
			if err != nil || info.MID != 1_000_001 || info.ShowClassBalance != 1 {
				t.Fatalf("Info() = %+v, %v", info, err)
			}
		})
	}
}

func TestWalletInfoRequiresAuthenticationBeforeTransport(t *testing.T) {
	client, err := bpi.NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	params, _ := wallet.InfoParamsAtMilliseconds(1)
	_, err = client.Wallet().Info(context.Background(), params)
	if !errors.Is(err, bpi.ErrAuthenticationRequired) {
		t.Fatalf("Info() error = %v, want ErrAuthenticationRequired", err)
	}
}
