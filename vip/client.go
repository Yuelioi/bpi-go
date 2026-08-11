package vip

import (
	"context"
	"net/http"

	core "github.com/Yuelioi/bpi-go/client"
)

const vipCenterEndpoint = "https://api.bilibili.com/x/vip/web/vip_center/combine"

// Client provides promoted VIP-center reads.
type Client struct{ client *core.Client }

// NewClient binds the vip module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (v Client) Center(ctx context.Context, params CenterParams) (Center, error) {
	var zero Center
	if v.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "VIP client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, vipCenterEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[Center](ctx, v.client, request, "vip.center_info")
}
