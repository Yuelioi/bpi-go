package wallet

import (
	"context"

	core "github.com/Yuelioi/bpi-go/client"
)

const walletInfoEndpoint = "https://pay.bilibili.com/paywallet/wallet/getUserWallet"

// Client provides explicitly authenticated wallet reads.
type Client struct{ client *core.Client }

// NewClient binds the wallet module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (w Client) Info(ctx context.Context, params InfoParams) (Info, error) {
	var zero Info
	if w.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "wallet client is not initialized"}
	}
	csrf, err := w.client.CSRF()
	if err != nil {
		return zero, err
	}
	body, err := params.RequestBody(csrf)
	if err != nil {
		return zero, err
	}
	request, err := core.NewJSONRequest(ctx, walletInfoEndpoint, nil, body)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[Info](ctx, w.client, request, "wallet.info")
}
