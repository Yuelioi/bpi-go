package webwidget

import (
	"context"
	"net/http"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	webWidgetRegionBannerEndpoint = "https://api.bilibili.com/x/web-show/region/banner"
	webWidgetHeaderPageEndpoint   = "https://api.bilibili.com/x/web-show/page/header"
	webWidgetOnlineEndpoint       = "https://api.bilibili.com/x/web-interface/online"
)

// Client provides Web-widget operations.
type Client struct{ client *core.Client }

// NewClient binds the webwidget module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

// RegionBanner gets the public carousel for a video region.
func (w Client) RegionBanner(ctx context.Context, params RegionBannerParams) (RegionBannerData, error) {
	var zero RegionBannerData
	if w.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "web-widget client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, webWidgetRegionBannerEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[RegionBannerData](ctx, w.client, request, "web_widget.region_banner")
}

// HeaderPage gets the public home-page header and parses its embedded layer
// definition during response decoding.
func (w Client) HeaderPage(ctx context.Context, params HeaderPageParams) (HeaderData, error) {
	var zero HeaderData
	if w.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "web-widget client is not initialized"}
	}
	query, err := params.EncodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, webWidgetHeaderPageEndpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[HeaderData](ctx, w.client, request, "web_widget.header_page")
}

// Online gets today's upload counts grouped by video region.
func (w Client) Online(ctx context.Context) (OnlineData, error) {
	var zero OnlineData
	if w.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "web-widget client is not initialized"}
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, webWidgetOnlineEndpoint, nil)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[OnlineData](ctx, w.client, request, "web_widget.online")
}
