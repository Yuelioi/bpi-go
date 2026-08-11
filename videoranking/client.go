package videoranking

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	videoRankingPopularEndpoint       = "https://api.bilibili.com/x/web-interface/popular"
	videoRankingSeriesListEndpoint    = "https://api.bilibili.com/x/web-interface/popular/series/list"
	videoRankingSeriesOneEndpoint     = "https://api.bilibili.com/x/web-interface/popular/series/one"
	videoRankingPreciousEndpoint      = "https://api.bilibili.com/x/web-interface/popular/precious"
	videoRankingListEndpoint          = "https://api.bilibili.com/x/web-interface/ranking/v2"
	videoRankingRegionDynamicEndpoint = "https://api.bilibili.com/x/web-interface/dynamic/region"
	videoRankingRegionTagEndpoint     = "https://api.bilibili.com/x/web-interface/dynamic/tag"
	videoRankingRegionNewListEndpoint = "https://api.bilibili.com/x/web-interface/newlist"
	videoRankingRegionNewRankEndpoint = "https://api.bilibili.com/x/web-interface/newlist_rank"
)

// Client provides video-ranking operations.
type Client struct{ client *core.Client }

// NewClient binds the videoranking module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (v Client) PopularList(ctx context.Context, params PopularListParams) (PopularList, error) {
	return sendRankingPayload[PopularList](ctx, v.client, videoRankingPopularEndpoint, "video_ranking.popular_list", params.EncodeQuery)
}

func (v Client) PopularSeriesList(ctx context.Context) (PopularSeriesList, error) {
	return sendRankingNoParams[PopularSeriesList](ctx, v.client, videoRankingSeriesListEndpoint, "video_ranking.popular_series_list")
}

func (v Client) PopularSeries(ctx context.Context, params PopularSeriesParams) (PopularSeries, error) {
	return sendRankingWBI[PopularSeries](ctx, v.client, videoRankingSeriesOneEndpoint, "video_ranking.popular_series_one", params.EncodeQuery)
}

func (v Client) Precious(ctx context.Context) (PreciousVideos, error) {
	return sendRankingNoParams[PreciousVideos](ctx, v.client, videoRankingPreciousEndpoint, "video_ranking.popular_precious")
}

func (v Client) RankingList(ctx context.Context, params RankingListParams) (RankingList, error) {
	return sendRankingPayload[RankingList](ctx, v.client, videoRankingListEndpoint, "video_ranking.ranking_list", params.EncodeQuery)
}

func (v Client) RegionDynamic(ctx context.Context, params RegionDynamicParams) (RegionArchives, error) {
	return sendRankingPayload[RegionArchives](ctx, v.client, videoRankingRegionDynamicEndpoint, "video_ranking.region_dynamic", params.EncodeQuery)
}

func (v Client) RegionTagDynamic(ctx context.Context, params RegionTagDynamicParams) (RegionArchives, error) {
	return sendRankingPayload[RegionArchives](ctx, v.client, videoRankingRegionTagEndpoint, "video_ranking.region_tag_dynamic", params.EncodeQuery)
}

func (v Client) RegionNewList(ctx context.Context, params RegionNewListParams) (RegionArchives, error) {
	return sendRankingPayload[RegionArchives](ctx, v.client, videoRankingRegionNewListEndpoint, "video_ranking.region_newlist", params.EncodeQuery)
}

func (v Client) RegionNewListRank(ctx context.Context, params RegionNewListRankParams) (NewListRank, error) {
	return sendRankingPayload[NewListRank](ctx, v.client, videoRankingRegionNewRankEndpoint, "video_ranking.region_newlist_rank", params.EncodeQuery)
}

func sendRankingNoParams[T any](ctx context.Context, client *core.Client, endpoint, operation string) (T, error) {
	return sendRankingPayload[T](ctx, client, endpoint, operation, func() (url.Values, error) { return nil, nil })
}

func sendRankingPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "video-ranking client is not initialized"}
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

func sendRankingWBI[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "video-ranking client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	signed, err := client.WBIValues(ctx, query)
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, signed)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
