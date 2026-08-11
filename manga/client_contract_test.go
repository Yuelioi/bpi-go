package manga_test

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
	"github.com/Yuelioi/bpi-go/manga"
)

func TestMangaPromotedReadsUseContractRequestsAndProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			transport := testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodPost || request.URL.Host != "manga.bilibili.com" || request.URL.RawQuery != "" {
					t.Fatalf("request = %s %s", request.Method, request.URL)
				}
				for _, header := range []string{"User-Agent", "Referer", "Origin"} {
					if request.Header.Get(header) == "" {
						t.Fatalf("%s header is empty", header)
					}
				}
				if profile == "anonymous" && request.Header.Get("Cookie") != "" {
					t.Fatalf("anonymous Cookie = %q", request.Header.Get("Cookie"))
				}
				if profile != "anonymous" && request.Header.Get("Cookie") == "" {
					t.Fatal("authenticated Cookie is empty")
				}

				folder := ""
				switch request.URL.Path {
				case "/twirp/user.v1.Season/GetSeasonInfo":
					folder = "season-info"
				case "/twirp/activity.v1.Activity/GetClockInInfo":
					folder = "clock-in-info"
				case "/twirp/pointshop.v1.Pointshop/GetUserPoint":
					folder = "user-point"
				case "/twirp/pointshop.v1.Pointshop/ListProduct":
					folder = "point-products"
				case "/twirp/user.v1.User/GetCoupons":
					if got := request.Header.Get("Content-Type"); got != "application/json" {
						t.Fatalf("Content-Type = %q", got)
					}
					body, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("read body: %v", err)
					}
					var got map[string]any
					if err := json.Unmarshal(body, &got); err != nil {
						t.Fatalf("decode body: %v", err)
					}
					if got["pageNum"] != float64(1) || got["pageSize"] != float64(20) || got["notExpired"] != true || got["tabType"] != float64(1) || got["type"] != float64(0) {
						t.Fatalf("coupon body = %v", got)
					}
					if profile == "anonymous" {
						body := contracttest.Fixture(t, "manga", "read-core", "coupons", "responses", "anonymous.requires_login.json")
						return testutil.JSONResponse(http.StatusUnauthorized, string(body)), nil
					}
					fixtureBody := contracttest.Fixture(t, "manga", "read-core", "coupons", "responses", "authenticated.success.json")
					return testutil.JSONResponse(http.StatusOK, string(fixtureBody)), nil
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				if request.Body != nil {
					body, err := io.ReadAll(request.Body)
					if err != nil || len(body) != 0 {
						t.Fatalf("empty POST body = %q, %v", body, err)
					}
				}
				body := contracttest.Fixture(t, "manga", "read-core", folder, "responses", "success.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})

			options := []bpi.Option{bpi.WithHTTPClient(&http.Client{Transport: transport})}
			if profile != "anonymous" {
				options = append(options, bpi.WithCookie("SESSDATA=fixture-session"))
			}
			client, err := bpi.NewClient(options...)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			domain := client.Manga()
			ctx := context.Background()

			season, err := domain.SeasonInfo(ctx)
			if err != nil || season.SeasonID == "" || season.CurrentTime == "" {
				t.Fatalf("SeasonInfo() = %+v, %v", season, err)
			}
			clock, err := domain.ClockInInfo(ctx)
			if err != nil || len(clock.Points) != 7 || len(clock.PointInfos) == 0 {
				t.Fatalf("ClockInInfo() = %+v, %v", clock, err)
			}
			point, err := domain.UserPoint(ctx)
			if err != nil || point.Point != "0" {
				t.Fatalf("UserPoint() = %+v, %v", point, err)
			}
			products, err := domain.PointProducts(ctx)
			if err != nil || len(products) != 1 || products[0].ID != 1938 {
				t.Fatalf("PointProducts() = %+v, %v", products, err)
			}
			coupons, err := domain.Coupons(ctx, manga.NewCouponsParams())
			if profile == "anonymous" {
				var httpError *bpi.HTTPError
				if !errors.As(err, &httpError) || !bpi.RequiresLogin(err) {
					t.Fatalf("Coupons() error = %v, want HTTP 401 login error", err)
				}
				return
			}
			if err != nil || coupons.TotalRemainAmount != 0 || coupons.Items == nil {
				t.Fatalf("Coupons() = %+v, %v", coupons, err)
			}
		})
	}
}

func TestMangaNilDomainClientReturnsParameterError(t *testing.T) {
	var domain bpi.MangaClient
	_, err := domain.SeasonInfo(context.Background())
	var parameterError *bpi.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("SeasonInfo() error = %v, want ParameterError", err)
	}
}
