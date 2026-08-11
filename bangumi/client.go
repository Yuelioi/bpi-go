package bangumi

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
	"github.com/Yuelioi/bpi-go/ids"
)

const bangumiTimelineEndpoint = "https://api.bilibili.com/pgc/web/timeline"

const (
	bangumiInfoEndpoint     = "https://api.bilibili.com/pgc/review/user"
	bangumiDetailEndpoint   = "https://api.bilibili.com/pgc/view/web/season"
	bangumiSectionsEndpoint = "https://api.bilibili.com/pgc/web/season/section"
	bangumiPlayURLEndpoint  = "https://api.bilibili.com/pgc/player/web/playurl"
)

// Client provides bangumi-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the bangumi module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// Review gets the public review and media summary for a media ID.
func (b Client) Review(ctx context.Context, params InfoParams) (Info, error) {
	return sendBangumiPayload[Info](ctx, b.client, bangumiInfoEndpoint, "bangumi.info.review_user", params.EncodeQuery)
}

// SeasonDetail gets season detail by season ID.
func (b Client) SeasonDetail(ctx context.Context, seasonID ids.SeasonID) (Detail, error) {
	params := DetailBySeasonID(seasonID)
	return sendBangumiPayload[Detail](ctx, b.client, bangumiDetailEndpoint, "bangumi.info.season_detail_by_season_id", params.EncodeQuery)
}

// EpisodeDetail gets season detail by episode ID.
func (b Client) EpisodeDetail(ctx context.Context, episodeID ids.EpisodeID) (Detail, error) {
	params := DetailByEpisodeID(episodeID)
	return sendBangumiPayload[Detail](ctx, b.client, bangumiDetailEndpoint, "bangumi.info.season_detail_by_ep_id", params.EncodeQuery)
}

// Sections gets the main and auxiliary sections for a season.
func (b Client) Sections(ctx context.Context, params SectionsParams) (Sections, error) {
	return sendBangumiPayload[Sections](ctx, b.client, bangumiSectionsEndpoint, "bangumi.info.season_section", params.EncodeQuery)
}

// PlayURL gets the playback payload for an episode or content ID.
func (b Client) PlayURL(ctx context.Context, params PlayURLParams) (PlayURL, error) {
	return sendBangumiPayload[PlayURL](ctx, b.client, bangumiPlayURLEndpoint, "bangumi.playurl", params.EncodeQuery)
}

// Timeline gets a bangumi or movie release timeline. The endpoint's `result`
// payload alias is normalized by the root response decoder.
func (b Client) Timeline(ctx context.Context, params TimelineParams) ([]TimelineDay, error) {
	return sendBangumiPayload[[]TimelineDay](ctx, b.client, bangumiTimelineEndpoint, "bangumi.timeline", params.EncodeQuery)
}

func sendBangumiPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "bangumi client is not initialized"}
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
