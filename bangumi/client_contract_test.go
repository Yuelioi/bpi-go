package bangumi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/bangumi"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestBangumiTimelineDecodesResultAliasForAllPromotedProfiles(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "bangumi", "timeline", "responses", profile+".success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.Path != "/pgc/web/timeline" || request.URL.Query().Encode() != "after=7&before=3&types=1" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := bangumi.NewTimelineParams(bangumi.TimelineAnime, 3, 7)
			days, err := client.Bangumi().Timeline(context.Background(), params)
			if err != nil {
				t.Fatalf("Bangumi().Timeline() error = %v", err)
			}
			if len(days) != 11 || days[3].IsToday != 1 || len(days[5].Episodes) == 0 || days[5].Episodes[0].Title == "" {
				t.Fatalf("Timeline() returned unexpected promoted payload: %+v", days)
			}
		})
	}
}
