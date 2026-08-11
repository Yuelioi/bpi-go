package login_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/login"
)

func TestLoginPrivateReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	mid, _ := ids.NewMID(1000001)
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				type endpointSpec struct {
					parts []string
					query string
				}
				specs := map[string]endpointSpec{
					"/x/member/web/account":           {[]string{"account-info"}, ""},
					"/site/getCoin":                   {[]string{"coin"}, ""},
					"/x/web-interface/nav":            {[]string{"nav"}, ""},
					"/x/web-interface/nav/stat":       {[]string{"stat"}, ""},
					"/x/web-interface/coin/today/exp": {[]string{"today-coin-exp"}, ""},
					"/x/vip/web/user/info":            {[]string{"vip-info"}, ""},
					"/x/safecenter/login_notice":      {[]string{"notice", "login-notice"}, "mid=1000001"},
					"/x/member/web/login/log":         {[]string{"notice", "login-log"}, "jsonp=jsonp&web_location=333.33"},
				}
				spec, ok := specs[request.URL.Path]
				if request.Method != http.MethodGet || !ok || request.URL.Query().Encode() != spec.query {
					t.Fatalf("unexpected request %s %s", request.Method, request.URL)
				}
				fileName := profile + ".success.json"
				if profile == "anonymous" {
					fileName = "anonymous.error.json"
				}
				parts := append([]string{"login"}, spec.parts...)
				parts = append(parts, "responses", fileName)
				body := contracttest.Fixture(t, parts...)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			domain, ctx := client.Login(), context.Background()
			account, accountErr := domain.AccountInfo(ctx)
			coin, coinErr := domain.Coin(ctx)
			nav, navErr := domain.Nav(ctx)
			stats, statsErr := domain.Stat(ctx)
			exp, expErr := domain.TodayCoinExp(ctx)
			vipInfo, vipErr := domain.VIPInfo(ctx)
			notice, noticeErr := domain.Notice(ctx, login.NewNoticeParams(mid))
			log, logErr := domain.Log(ctx, login.NewLogParams())
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"account": accountErr, "coin": coinErr, "nav": navErr, "stats": statsErr, "exp": expErr, "vip": vipErr, "notice": noticeErr, "log": logErr} {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if accountErr != nil || account.MID == 0 || coinErr != nil || coin.Money < 0 || navErr != nil || !nav.IsLogin || statsErr != nil || stats.Following == 0 || expErr != nil || vipErr != nil || noticeErr != nil || notice.MID == 0 || logErr != nil || len(log.List) != 1 {
				t.Fatalf("reads failed: account=%v coin=%v nav=%v stats=%v exp=%v vip=%v notice=%v log=%v", accountErr, coinErr, navErr, statsErr, expErr, vipErr, noticeErr, logErr)
			}
			if profile == "vip" && (!vipInfo.Active() || exp != 50) {
				t.Fatalf("vip payload = %+v exp=%d", vipInfo, exp)
			}
		})
	}
}

func TestLoginDailyRewardRecordsRiskControlProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path != "/x/member/web/exp/reward" {
					t.Fatalf("unexpected request %s", request.URL)
				}
				if profile != "normal" {
					return testutil.JSONResponse(http.StatusPreconditionFailed, `{}`), nil
				}
				body := contracttest.Fixture(t, "login", "daily-reward", "responses", "normal.success.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			reward, callErr := client.Login().DailyReward(context.Background())
			if profile == "normal" {
				if callErr != nil || !reward.Email {
					t.Fatalf("DailyReward() = %+v, %v", reward, callErr)
				}
			} else if !bpi.IsRiskControl(callErr) {
				t.Fatalf("DailyReward() error = %v, want risk control", callErr)
			}
		})
	}
}

func TestLoginCaptchaAndQRSessionFlow(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				var body []byte
				switch request.URL.Path {
				case "/x/passport-login/captcha":
					if request.URL.Query().Encode() != "source=main_web" {
						t.Fatalf("captcha query = %q", request.URL.RawQuery)
					}
					body = contracttest.Fixture(t, "login", "captcha", "generate", "responses", "success.json")
				case "/x/passport-login/web/qrcode/generate":
					body = contracttest.Fixture(t, "login", "qr", "generate", "responses", "anonymous.success.json")
				case "/x/passport-login/web/qrcode/poll":
					if request.URL.Query().Get("qrcode_key") != "sanitized-qrcode-key" {
						t.Fatalf("poll query = %q", request.URL.RawQuery)
					}
					body = contracttest.Fixture(t, "login", "qr", "poll", "responses", "waiting.success.json")
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			domain, ctx := client.Login(), context.Background()
			captcha, err := domain.GenerateCaptcha(ctx)
			if err != nil || captcha.Token == "" || captcha.GT == "" {
				t.Fatalf("GenerateCaptcha() = %+v, %v", captcha, err)
			}
			generated, err := domain.GenerateQR(ctx)
			if err != nil || generated.Key == "" {
				t.Fatalf("GenerateQR() = %+v, %v", generated, err)
			}
			pollParams, _ := login.NewQRPollParams(generated.Key)
			status, err := domain.PollQR(ctx, pollParams)
			if err != nil || status.Code != login.QRWaiting {
				t.Fatalf("PollQR() = %+v, %v", status, err)
			}
			flow, err := domain.QRFlow(ctx)
			if err != nil || flow.Generate.Key == "" || flow.Poll.Code != login.QRWaiting {
				t.Fatalf("QRFlow() = %+v, %v", flow, err)
			}
		})
	}
}
