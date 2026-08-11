package dynamic_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/dynamic"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestDynamicPromotedReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	type endpoint struct {
		folder        string
		expectedQuery string
		authenticated bool
	}
	endpoints := map[string]endpoint{
		"/x/polymer/web-dynamic/v1/detail":              {"detail/detail", "features=htmlNewStyle%2CitemOpusStyle%2CdecorationCard&id=1099138163191840776", false},
		"/x/polymer/web-dynamic/v1/detail/forward":      {"detail/forwards", "id=1099138163191840776", false},
		"/x/polymer/web-dynamic/v1/detail/forward/item": {"detail/forward-item", "id=1110902525317349376", true},
		"/x/polymer/web-dynamic/v1/detail/pic":          {"detail/pics", "id=1099138163191840776", false},
		"/x/polymer/web-dynamic/v1/detail/reaction":     {"detail/reactions", "id=1099138163191840776", false},
		"/x/polymer/web-dynamic/v1/feed/all":            {"feed/all", "features=itemOpusStyle%2ClistOnlyfans%2CopusBigCover%2ConlyfansVote%2CdecorationCard%2ConlyfansAssetsV2%2CforwardListHidden%2CugcDelete&web_location=333.1365", true},
		"/x/polymer/web-dynamic/v1/feed/all/update":     {"feed/check-new", "update_baseline=0", true},
		"/x/dynamic/feed/dyn/banner":                    {"feed/banner", "platform=1&position=web%E5%8A%A8%E6%80%81&web_location=333.1365", false},
		"/x/polymer/web-dynamic/v1/feed/nav":            {"feed/nav", "", true},
		"/dynamic_svr/v1/dynamic_svr/w_live_users":      {"content/live-users", "size=1", true},
		"/lottery_svr/v1/lottery_svr/lottery_notice":    {"lottery-notice-read/lottery-notice", "", false},
		"/x/polymer/web-dynamic/v1/portal":              {"content/recent-up", "", true},
		"/dynamic_svr/v1/dynamic_svr/w_dyn_uplist":      {"content/up-users", "teenagers_mode=0", true},
	}

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			options := []bpi.Option{
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					spec, ok := endpoints[request.URL.Path]
					if !ok {
						t.Fatalf("unexpected request %s", request.URL)
					}
					expectedQuery := spec.expectedQuery
					if request.URL.Path == "/lottery_svr/v1/lottery_svr/lottery_notice" {
						csrf := ""
						if profile != "anonymous" {
							csrf = "fixture-csrf"
						}
						expectedQuery = "business_id=969916293954142214&business_type=1&csrf=" + csrf
					}
					if request.Method != http.MethodGet || request.URL.Query().Encode() != expectedQuery {
						t.Fatalf("request = %s %s, want query %q", request.Method, request.URL, expectedQuery)
					}
					if profile == "anonymous" && request.Header.Get("Cookie") != "" {
						t.Fatalf("anonymous Cookie = %q", request.Header.Get("Cookie"))
					}
					if profile != "anonymous" && request.Header.Get("Cookie") == "" {
						t.Fatal("authenticated Cookie is empty")
					}
					fileName := profile + ".success.json"
					if spec.authenticated && profile == "anonymous" {
						fileName = "anonymous.requires_login.json"
					}
					if request.URL.Path == "/lottery_svr/v1/lottery_svr/lottery_notice" {
						fileName = "success.json"
					}
					body := contracttest.Fixture(t, append([]string{"dynamic"}, append(splitFixturePath(spec.folder), "responses", fileName)...)...)
					return testutil.JSONResponse(http.StatusOK, string(body)), nil
				})}),
			}
			if profile != "anonymous" {
				options = append(options, bpi.WithAccount(bpi.Account{DedeUserID: "1", SESSDATA: "fixture-session", BiliJCT: "fixture-csrf", Buvid3: "fixture-buvid"}))
			}
			client, err := bpi.NewClient(options...)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			ctx := context.Background()
			domain := client.Dynamic()
			detailID, _ := ids.NewDynamicID("1099138163191840776")
			forwardID, _ := ids.NewDynamicID("1110902525317349376")
			lotteryID, _ := ids.NewDynamicID("969916293954142214")

			detail, err := domain.Detail(ctx, dynamic.NewDetailParams(detailID))
			if err != nil || detail.Item.ID != detailID.String() {
				t.Fatalf("Detail() = %+v, %v", detail, err)
			}
			reactions, err := domain.Reactions(ctx, dynamic.NewOffsetParams(detailID))
			if err != nil || reactions.Offset == "" {
				t.Fatalf("Reactions() = %+v, %v", reactions, err)
			}
			forwards, err := domain.Forwards(ctx, dynamic.NewOffsetParams(detailID))
			if err != nil || len(forwards.Items) != 1 {
				t.Fatalf("Forwards() = %+v, %v", forwards, err)
			}
			pictures, err := domain.Pictures(ctx, dynamic.NewItemParams(detailID))
			if err != nil || len(pictures) != 1 || pictures[0].Source == "" {
				t.Fatalf("Pictures() = %+v, %v", pictures, err)
			}
			banners, err := domain.FeedBanner(ctx)
			if err != nil || len(banners.Banners) == 0 {
				t.Fatalf("FeedBanner() = %+v, %v", banners, err)
			}
			lottery, err := domain.LotteryNotice(ctx, dynamic.NewLotteryNoticeParams(lotteryID))
			if err != nil || lottery.BusinessID != 969_916_293_954_142_214 {
				t.Fatalf("LotteryNotice() = %+v, %v", lottery, err)
			}

			all, allErr := domain.All(ctx, dynamic.NewAllParams())
			checkParams, _ := dynamic.NewCheckNewParams("0")
			update, updateErr := domain.CheckNew(ctx, checkParams)
			nav, navErr := domain.NavFeed(ctx, dynamic.NewNavFeedParams())
			forward, forwardErr := domain.ForwardItem(ctx, dynamic.NewItemParams(forwardID))
			liveParams, _ := dynamic.NewLiveUsersParams().WithSize(1)
			live, liveErr := domain.LiveUsers(ctx, liveParams)
			up, upErr := domain.UpUsers(ctx, dynamic.NewUpUsersParams())
			recent, recentErr := domain.RecentUp(ctx)
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"all": allErr, "update": updateErr, "nav": navErr, "forward": forwardErr, "live": liveErr, "up": upErr, "recent": recentErr} {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if allErr != nil || len(all.Items) != 1 || updateErr != nil || update.UpdateNum != 0 || navErr != nil || len(nav.Items) != 1 || nav.Items[0].Author.MID != 1 {
				t.Fatalf("feed reads = all %+v/%v update %+v/%v nav %+v/%v", all, allErr, update, updateErr, nav, navErr)
			}
			if forwardErr != nil || forward.Item.ID != forwardID.String() || liveErr != nil || live.Group != "default" || upErr != nil || recentErr != nil || recent.MyInfo == nil || recent.MyInfo.MID != 1 {
				t.Fatalf("authenticated reads = forward %+v/%v live %+v/%v up %+v/%v recent %+v/%v", forward, forwardErr, live, liveErr, up, upErr, recent, recentErr)
			}
		})
	}
}

func splitFixturePath(path string) []string {
	return strings.Split(path, "/")
}

func TestDynamicNilClientReturnsParameterError(t *testing.T) {
	var domain bpi.DynamicClient
	_, err := domain.RecentUp(context.Background())
	var parameterError *bpi.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("RecentUp() error = %v, want ParameterError", err)
	}
}
