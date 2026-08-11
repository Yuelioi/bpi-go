package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	searchTypedEndpoint    = "https://api.bilibili.com/x/web-interface/wbi/search/type"
	searchDefaultEndpoint  = "https://api.bilibili.com/x/web-interface/wbi/search/default"
	searchSuggestEndpoint  = "https://s.search.bilibili.com/main/suggest"
	searchHotWordsEndpoint = "https://s.search.bilibili.com/main/hotword"
)

// Client provides public search operations.
type Client struct{ client *core.Client }

// NewClient binds the search module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (s Client) Article(ctx context.Context, params ArticleParams) (Data[[]Article], error) {
	return sendSearchWBI[Data[[]Article]](ctx, s.client, "search.article", params.EncodeQuery)
}

func (s Client) Bangumi(ctx context.Context, params BangumiParams) (Data[[]Bangumi], error) {
	return sendSearchWBI[Data[[]Bangumi]](ctx, s.client, "search.bangumi", params.EncodeQuery)
}

func (s Client) Users(ctx context.Context, params UserParams) (Data[[]User], error) {
	return sendSearchWBI[Data[[]User]](ctx, s.client, "search.bili_user", params.EncodeQuery)
}

func (s Client) Live(ctx context.Context, params LiveParams) (Data[LiveData], error) {
	return sendSearchWBI[Data[LiveData]](ctx, s.client, "search.live", params.EncodeQuery)
}

func (s Client) LiveRooms(ctx context.Context, params LiveRoomParams) (Data[[]LiveRoom], error) {
	return sendSearchWBI[Data[[]LiveRoom]](ctx, s.client, "search.live_room", params.EncodeQuery)
}

func (s Client) LiveUsers(ctx context.Context, params LiveUserParams) (Data[[]LiveUser], error) {
	return sendSearchWBI[Data[[]LiveUser]](ctx, s.client, "search.live_user", params.EncodeQuery)
}

func (s Client) Movies(ctx context.Context, params MovieParams) (Data[[]Movie], error) {
	return sendSearchWBI[Data[[]Movie]](ctx, s.client, "search.movie", params.EncodeQuery)
}

func (s Client) Videos(ctx context.Context, params VideoParams) (Data[[]Video], error) {
	return sendSearchWBI[Data[[]Video]](ctx, s.client, "search.video", params.EncodeQuery)
}

// Default gets the WBI-signed default Web search target.
func (s Client) Default(ctx context.Context) (Default, error) {
	return sendSearchWBI[Default](ctx, s.client, "search.default", func() (url.Values, error) {
		return url.Values{"foo": {"bar"}}, nil
	})
}

// Suggest gets suggestions for a partial search term.
func (s Client) Suggest(ctx context.Context, params SuggestParams) (Suggest, error) {
	return sendSearchPayload[Suggest](ctx, s.client, searchSuggestEndpoint, "search.suggest", params.EncodeQuery)
}

// HotWords gets the raw public hot-word response. Unlike ordinary Bilibili
// endpoints, this service places list directly beside code without data.
func (s Client) HotWords(ctx context.Context) (HotWords, error) {
	var zero HotWords
	if s.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "search client is not initialized"}
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, searchHotWordsEndpoint, nil)
	if err != nil {
		return zero, err
	}
	response, err := s.client.Do(ctx, request, "search.hotwords")
	if err != nil {
		return zero, err
	}
	var payload HotWords
	if err := json.Unmarshal(response.Body, &payload); err != nil {
		return zero, core.NewResponseDecodeError(err, response.Body)
	}
	if payload.Code != 0 {
		return zero, &core.APIError{Code: int(payload.Code), Message: "hot-word request failed"}
	}
	return payload, nil
}

func sendSearchWBI[T any](ctx context.Context, client *core.Client, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "search client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	signed, err := client.WBIValues(ctx, query)
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, searchEndpoint(operation), signed)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}

func sendSearchPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "search client is not initialized"}
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

func searchEndpoint(operation string) string {
	if operation == "search.default" {
		return searchDefaultEndpoint
	}
	return searchTypedEndpoint
}
