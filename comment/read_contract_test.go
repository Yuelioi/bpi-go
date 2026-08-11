package comment_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/comment"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestCommentRemainingPromotedReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder := ""
				expectedQuery := ""
				switch request.URL.Path {
				case "/x/v2/reply":
					folder = "list"
					expectedQuery = "nohot=0&oid=23199&pn=1&ps=5&sort=0&type=1"
				case "/x/v2/reply/reply":
					folder = "replies"
					expectedQuery = "oid=23199&pn=1&ps=5&root=2554491176&type=1"
				case "/x/v2/reply/count":
					folder = "count"
					expectedQuery = "oid=23199&type=1"
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				if request.Method != http.MethodGet || request.URL.Host != "api.bilibili.com" || request.URL.Query().Encode() != expectedQuery {
					t.Fatalf("request = %s %s", request.Method, request.URL)
				}
				body := contracttest.Fixture(t, "comment", "read", folder, "responses", profile+".success.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			target, _ := comment.NewTarget(1, 23199)
			listParams, _ := comment.NewListParams(target).WithPage(1)
			listParams, _ = listParams.WithPageSize(5)
			listParams, _ = listParams.WithSort(comment.SortByTime)
			listParams = listParams.WithoutHot(false)
			ctx := context.Background()
			list, err := client.Comment().List(ctx, listParams)
			if err != nil || list.Page == nil || list.Page.AllCount == nil {
				t.Fatalf("List() = %+v, %v", list, err)
			}

			repliesParams, _ := comment.NewRepliesParams(target, 2_554_491_176)
			repliesParams, _ = repliesParams.WithPage(1)
			repliesParams, _ = repliesParams.WithPageSize(5)
			replies, err := client.Comment().Replies(ctx, repliesParams)
			if err != nil || replies.Page == nil || len(replies.Replies) != 1 || replies.Replies[0].Content.Message == "" {
				t.Fatalf("Replies() = %+v, %v", replies, err)
			}

			count, err := client.Comment().Count(ctx, comment.NewCountParams(target))
			if err != nil || count.Count != 10 {
				t.Fatalf("Count() = %+v, %v", count, err)
			}
		})
	}
}
