package video

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	videoViewEndpoint                    = "https://api.bilibili.com/x/web-interface/view"
	videoDetailEndpoint                  = "https://api.bilibili.com/x/web-interface/view/detail"
	videoPageListEndpoint                = "https://api.bilibili.com/x/player/pagelist"
	videoDescEndpoint                    = "https://api.bilibili.com/x/web-interface/archive/desc"
	videoPlayURLEndpoint                 = "https://api.bilibili.com/x/player/wbi/playurl"
	videoOnlineTotalEndpoint             = "https://api.bilibili.com/x/player/online/total"
	videoRelatedEndpoint                 = "https://api.bilibili.com/x/web-interface/archive/related"
	videoTagsEndpoint                    = "https://api.bilibili.com/x/web-interface/view/detail/tag"
	videoInteractiveInfoEndpoint         = "https://api.bilibili.com/x/stein/edgeinfo_v2"
	videoPlayerInfoV2Endpoint            = "https://api.bilibili.com/x/player/wbi/v2"
	videoHomepageRecommendationsEndpoint = "https://api.bilibili.com/x/web-interface/wbi/index/top/feed/rcmd"
	videoAISummaryEndpoint               = "https://api.bilibili.com/x/web-interface/view/conclusion/get"
	videoSeasonsArchivesEndpoint         = "https://api.bilibili.com/x/polymer/web-space/seasons_archives_list"
	videoHomeSeasonsSeriesEndpoint       = "https://api.bilibili.com/x/polymer/web-space/home/seasons_series"
	videoSeasonsSeriesEndpoint           = "https://api.bilibili.com/x/polymer/web-space/seasons_series_list"
	videoSeriesInfoEndpoint              = "https://api.bilibili.com/x/series/series"
	videoSeriesArchivesEndpoint          = "https://api.bilibili.com/x/series/archives"
)

// Client provides video-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the video module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// View gets Web video metadata by AID or BVID.
func (v Client) View(ctx context.Context, params ViewParams) (View, error) {
	return sendVideoPayload[View](ctx, v.client, videoViewEndpoint, "video.view", params.EncodeQuery)
}

// Detail gets Web video metadata plus tags and related videos.
func (v Client) Detail(ctx context.Context, params DetailParams) (Detail, error) {
	return sendVideoPayload[Detail](ctx, v.client, videoDetailEndpoint, "video.detail", params.EncodeQuery)
}

// PageList gets a video's content/part list.
func (v Client) PageList(ctx context.Context, params PageListParams) ([]Page, error) {
	return sendVideoPayload[[]Page](ctx, v.client, videoPageListEndpoint, "video.pagelist", params.EncodeQuery)
}

// Desc gets a video's plain-text description.
func (v Client) Desc(ctx context.Context, params DescParams) (string, error) {
	return sendVideoPayload[string](ctx, v.client, videoDescEndpoint, "video.desc", params.EncodeQuery)
}

// PlayURL gets WBI-signed playback metadata by video and content ID.
func (v Client) PlayURL(ctx context.Context, params PlayURLParams) (PlayURL, error) {
	var zero PlayURL
	if v.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "video client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	signed, err := v.client.WBIValues(ctx, query)
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, videoPlayURLEndpoint, signed)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[PlayURL](ctx, v.client, request, "video.play_url")
}

// OnlineTotal gets the current audience counts for one video part.
func (v Client) OnlineTotal(ctx context.Context, params OnlineTotalParams) (OnlineTotal, error) {
	return sendVideoPayload[OnlineTotal](ctx, v.client, videoOnlineTotalEndpoint, "video.online_total", params.EncodeQuery)
}

// RelatedVideos gets videos related to the selected video.
func (v Client) RelatedVideos(ctx context.Context, params RelatedParams) ([]RelatedVideo, error) {
	return sendVideoPayload[[]RelatedVideo](ctx, v.client, videoRelatedEndpoint, "video.related_videos", params.EncodeQuery)
}

// Tags gets tags associated with a video and optional content part.
func (v Client) Tags(ctx context.Context, params TagsParams) ([]VideoTag, error) {
	return sendVideoPayload[[]VideoTag](ctx, v.client, videoTagsEndpoint, "video.tags", params.EncodeQuery)
}

// InteractiveVideoInfo gets metadata for one interactive-video graph edge.
func (v Client) InteractiveVideoInfo(ctx context.Context, params InteractiveInfoParams) (InteractiveInfo, error) {
	return sendVideoPayload[InteractiveInfo](ctx, v.client, videoInteractiveInfoEndpoint, "video.interactive_video_info", params.EncodeQuery)
}

// PlayerInfoV2 gets WBI-signed Web-player metadata for one video part.
func (v Client) PlayerInfoV2(ctx context.Context, params PlayerInfoParams) (PlayerInfo, error) {
	return sendVideoWBIPayload[PlayerInfo](ctx, v.client, videoPlayerInfoV2Endpoint, "video.player_info_v2", params.EncodeQuery)
}

// HomepageRecommendations gets the WBI-signed public homepage feed.
func (v Client) HomepageRecommendations(ctx context.Context, params HomepageRecommendationsParams) (HomepageRecommendations, error) {
	return sendVideoWBIPayload[HomepageRecommendations](ctx, v.client, videoHomepageRecommendationsEndpoint, "video.homepage_recommendations", params.EncodeQuery)
}

// AISummary gets the WBI-signed AI summary for one video part. The endpoint
// normally requires an authenticated session even though key discovery can
// operate anonymously.
func (v Client) AISummary(ctx context.Context, params AISummaryParams) (AISummary, error) {
	return sendVideoWBIPayload[AISummary](ctx, v.client, videoAISummaryEndpoint, "video.ai_summary", params.EncodeQuery)
}

// SeasonsArchives gets one WBI-signed page of a user's video collection.
func (v Client) SeasonsArchives(ctx context.Context, params SeasonsArchivesParams) (SeasonsArchives, error) {
	return sendVideoWBIPayload[SeasonsArchives](ctx, v.client, videoSeasonsArchivesEndpoint, "video.collection.seasons_archives_list", params.EncodeQuery)
}

// HomeSeasonsSeries gets a WBI-signed overview of collections and series on a
// user's public space.
func (v Client) HomeSeasonsSeries(ctx context.Context, params HomeSeasonsSeriesParams) (SeasonsSeries, error) {
	return sendVideoWBIPayload[SeasonsSeries](ctx, v.client, videoHomeSeasonsSeriesEndpoint, "video.collection.home_seasons_series", params.EncodeQuery)
}

// SeasonsSeries gets a WBI-signed paginated collection/series listing.
func (v Client) SeasonsSeries(ctx context.Context, params SeasonsSeriesParams) (SeasonsSeries, error) {
	return sendVideoWBIPayload[SeasonsSeries](ctx, v.client, videoSeasonsSeriesEndpoint, "video.collection.seasons_series_list", params.EncodeQuery)
}

// SeriesInfo gets public metadata for one video series.
func (v Client) SeriesInfo(ctx context.Context, params SeriesInfoParams) (SeriesInfo, error) {
	return sendVideoPayload[SeriesInfo](ctx, v.client, videoSeriesInfoEndpoint, "video.collection.series_info", params.EncodeQuery)
}

// SeriesArchives gets one public page of videos in a series.
func (v Client) SeriesArchives(ctx context.Context, params SeriesArchivesParams) (SeriesArchives, error) {
	return sendVideoPayload[SeriesArchives](ctx, v.client, videoSeriesArchivesEndpoint, "video.collection.series_archives", params.EncodeQuery)
}

func sendVideoPayload[T any](
	ctx context.Context,
	client *core.Client,
	endpoint, operation string,
	encodeQuery func() (url.Values, error),
) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "video client is not initialized"}
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

func sendVideoWBIPayload[T any](
	ctx context.Context,
	client *core.Client,
	endpoint, operation string,
	encodeQuery func() (url.Values, error),
) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "video client is not initialized"}
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
