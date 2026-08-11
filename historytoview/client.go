package historytoview

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	historyListEndpoint   = "https://api.bilibili.com/x/web-interface/history/cursor"
	historyShadowEndpoint = "https://api.bilibili.com/x/v2/history/shadow"
	toViewListEndpoint    = "https://api.bilibili.com/x/v2/history/toview"
)

// Client provides viewing-history and watch-later operations.
type Client struct{ client *core.Client }

// NewClient binds the historytoview module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (h Client) HistoryList(ctx context.Context, params ListParams) (HistoryList, error) {
	return sendHistoryPayload[HistoryList](ctx, h.client, historyListEndpoint, "historytoview.history_list", params.EncodeQuery)
}

func (h Client) HistoryShadow(ctx context.Context) (bool, error) {
	return sendHistoryPayload[bool](ctx, h.client, historyShadowEndpoint, "historytoview.history_shadow", func() (url.Values, error) { return nil, nil })
}

func (h Client) ToViewList(ctx context.Context) (ToViewList, error) {
	return sendHistoryPayload[ToViewList](ctx, h.client, toViewListEndpoint, "historytoview.toview_list", func() (url.Values, error) { return nil, nil })
}

func sendHistoryPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "history client is not initialized"}
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
