package activity

import (
	"context"
	"net/http"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	activityInfoEndpoint = "https://api.bilibili.com/x/activity/subject/info"
	activityListEndpoint = "https://api.bilibili.com/x/activity/page/list"
)

// Client provides activity-domain operations.
type Client struct {
	client *core.Client
}

// NewClient binds the activity module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// Info gets an activity subject by ID and optional source video.
func (a Client) Info(ctx context.Context, params InfoParams) (Info, error) {
	var zero Info
	if a.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "activity client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, activityInfoEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[Info](ctx, a.client, request, "activity.info")
}

// List gets one page of activities.
func (a Client) List(ctx context.Context, params ListParams) (ListPage, error) {
	var zero ListPage
	if a.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "activity client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, activityListEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[ListPage](ctx, a.client, request, "activity.list")
}

// ListDefault gets the first page using Bilibili's Web defaults.
func (a Client) ListDefault(ctx context.Context) (ListPage, error) {
	return a.List(ctx, NewListParams())
}
