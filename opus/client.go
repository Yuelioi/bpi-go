package opus

import (
	"context"
	"net/http"

	core "github.com/Yuelioi/bpi-go/client"
)

const opusSpaceFeedEndpoint = "https://api.bilibili.com/x/polymer/web-dynamic/v1/opus/feed/space"

// Client provides public opus operations.
type Client struct{ client *core.Client }

// NewClient binds the opus module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// SpaceFeed gets one page of opus entries from a member's public space.
func (o Client) SpaceFeed(ctx context.Context, params SpaceFeedParams) (SpaceFeed, error) {
	var zero SpaceFeed
	if o.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "opus client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, opusSpaceFeedEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[SpaceFeed](ctx, o.client, request, "opus.space_feed")
}
