package manga

import (
	"context"
	"net/http"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	mangaSeasonInfoEndpoint    = "https://manga.bilibili.com/twirp/user.v1.Season/GetSeasonInfo"
	mangaClockInInfoEndpoint   = "https://manga.bilibili.com/twirp/activity.v1.Activity/GetClockInInfo"
	mangaUserPointEndpoint     = "https://manga.bilibili.com/twirp/pointshop.v1.Pointshop/GetUserPoint"
	mangaPointProductsEndpoint = "https://manga.bilibili.com/twirp/pointshop.v1.Pointshop/ListProduct"
	mangaCouponsEndpoint       = "https://manga.bilibili.com/twirp/user.v1.User/GetCoupons"
)

// Client provides the promoted Manga read operations.
type Client struct{ client *core.Client }

// NewClient binds the manga module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (m Client) SeasonInfo(ctx context.Context) (SeasonInfo, error) {
	return sendMangaEmptyPOST[SeasonInfo](ctx, m.client, mangaSeasonInfoEndpoint, "manga.season_info")
}

func (m Client) ClockInInfo(ctx context.Context) (ClockInInfo, error) {
	return sendMangaEmptyPOST[ClockInInfo](ctx, m.client, mangaClockInInfoEndpoint, "manga.clock_in_info")
}

func (m Client) UserPoint(ctx context.Context) (UserPoint, error) {
	return sendMangaEmptyPOST[UserPoint](ctx, m.client, mangaUserPointEndpoint, "manga.user_point")
}

func (m Client) PointProducts(ctx context.Context) ([]Product, error) {
	return sendMangaEmptyPOST[[]Product](ctx, m.client, mangaPointProductsEndpoint, "manga.point_products")
}

func (m Client) Coupons(ctx context.Context, params CouponsParams) (Coupons, error) {
	var zero Coupons
	if m.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "manga client is not initialized"}
	}
	body, err := params.RequestBody()
	if err != nil {
		return zero, err
	}
	request, err := core.NewJSONRequest(ctx, mangaCouponsEndpoint, nil, body)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[Coupons](ctx, m.client, request, "manga.coupons")
}

func sendMangaEmptyPOST[T any](ctx context.Context, client *core.Client, endpoint, operation string) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "manga client is not initialized"}
	}
	request, err := core.NewQueryRequest(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
