package message_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/message"
)

func TestMessageReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder, query := "", ""
				switch request.URL.Path {
				case "/x/im/web/msgfeed/unread":
					folder, query = "unread-count", "build=0&mobi_app=web"
				case "/x/msgfeed/reply":
					folder, query = "reply-feed", "build=0&mobi_app=web&platform=web&web_location="
				case "/session_svr/v1/session_svr/single_unread":
					folder, query = "single-unread", "build=0&mobi_app=web&show_dustbin=0&show_unfollow_list=0&unread_type=0"
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				if request.Method != http.MethodGet || request.URL.Query().Encode() != query {
					t.Fatalf("request = %s %s", request.Method, request.URL)
				}
				fileName := "authenticated.success.json"
				if profile == "anonymous" {
					fileName = "anonymous.requires_login.json"
				}
				body := contracttest.Fixture(t, "message", "read", folder, "responses", fileName)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			domain, ctx := client.Message(), context.Background()
			counts, countsErr := domain.UnreadCount(ctx, message.NewUnreadCountParams())
			feed, feedErr := domain.ReplyFeed(ctx, message.NewReplyFeedParams())
			single, singleErr := domain.SingleUnread(ctx, message.NewSingleUnreadParams())
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"counts": countsErr, "feed": feedErr, "single": singleErr} {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if countsErr != nil || counts.System != 1 {
				t.Fatalf("UnreadCount() = %+v, %v", counts, countsErr)
			}
			if feedErr != nil || len(feed.Items) != 1 || feed.Items[0].User.Nickname != "sanitized user" {
				t.Fatalf("ReplyFeed() = %+v, %v", feed, feedErr)
			}
			if singleErr != nil || single.FollowUnread != 0 {
				t.Fatalf("SingleUnread() = %+v, %v", single, singleErr)
			}
		})
	}
}
