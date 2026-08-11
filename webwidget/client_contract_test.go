package webwidget_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/webwidget"
)

func TestWebWidgetRegionBannerUsesPromotedFixtures(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "web_widget", "region-banner", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.String() != "https://api.bilibili.com/x/web-show/region/banner?region_id=1005" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := webwidget.NewRegionBannerParams(1005)
			data, err := client.WebWidget().RegionBanner(context.Background(), params)
			if err != nil {
				t.Fatalf("WebWidget().RegionBanner() error = %v", err)
			}
			if len(data.Items) == 0 || data.Items[0].RegionID != 1005 {
				t.Fatalf("RegionBannerData = %+v, want promoted items", data)
			}
		})
	}
}

func TestWebWidgetHeaderPageUsesPromotedFixtures(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "web_widget", "header-page", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.String() != "https://api.bilibili.com/x/web-show/page/header?resource_id=142" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			data, err := client.WebWidget().HeaderPage(context.Background(), webwidget.NewHeaderPageParams())
			if err != nil {
				t.Fatalf("WebWidget().HeaderPage() error = %v", err)
			}
			if data.IsSplitLayer != 1 || data.ParsedSplitLayer == nil || len(data.ParsedSplitLayer.Layers) == 0 {
				t.Fatalf("HeaderData = %+v, want parsed promoted layers", data)
			}
		})
	}
}

func TestWebWidgetOnlineUsesPromotedFixtures(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "web_widget", "online", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.String() != "https://api.bilibili.com/x/web-interface/online" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			data, err := client.WebWidget().Online(context.Background())
			if err != nil {
				t.Fatalf("WebWidget().Online() error = %v", err)
			}
			if len(data.RegionCount) == 0 {
				t.Fatal("OnlineData.RegionCount is empty")
			}
		})
	}
}

func TestWebWidgetReturnsTypedAPIErrors(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return testutil.JSONResponse(http.StatusOK, `{"code":-101,"message":"not logged in"}`), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.WebWidget().Online(context.Background())
	var apiError *bpi.APIError
	if !errors.As(err, &apiError) || apiError.Code != -101 || !bpi.RequiresLogin(err) {
		t.Fatalf("WebWidget().Online() error = %v, want login APIError", err)
	}
}
