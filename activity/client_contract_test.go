package activity_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/activity"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/contracttest"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestActivityInfoUsesPromotedContractAndDecodesFixture(t *testing.T) {
	t.Parallel()

	body := contracttest.Fixture(t, "activity", "info", "responses", "anonymous.success.json")
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet || request.URL.String() != "https://api.bilibili.com/x/activity/subject/info?bvid=BV1mKY4e8ELy&sid=4017552" {
			t.Errorf("request = %s %s", request.Method, request.URL)
		}
		return testutil.JSONResponse(http.StatusOK, string(body)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	params, _ := activity.NewInfoParams(4_017_552)
	bvid, _ := ids.NewBVID("BV1mKY4e8ELy")
	params, _ = params.WithBVID(bvid)
	info, err := client.Activity().Info(context.Background(), params)
	if err != nil {
		t.Fatalf("Activity().Info() error = %v", err)
	}
	if info.ID != 4_017_552 || info.LID == nil || info.Name == "" {
		t.Fatalf("info = %+v, want promoted payload", info)
	}
}

func TestActivityListUsesPromotedContractAndDecodesAllProfiles(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "activity", "list", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				query := request.URL.Query()
				if request.Method != http.MethodGet || request.URL.Path != "/x/activity/page/list" || query.Get("plat") != "1,3" || query.Get("ps") != "1" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := activity.NewListParams().WithPageSize(1)
			page, err := client.Activity().List(context.Background(), params)
			if err != nil {
				t.Fatalf("Activity().List() error = %v", err)
			}
			if page.Page != 1 || page.Size != 1 || len(page.Items) != 1 {
				t.Fatalf("page = %+v, want one promoted item", page)
			}
		})
	}
}

func TestActivityListDefaultAndAPIErrorsUsePublicSurface(t *testing.T) {
	t.Parallel()

	var calls int
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			if request.URL.Query().Get("ps") != "15" || request.URL.Query().Get("pn") != "1" {
				t.Errorf("default query = %v", request.URL.Query())
			}
			return testutil.JSONResponse(http.StatusOK, string(contracttest.Fixture(t, "activity", "list", "responses", "anonymous.success.json"))), nil
		}
		return testutil.JSONResponse(http.StatusOK, `{"code":-101,"message":"not logged in"}`), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if _, err := client.Activity().ListDefault(context.Background()); err != nil {
		t.Fatalf("ListDefault() error = %v", err)
	}
	params, _ := activity.NewInfoParams(1)
	_, err = client.Activity().Info(context.Background(), params)
	if !bpi.RequiresLogin(err) {
		t.Fatalf("Info() error = %v, want login API error", err)
	}
	var apiErr *bpi.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != -101 {
		t.Fatalf("Info() error = %v, want APIError -101", err)
	}
}
