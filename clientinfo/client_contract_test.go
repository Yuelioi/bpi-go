package clientinfo_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/clientinfo"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestClientInfoIPUsesPromotedContractAndDecodesAllProfiles(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "clientinfo", "ip", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.String() != "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr?ip=8.8.8.8" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := clientinfo.NewIPParams().WithIP("8.8.8.8")
			info, err := client.ClientInfo().IP(context.Background(), params)
			if err != nil {
				t.Fatalf("ClientInfo().IP() error = %v", err)
			}
			if info.Address == nil || *info.Address != "8.8.8.8" {
				t.Fatalf("IPInfo.Address = %v, want 8.8.8.8", info.Address)
			}
		})
	}
}

func TestClientInfoIPUsesPublicTypedErrors(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return testutil.JSONResponse(http.StatusOK, `{"code":-400,"message":"bad request"}`), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.ClientInfo().IP(context.Background(), clientinfo.NewIPParams())
	var apiError *bpi.APIError
	if !errors.As(err, &apiError) || apiError.Code != -400 {
		t.Fatalf("ClientInfo().IP() error = %v, want APIError -400", err)
	}
	_, err = clientinfo.NewIPParams().WithIP("not-an-ip")
	var parameterError *bpi.ParameterError
	if !errors.As(err, &parameterError) || parameterError.Field != "ip" {
		t.Fatalf("WithIP() error = %v, want public ip ParameterError", err)
	}
}
