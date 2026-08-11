package opus_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/opus"
)

func TestOpusSpaceFeedUsesPromotedContractAndFixture(t *testing.T) {
	t.Parallel()

	body := contracttest.Fixture(t, "opus", "space-read", "space-feed", "responses", "success.json")
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		if request.Method != http.MethodGet || request.URL.Path != "/x/polymer/web-dynamic/v1/opus/feed/space" || query.Get("host_mid") != "1000001" || query.Get("page") != "0" || query.Get("type") != "all" || query.Get("web_location") != "333.1387" {
			t.Errorf("request = %s %s", request.Method, request.URL)
		}
		return testutil.JSONResponse(http.StatusOK, string(body)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	mid, _ := ids.NewMID(1_000_001)
	params, _ := opus.NewSpaceFeedParams(mid)
	feed, err := client.Opus().SpaceFeed(context.Background(), params)
	if err != nil {
		t.Fatalf("Opus().SpaceFeed() error = %v", err)
	}
	if !feed.HasMore || len(feed.Items) != 2 || feed.Items[0].Cover == nil || feed.Items[1].Stat.View == nil {
		t.Fatalf("SpaceFeed = %+v, want promoted payload", feed)
	}
}

func TestOpusSpaceFeedUsesTypedErrors(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		return testutil.JSONResponse(http.StatusOK, `{"code":-400,"message":"bad request"}`), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	mid, _ := ids.NewMID(1)
	params, _ := opus.NewSpaceFeedParams(mid)
	_, err = client.Opus().SpaceFeed(context.Background(), params)
	var apiError *bpi.APIError
	if !errors.As(err, &apiError) || apiError.Code != -400 {
		t.Fatalf("Opus().SpaceFeed() error = %v, want APIError -400", err)
	}

	_, err = client.Opus().SpaceFeed(context.Background(), opus.SpaceFeedParams{})
	var parameterError *bpi.ParameterError
	if !errors.As(err, &parameterError) || parameterError.Field != "host_mid" {
		t.Fatalf("zero SpaceFeedParams error = %v, want host_mid ParameterError", err)
	}
}
