package electric_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/electric"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestElectricReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	up, _ := ids.NewMID(1265680561)
	publicUP, _ := ids.NewMID(53456)
	bvid, _ := ids.NewBVID("BV1Dh411S7sS")
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				type endpointSpec struct {
					batch, folder, query string
					private              bool
				}
				endpoints := map[string]endpointSpec{
					"/x/ugcpay-rank/elec/month/up":                {"public-read", "month-up-list", "up_mid=53456", false},
					"/x/web-interface/elec/show":                  {"public-read", "video-show", "bvid=BV1Dh411S7sS&mid=53456", false},
					"/bk/brokerage/listForCustomerRechargeRecord": {"private-read", "recharge-list", "currentPage=1&customerId=10026&pageSize=10", true},
					"/x/h5/elec/rank/recent":                      {"private-read", "rank-recent", "pn=1&ps=10", true},
					"/xlive/revenue/v1/guard/getChargeRecord":     {"private-read", "charge-record", "page=1&type=1", true},
					"/x/upower/item/detail":                       {"public-read", "upower-item-detail", "up_mid=1265680561", false},
					"/x/upower/charge/follow/info":                {"private-read", "charge-follow-info", "up_mid=1265680561", true},
					"/x/upower/up/member/rank/v2":                 {"public-read", "upower-member-rank", "pn=1&ps=10&up_mid=1265680561", false},
					"/x/web/elec/remark/list":                     {"private-read", "remark-list", "pn=1&ps=10", true},
					"/x/web/elec/remark/detail":                   {"private-read", "remark-detail", "id=1", true},
				}
				spec, ok := endpoints[request.URL.Path]
				if request.Method != http.MethodGet || !ok || request.URL.Query().Encode() != spec.query {
					t.Fatalf("unexpected request %s %s", request.Method, request.URL)
				}
				fileName := "success.json"
				if spec.folder == "upower-member-rank" {
					fileName = "authenticated.success.json"
					if profile == "anonymous" {
						fileName = "anonymous.success.json"
					}
				} else if spec.private {
					fileName = "authenticated.success.json"
					if profile == "anonymous" {
						fileName = "anonymous.requires_login.json"
					}
					if spec.folder == "rank-recent" && profile != "anonymous" {
						fileName = profile + ".success.json"
					}
				}
				body := contracttest.Fixture(t, "electric", spec.batch, spec.folder, "responses", fileName)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			domain, ctx := client.Electric(), context.Background()
			recharge, _ := electric.NewRechargeListParams(1, 10)
			page, _ := electric.NewPaginationParams(1, 10)
			record, _ := electric.NewChargeRecordParams(1, 1)
			memberRank, _ := electric.NewMemberRankParams(up, 1, 10)
			remarks, _ := electric.NewRemarkListParams(1, 10)
			detail, _ := electric.NewRemarkDetailParams(1)
			month, monthErr := domain.MonthUpList(ctx, electric.NewMonthUpListParams(publicUP))
			video, videoErr := domain.VideoShow(ctx, electric.NewVideoShowParams(publicUP).WithBVID(bvid))
			recharges, rechargeErr := domain.RechargeList(ctx, recharge)
			recent, recentErr := domain.RankRecent(ctx, page)
			records, recordErr := domain.ChargeRecord(ctx, record)
			item, itemErr := domain.ItemDetail(ctx, electric.NewUpMIDParams(up))
			follow, followErr := domain.ChargeFollowInfo(ctx, electric.NewUpMIDParams(up))
			members, membersErr := domain.MemberRank(ctx, memberRank)
			remarkList, remarksErr := domain.RemarkList(ctx, remarks)
			remark, detailErr := domain.RemarkDetail(ctx, detail)
			if monthErr != nil || len(month.List) != 1 || videoErr != nil || !video.ShowInfo.Show || itemErr != nil || item.Rank.Total == 0 || membersErr != nil || len(members.RankInfo) != 1 {
				t.Fatalf("public reads failed: month=%v video=%v item=%v members=%v", monthErr, videoErr, itemErr, membersErr)
			}
			if profile == "anonymous" {
				for name, callErr := range map[string]error{"recharge": rechargeErr, "recent": recentErr, "record": recordErr, "follow": followErr, "remarks": remarksErr, "detail": detailErr} {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if rechargeErr != nil || len(recharges.Result) != 1 || recentErr != nil || recent.Pager.Current != 1 || recordErr != nil || records.TotalCount != 0 || followErr != nil || follow.UPCard.MID != 1265680561 || remarksErr != nil || len(remarkList.List) != 1 || detailErr != nil || remark.ID != 1 {
				t.Fatalf("private reads failed: recharge=%v recent=%v record=%v follow=%v remarks=%v detail=%v", rechargeErr, recentErr, recordErr, followErr, remarksErr, detailErr)
			}
		})
	}
}
