package fav

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	favFolderInfoEndpoint    = "https://api.bilibili.com/x/v3/fav/folder/info"
	favCreatedListEndpoint   = "https://api.bilibili.com/x/v3/fav/folder/created/list-all"
	favCollectedListEndpoint = "https://api.bilibili.com/x/v3/fav/folder/collected/list"
	favResourceInfosEndpoint = "https://api.bilibili.com/x/v3/fav/resource/infos"
	favListDetailEndpoint    = "https://api.bilibili.com/x/v3/fav/resource/list"
	favResourceIDsEndpoint   = "https://api.bilibili.com/x/v3/fav/resource/ids"
)

// Client provides favorites-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the fav module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (f Client) FolderInfo(ctx context.Context, params FolderInfoParams) (FolderInfo, error) {
	return sendFavPayload[FolderInfo](ctx, f.client, favFolderInfoEndpoint, "fav.folder_info", params.EncodeQuery)
}

func (f Client) CreatedList(ctx context.Context, params CreatedListParams) (CreatedList, error) {
	return sendFavPayload[CreatedList](ctx, f.client, favCreatedListEndpoint, "fav.created_list", params.EncodeQuery)
}

func (f Client) CollectedList(ctx context.Context, params CollectedListParams) (CollectedList, error) {
	return sendFavPayload[CollectedList](ctx, f.client, favCollectedListEndpoint, "fav.collected_list", params.EncodeQuery)
}

func (f Client) ResourceInfos(ctx context.Context, params ResourceInfosParams) ([]ResourceInfo, error) {
	return sendFavPayload[[]ResourceInfo](ctx, f.client, favResourceInfosEndpoint, "fav.resource_infos", params.EncodeQuery)
}

func (f Client) ListDetail(ctx context.Context, params ListDetailParams) (ListDetail, error) {
	return sendFavPayload[ListDetail](ctx, f.client, favListDetailEndpoint, "fav.list_detail", params.EncodeQuery)
}

func (f Client) ResourceIDs(ctx context.Context, params ResourceIDsParams) ([]ResourceID, error) {
	return sendFavPayload[[]ResourceID](ctx, f.client, favResourceIDsEndpoint, "fav.resource_ids", params.EncodeQuery)
}

func sendFavPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "favorite client is not initialized"}
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
