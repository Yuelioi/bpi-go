package electric

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	electricMonthUpListEndpoint      = "https://api.bilibili.com/x/ugcpay-rank/elec/month/up"
	electricVideoShowEndpoint        = "https://api.bilibili.com/x/web-interface/elec/show"
	electricRechargeListEndpoint     = "https://pay.bilibili.com/bk/brokerage/listForCustomerRechargeRecord"
	electricRankRecentEndpoint       = "https://member.bilibili.com/x/h5/elec/rank/recent"
	electricChargeRecordEndpoint     = "https://api.live.bilibili.com/xlive/revenue/v1/guard/getChargeRecord"
	electricItemDetailEndpoint       = "https://api.bilibili.com/x/upower/item/detail"
	electricChargeFollowInfoEndpoint = "https://api.bilibili.com/x/upower/charge/follow/info"
	electricMemberRankEndpoint       = "https://api.bilibili.com/x/upower/up/member/rank/v2"
	electricRemarkListEndpoint       = "https://member.bilibili.com/x/web/elec/remark/list"
	electricRemarkDetailEndpoint     = "https://member.bilibili.com/x/web/elec/remark/detail"
)

// Client provides charging-support operations.
type Client struct{ client *core.Client }

// NewClient binds the electric module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (e Client) MonthUpList(ctx context.Context, params MonthUpListParams) (MonthUpList, error) {
	return sendElectricPayload[MonthUpList](ctx, e.client, electricMonthUpListEndpoint, "electric.month_up_list", params.EncodeQuery)
}

func (e Client) VideoShow(ctx context.Context, params VideoShowParams) (VideoShow, error) {
	return sendElectricPayload[VideoShow](ctx, e.client, electricVideoShowEndpoint, "electric.video_show", params.EncodeQuery)
}

func (e Client) RechargeList(ctx context.Context, params RechargeListParams) (RechargeList, error) {
	return sendElectricPayload[RechargeList](ctx, e.client, electricRechargeListEndpoint, "electric.recharge_list", params.EncodeQuery)
}

func (e Client) RankRecent(ctx context.Context, params PaginationParams) (RecentRank, error) {
	return sendElectricPayload[RecentRank](ctx, e.client, electricRankRecentEndpoint, "electric.rank_recent", params.EncodeQuery)
}

func (e Client) ChargeRecord(ctx context.Context, params ChargeRecordParams) (ChargeRecord, error) {
	return sendElectricPayload[ChargeRecord](ctx, e.client, electricChargeRecordEndpoint, "electric.charge_record", params.EncodeQuery)
}

func (e Client) ItemDetail(ctx context.Context, params UpMIDParams) (ItemDetail, error) {
	return sendElectricPayload[ItemDetail](ctx, e.client, electricItemDetailEndpoint, "electric.upower_item_detail", params.EncodeQuery)
}

func (e Client) ChargeFollowInfo(ctx context.Context, params UpMIDParams) (FollowInfo, error) {
	return sendElectricPayload[FollowInfo](ctx, e.client, electricChargeFollowInfoEndpoint, "electric.charge_follow_info", params.EncodeQuery)
}

func (e Client) MemberRank(ctx context.Context, params MemberRankParams) (MemberRank, error) {
	return sendElectricPayload[MemberRank](ctx, e.client, electricMemberRankEndpoint, "electric.upower_member_rank", params.EncodeQuery)
}

func (e Client) RemarkList(ctx context.Context, params RemarkListParams) (RemarkList, error) {
	return sendElectricPayload[RemarkList](ctx, e.client, electricRemarkListEndpoint, "electric.remark_list", params.EncodeQuery)
}

func (e Client) RemarkDetail(ctx context.Context, params RemarkDetailParams) (RemarkDetail, error) {
	return sendElectricPayload[RemarkDetail](ctx, e.client, electricRemarkDetailEndpoint, "electric.remark_detail", params.EncodeQuery)
}

func sendElectricPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "electric client is not initialized"}
	}
	query, err := encode()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
