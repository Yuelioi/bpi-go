package dynamic

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	dynamicAllEndpoint         = "https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/all"
	dynamicCheckNewEndpoint    = "https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/all/update"
	dynamicNavFeedEndpoint     = "https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/nav"
	dynamicFeedBannerEndpoint  = "https://api.bilibili.com/x/dynamic/feed/dyn/banner"
	dynamicDetailEndpoint      = "https://api.bilibili.com/x/polymer/web-dynamic/v1/detail"
	dynamicReactionsEndpoint   = "https://api.bilibili.com/x/polymer/web-dynamic/v1/detail/reaction"
	dynamicLotteryEndpoint     = "https://api.vc.bilibili.com/lottery_svr/v1/lottery_svr/lottery_notice"
	dynamicForwardsEndpoint    = "https://api.bilibili.com/x/polymer/web-dynamic/v1/detail/forward"
	dynamicPicturesEndpoint    = "https://api.bilibili.com/x/polymer/web-dynamic/v1/detail/pic"
	dynamicForwardItemEndpoint = "https://api.bilibili.com/x/polymer/web-dynamic/v1/detail/forward/item"
	dynamicLiveUsersEndpoint   = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/w_live_users"
	dynamicUpUsersEndpoint     = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/w_dyn_uplist"
	dynamicRecentUpEndpoint    = "https://api.bilibili.com/x/polymer/web-dynamic/v1/portal"
)

// Client provides promoted dynamic-feed read operations.
type Client struct{ client *core.Client }

// NewClient binds the dynamic module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (d Client) All(ctx context.Context, params AllParams) (Feed, error) {
	return sendDynamicPayload[Feed](ctx, d.client, dynamicAllEndpoint, "dynamic.feed_all", params.EncodeQuery)
}

func (d Client) CheckNew(ctx context.Context, params CheckNewParams) (Update, error) {
	return sendDynamicPayload[Update](ctx, d.client, dynamicCheckNewEndpoint, "dynamic.feed_all_update", params.EncodeQuery)
}

func (d Client) NavFeed(ctx context.Context, params NavFeedParams) (NavFeed, error) {
	return sendDynamicPayload[NavFeed](ctx, d.client, dynamicNavFeedEndpoint, "dynamic.feed_nav", params.EncodeQuery)
}

func (d Client) FeedBanner(ctx context.Context) (BannerFeed, error) {
	return sendDynamicPayload[BannerFeed](ctx, d.client, dynamicFeedBannerEndpoint, "dynamic.feed_banner", func() (url.Values, error) {
		return FeedBannerQuery(), nil
	})
}

func (d Client) Detail(ctx context.Context, params DetailParams) (Detail, error) {
	return sendDynamicPayload[Detail](ctx, d.client, dynamicDetailEndpoint, "dynamic.detail", params.EncodeQuery)
}

func (d Client) Reactions(ctx context.Context, params OffsetParams) (Reactions, error) {
	return sendDynamicPayload[Reactions](ctx, d.client, dynamicReactionsEndpoint, "dynamic.detail_reaction", params.EncodeQuery)
}

func (d Client) LotteryNotice(ctx context.Context, params LotteryNoticeParams) (LotteryNotice, error) {
	return sendDynamicPayload[LotteryNotice](ctx, d.client, dynamicLotteryEndpoint, "dynamic.lottery_notice", func() (url.Values, error) {
		csrf := ""
		if d.client != nil {
			csrf, _ = d.client.CSRF()
		}
		return params.EncodeQuery(csrf)
	})
}

func (d Client) Forwards(ctx context.Context, params OffsetParams) (Forwards, error) {
	return sendDynamicPayload[Forwards](ctx, d.client, dynamicForwardsEndpoint, "dynamic.detail_forward", params.EncodeQuery)
}

func (d Client) Pictures(ctx context.Context, params ItemParams) ([]Picture, error) {
	return sendDynamicPayload[[]Picture](ctx, d.client, dynamicPicturesEndpoint, "dynamic.detail_pic", params.EncodeQuery)
}

func (d Client) ForwardItem(ctx context.Context, params ItemParams) (ForwardInfo, error) {
	return sendDynamicPayload[ForwardInfo](ctx, d.client, dynamicForwardItemEndpoint, "dynamic.detail_forward_item", params.EncodeQuery)
}

func (d Client) LiveUsers(ctx context.Context, params LiveUsersParams) (LiveUsers, error) {
	return sendDynamicPayload[LiveUsers](ctx, d.client, dynamicLiveUsersEndpoint, "dynamic.live_users", params.EncodeQuery)
}

func (d Client) UpUsers(ctx context.Context, params UpUsersParams) (UpUsers, error) {
	return sendDynamicPayload[UpUsers](ctx, d.client, dynamicUpUsersEndpoint, "dynamic.up_users", params.EncodeQuery)
}

func (d Client) RecentUp(ctx context.Context) (RecentUp, error) {
	return sendDynamicPayload[RecentUp](ctx, d.client, dynamicRecentUpEndpoint, "dynamic.recent_up", func() (url.Values, error) { return nil, nil })
}

func sendDynamicPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "dynamic client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
