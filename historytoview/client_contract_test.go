package historytoview_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/historytoview"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestHistoryToViewReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder, query := "", ""
				switch request.URL.Path {
				case "/x/web-interface/history/cursor":
					folder, query = "history-list", "ps=5"
				case "/x/v2/history/shadow":
					folder = "history-shadow"
				case "/x/v2/history/toview":
					folder = "toview-list"
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
				body := contracttest.Fixture(t, "historytoview", "read", folder, "responses", fileName)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := historytoview.NewListParams().WithPageSize(5)
			domain, ctx := client.HistoryToView(), context.Background()
			list, listErr := domain.HistoryList(ctx, params)
			shadow, shadowErr := domain.HistoryShadow(ctx)
			toView, toViewErr := domain.ToViewList(ctx)
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"list": listErr, "shadow": shadowErr, "toview": toViewErr} {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if listErr != nil || list.Cursor.PageSize != 5 || len(list.Items) != 1 || list.Items[0].Total == nil || *list.Items[0].Total != -1 {
				t.Fatalf("HistoryList() = %+v, %v", list, listErr)
			}
			if shadowErr != nil || shadow {
				t.Fatalf("HistoryShadow() = %v, %v", shadow, shadowErr)
			}
			if toViewErr != nil || toView.Count != 1 || len(toView.Items) != 1 || toView.Items[0].Stat.VT != -1 {
				t.Fatalf("ToViewList() = %+v, %v", toView, toViewErr)
			}
		})
	}
}
