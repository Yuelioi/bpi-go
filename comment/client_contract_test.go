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

func TestCommentHotReturnsNilForPromotedOptionalPayloads(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "comment", "read", "hot", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.Path != "/x/v2/reply/hot" || request.URL.Query().Encode() != "oid=23199&pn=1&ps=5&root=2554491176&type=1" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			target, _ := comment.NewTarget(1, 23199)
			params, _ := comment.NewHotParams(target, 2_554_491_176)
			params, _ = params.WithPage(1)
			params, _ = params.WithPageSize(5)
			payload, err := client.Comment().Hot(context.Background(), params)
			if err != nil {
				t.Fatalf("Comment().Hot() error = %v", err)
			}
			if payload != nil {
				t.Fatalf("Comment().Hot() = %+v, want nil optional payload", payload)
			}
		})
	}
}
