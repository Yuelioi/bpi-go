package cheese

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
	"github.com/Yuelioi/bpi-go/ids"
)

const (
	cheeseEpisodeListEndpoint = "https://api.bilibili.com/pugv/view/web/ep/list"
	cheeseSeasonEndpoint      = "https://api.bilibili.com/pugv/view/web/season"
	cheesePlayURLEndpoint     = "https://api.bilibili.com/pugv/player/web/playurl"
)

// Client provides public PUGV/course operations.
type Client struct{ client *core.Client }

// NewClient binds the cheese module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// EpisodeList gets one page of episodes in a course season.
func (c Client) EpisodeList(ctx context.Context, params EpisodeListParams) (EpisodeList, error) {
	return sendCheesePayload[EpisodeList](ctx, c.client, cheeseEpisodeListEndpoint, "cheese.info.ep_list", params.EncodeQuery)
}

// SeasonDetail gets course-season detail by season ID.
func (c Client) SeasonDetail(ctx context.Context, seasonID ids.SeasonID) (Course, error) {
	return sendCheesePayload[Course](ctx, c.client, cheeseSeasonEndpoint, "cheese.info.season_detail_by_season_id", func() (url.Values, error) {
		if err := seasonID.Validate(); err != nil {
			return nil, &core.ParameterError{Field: "season_id", Message: "season ID is invalid"}
		}
		return url.Values{"season_id": {seasonID.String()}}, nil
	})
}

// EpisodeDetail gets the containing course-season detail by episode ID.
func (c Client) EpisodeDetail(ctx context.Context, episodeID ids.EpisodeID) (Course, error) {
	return sendCheesePayload[Course](ctx, c.client, cheeseSeasonEndpoint, "cheese.info.season_detail_by_ep_id", func() (url.Values, error) {
		if err := episodeID.Validate(); err != nil {
			return nil, &core.ParameterError{Field: "ep_id", Message: "episode ID is invalid"}
		}
		return url.Values{"ep_id": {episodeID.String()}}, nil
	})
}

// PlayURL gets a course episode's playback payload.
func (c Client) PlayURL(ctx context.Context, params PlayURLParams) (PlayURL, error) {
	return sendCheesePayload[PlayURL](ctx, c.client, cheesePlayURLEndpoint, "cheese.playurl", params.EncodeQuery)
}

func sendCheesePayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "cheese client is not initialized"}
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
