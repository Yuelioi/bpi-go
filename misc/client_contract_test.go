package misc_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	bpi "github.com/Yuelioi/bpi-go"
	core "github.com/Yuelioi/bpi-go/client"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/sign"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/misc"
)

func TestMiscPromotedContractsUseAllProfiles(t *testing.T) {
	t.Parallel()
	const timestamp = uint64(1_234_567_890)
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					if request.Header.Get("User-Agent") == "" {
						t.Fatal("User-Agent is empty")
					}
					switch request.URL.Path {
					case "/x/web-frontend/getbuvid":
						if request.Method != http.MethodGet || request.URL.RawQuery != "" {
							t.Fatalf("Buvid3 request = %s %s", request.Method, request.URL)
						}
						body := contracttest.Fixture(t, "misc", "buvid3", "responses", profile+".success.json")
						return testutil.JSONResponse(http.StatusOK, string(body)), nil
					case "/x/frontend/finger/spi":
						if request.Method != http.MethodGet || request.URL.RawQuery != "" {
							t.Fatalf("Buvid request = %s %s", request.Method, request.URL)
						}
						body := contracttest.Fixture(t, "misc", "buvid", "responses", profile+".success.json")
						return testutil.JSONResponse(http.StatusOK, string(body)), nil
					case "/x/share/click":
						if request.Method != http.MethodPost || request.URL.Host != "api.biliapi.net" || request.Header.Get("Referer") == "" || request.Header.Get("Origin") == "" {
							t.Fatalf("ShortLink request = %s %s headers %v", request.Method, request.URL, request.Header)
						}
						if cookie := request.Header.Get("Cookie"); cookie != "" {
							t.Fatalf("cross-host Cookie leaked: %q", cookie)
						}
						body, err := io.ReadAll(request.Body)
						if err != nil {
							t.Fatalf("read form: %v", err)
						}
						form, err := url.ParseQuery(string(body))
						if err != nil || form.Encode() != "build=6114514&buvid=qwq&oid=10001&platform=unix&share_channel=COPY&share_id=main.ugc-video-detail.0.0.pv&share_mode=4" {
							t.Fatalf("form = %v, %v", form, err)
						}
						fixture := contracttest.Fixture(t, "misc", "b23tv", "short-link", "responses", "success.json")
						return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
					case "/bapis/bilibili.api.ticket.v1.Ticket/GenWebTicket":
						csrf := ""
						if profile != "anonymous" {
							csrf = "fixture-csrf"
						}
						query := request.URL.Query()
						if request.Method != http.MethodPost || query.Get("key_id") != "ec02" || query.Get("context[ts]") != "1234567890" || query.Get("hexsign") != sign.TicketHexSign(timestamp) || query.Get("csrf") != csrf {
							t.Fatalf("BiliTicket request = %s %s", request.Method, request.URL)
						}
						fixture := contracttest.Fixture(t, "misc", "sign", "bili-ticket", "responses", "success.json")
						return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
					default:
						t.Fatalf("unexpected request %s", request.URL)
					}
					return nil, nil
				})}),
				contracttest.ProfileOption(profile),
				core.WithClock(func() time.Time { return time.Unix(int64(timestamp), 0) }),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			domain := client.Misc()
			ctx := context.Background()
			buvid3, err := domain.Buvid3(ctx)
			if err != nil || buvid3.Buvid != "BUVID3_SANITIZED" {
				t.Fatalf("Buvid3() = %+v, %v", buvid3, err)
			}
			buvid, err := domain.Buvid(ctx)
			if err != nil || buvid.Buvid3 != "BUVID3_SANITIZED" || buvid.Buvid4 != "BUVID4_SANITIZED" {
				t.Fatalf("Buvid() = %+v, %v", buvid, err)
			}
			aid, _ := ids.NewAID(10001)
			short, err := domain.ShortLink(ctx, misc.NewShortLinkParams(aid))
			if err != nil || short.Title != "sanitized-title" || short.Link != "https://b23.tv/sanitized" {
				t.Fatalf("ShortLink() = %+v, %v", short, err)
			}
			ticket, err := domain.BiliTicket(ctx)
			if err != nil || ticket.TTL != 259_200 || len(ticket.Ticket) == 0 || ticket.Navigation.Image == "" {
				t.Fatalf("BiliTicket() = %+v, %v", ticket, err)
			}
			ticketString, err := domain.BiliTicketString(ctx)
			if err != nil || ticketString != ticket.Ticket {
				t.Fatalf("BiliTicketString() = %q, %v", ticketString, err)
			}
		})
	}
}
