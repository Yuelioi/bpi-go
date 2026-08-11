package creativecenter

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	creativeSeasonListEndpoint      = "https://member.bilibili.com/x2/creative/web/seasons"
	creativeSeasonInfoEndpoint      = "https://member.bilibili.com/x2/creative/web/season"
	creativeSeasonByAIDEndpoint     = "https://member.bilibili.com/x2/creative/web/season/aid"
	creativeSeasonSectionEndpoint   = "https://member.bilibili.com/x2/creative/web/season/section"
	creativeArchivesListEndpoint    = "https://member.bilibili.com/x2/creative/web/archives/sp"
	creativeArchiveVideosEndpoint   = "https://member.bilibili.com/x/web/archive/videos"
	creativeUpStatEndpoint          = "https://member.bilibili.com/x/web/index/stat"
	creativeArchiveCompareEndpoint  = "https://member.bilibili.com/x/web/data/archive_diagnose/compare"
	creativeArticleStatEndpoint     = "https://member.bilibili.com/x/web/data/article"
	creativeVideoTrendEndpoint      = "https://member.bilibili.com/x/web/data/pandect"
	creativeArticleTrendEndpoint    = "https://member.bilibili.com/x/web/data/article/thirty"
	creativePlaySourceEndpoint      = "https://member.bilibili.com/x/web/data/playsource"
	creativeViewerDataEndpoint      = "https://member.bilibili.com/x/web/data/base"
	creativeElectromagneticEndpoint = "https://api.bilibili.com/studio/up-rating/v3/rating/info"
)

// Client provides creator-center operations.
type Client struct{ client *core.Client }

// NewClient binds the creativecenter module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (c Client) SeasonList(ctx context.Context, params SeasonListParams) (SeasonList, error) {
	return sendCreativePayload[SeasonList](ctx, c.client, creativeSeasonListEndpoint, "creativecenter.season.list", params.EncodeQuery)
}

func (c Client) SeasonInfo(ctx context.Context, params SeasonInfoParams) (SeasonInfo, error) {
	return sendCreativePayload[SeasonInfo](ctx, c.client, creativeSeasonInfoEndpoint, "creativecenter.season.info", params.EncodeQuery)
}

func (c Client) SeasonByAID(ctx context.Context, params SeasonByAIDParams) (Season, error) {
	return sendCreativePayload[Season](ctx, c.client, creativeSeasonByAIDEndpoint, "creativecenter.season.aid", params.EncodeQuery)
}

func (c Client) SeasonSection(ctx context.Context, params SectionParams) (SectionEpisodes, error) {
	return sendCreativePayload[SectionEpisodes](ctx, c.client, creativeSeasonSectionEndpoint, "creativecenter.season.section", params.EncodeQuery)
}

func (c Client) ArchivesList(ctx context.Context, params ArchivesListParams) (ArchivesList, error) {
	return sendCreativePayload[ArchivesList](ctx, c.client, creativeArchivesListEndpoint, "creativecenter.videos.archives_list", params.EncodeQuery)
}

func (c Client) ArchiveVideos(ctx context.Context, params ArchiveVideosParams) (ArchiveVideos, error) {
	return sendCreativePayload[ArchiveVideos](ctx, c.client, creativeArchiveVideosEndpoint, "creativecenter.videos.archive_videos", params.EncodeQuery)
}

func (c Client) UpStat(ctx context.Context) (UpStat, error) {
	return sendCreativePayload[UpStat](ctx, c.client, creativeUpStatEndpoint, "creativecenter.statistics.up_stat", noCreativeQuery)
}

func (c Client) ArchiveCompare(ctx context.Context, params ArchiveCompareParams) (ArchiveCompare, error) {
	return sendCreativePayload[ArchiveCompare](ctx, c.client, creativeArchiveCompareEndpoint, "creativecenter.statistics.archive_compare", params.EncodeQuery)
}

func (c Client) ArticleStat(ctx context.Context) (ArticleStat, error) {
	return sendCreativePayload[ArticleStat](ctx, c.client, creativeArticleStatEndpoint, "creativecenter.statistics.article_stat", noCreativeQuery)
}

func (c Client) VideoTrend(ctx context.Context, params VideoTrendParams) ([]Trend, error) {
	return sendCreativePayload[[]Trend](ctx, c.client, creativeVideoTrendEndpoint, "creativecenter.statistics.video_trend", params.EncodeQuery)
}

func (c Client) ArticleTrend(ctx context.Context, params ArticleTrendParams) ([]Trend, error) {
	return sendCreativeOptional[[]Trend](ctx, c.client, creativeArticleTrendEndpoint, "creativecenter.statistics.article_trend", params.EncodeQuery)
}

func (c Client) PlaySource(ctx context.Context) (*PlaySource, error) {
	return sendCreativeOptionalPointer[PlaySource](ctx, c.client, creativePlaySourceEndpoint, "creativecenter.statistics.play_source", noCreativeQuery)
}

func (c Client) ViewerData(ctx context.Context) (ViewerData, error) {
	return sendCreativePayload[ViewerData](ctx, c.client, creativeViewerDataEndpoint, "creativecenter.statistics.viewer_data", noCreativeQuery)
}

func (c Client) ElectromagneticInfo(ctx context.Context) (ElectromagneticInfo, error) {
	return sendCreativePayload[ElectromagneticInfo](ctx, c.client, creativeElectromagneticEndpoint, "creativecenter.railgun.electromagnetic_info", noCreativeQuery)
}

func sendCreativePayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "creative center client is not initialized"}
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

func sendCreativeOptional[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	payload, err := sendCreativeOptionalPointer[T](ctx, client, endpoint, operation, encode)
	if err != nil || payload == nil {
		return zero, err
	}
	return *payload, nil
}

func sendCreativeOptionalPointer[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (*T, error) {
	if client == nil {
		return nil, &core.ParameterError{Field: "client", Message: "creative center client is not initialized"}
	}
	query, err := encode()
	if err != nil {
		return nil, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return nil, err
	}
	return core.SendOptionalPayload[T](ctx, client, request, operation)
}

func noCreativeQuery() (url.Values, error) { return nil, nil }
