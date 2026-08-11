package clientinfo

import (
	"context"
	"net/http"

	core "github.com/Yuelioi/bpi-go/client"
)

const clientInfoIPEndpoint = "https://api.live.bilibili.com/ip_service/v1/ip_service/get_ip_addr"

// Client provides client-information operations.
type Client struct{ client *core.Client }

// NewClient binds the clientinfo module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// IP gets public geolocation information for an optional IP address.
func (c Client) IP(ctx context.Context, params IPParams) (IPInfo, error) {
	var zero IPInfo
	if c.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "client-info client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, clientInfoIPEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[IPInfo](ctx, c.client, request, "clientinfo.ip")
}
