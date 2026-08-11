package vip_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/vip"
)

func TestVIPCenterUsesAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(
				bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
					if request.Method != http.MethodGet || request.URL.Path != "/x/vip/web/vip_center/combine" || request.URL.Query().Encode() != "build=0" {
						t.Fatalf("request = %s %s", request.Method, request.URL)
					}
					body := contracttest.Fixture(t, "vip", "read", "center-info", "responses", profile+".success.json")
					return testutil.JSONResponse(http.StatusOK, string(body)), nil
				})}),
				contracttest.ProfileOption(profile),
			)
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			center, err := client.VIP().Center(context.Background(), vip.NewCenterParams())
			if err != nil || center.User.PanelTitle == "" {
				t.Fatalf("Center() = %+v, %v", center, err)
			}
		})
	}
}
