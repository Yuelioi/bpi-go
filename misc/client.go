package misc

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"

	"github.com/Yuelioi/bpi-go/internal/sign"
)

const (
	miscBuvid3Endpoint    = "https://api.bilibili.com/x/web-frontend/getbuvid"
	miscBuvidEndpoint     = "https://api.bilibili.com/x/frontend/finger/spi"
	miscShortLinkEndpoint = "https://api.biliapi.net/x/share/click"
	miscTicketEndpoint    = "https://api.bilibili.com/bapis/bilibili.api.ticket.v1.Ticket/GenWebTicket"
)

// Client provides promoted miscellaneous operations.
type Client struct{ client *core.Client }

// NewClient binds the misc module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (m Client) Buvid3(ctx context.Context) (Buvid3, error) {
	return sendMiscPayload[Buvid3](ctx, m.client, http.MethodGet, miscBuvid3Endpoint, "misc.buvid3", nil)
}

func (m Client) Buvid(ctx context.Context) (Buvid, error) {
	return sendMiscPayload[Buvid](ctx, m.client, http.MethodGet, miscBuvidEndpoint, "misc.buvid", nil)
}

func (m Client) ShortLink(ctx context.Context, params ShortLinkParams) (ShortLink, error) {
	var zero ShortLink
	if m.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "misc client is not initialized"}
	}
	form, err := params.EncodeForm()
	if err != nil {
		return zero, err
	}
	request, err := core.NewFormRequest(ctx, miscShortLinkEndpoint, nil, form)
	if err != nil {
		return zero, err
	}
	// The short-link host is not a bilibili.com host, so explicitly add public
	// browser headers while retaining the client's no-cross-host Cookie rule.
	request.Header.Set("Referer", m.client.Referer())
	request.Header.Set("Origin", m.client.Origin())
	result, err := core.SendPayload[ShortLink](ctx, m.client, request, "misc.b23tv.short_link")
	if err != nil {
		return zero, err
	}
	result.Extract()
	return result, nil
}

func (m Client) BiliTicket(ctx context.Context) (Ticket, error) {
	var zero Ticket
	if m.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "misc client is not initialized"}
	}
	csrf, _ := m.client.CSRF()
	parameters := sign.TicketRequestParams(uint64(m.client.Now().Unix()), csrf)
	query := make(url.Values, len(parameters))
	for key, value := range parameters {
		query.Set(key, value)
	}
	return sendMiscPayload[Ticket](ctx, m.client, http.MethodPost, miscTicketEndpoint, "misc.bili_ticket", query)
}

func (m Client) BiliTicketString(ctx context.Context) (string, error) {
	ticket, err := m.BiliTicket(ctx)
	if err != nil {
		return "", err
	}
	return ticket.Ticket, nil
}

func sendMiscPayload[T any](ctx context.Context, client *core.Client, method, endpoint, operation string, query url.Values) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "misc client is not initialized"}
	}
	request, err := core.NewQueryRequest(ctx, method, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
