package user

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	userAlbumCountEndpoint     = "https://api.vc.bilibili.com/link_draw/v1/doc/upload_count"
	userBangumiFollowEndpoint  = "https://api.bilibili.com/x/space/bangumi/follow/list"
	userCardEndpoint           = "https://api.bilibili.com/x/web-interface/card"
	userCardsEndpoint          = "https://api.vc.bilibili.com/account/v1/user/cards"
	userFollowersEndpoint      = "https://api.bilibili.com/x/relation/fans"
	userFollowingsEndpoint     = "https://api.bilibili.com/x/relation/followings"
	userFollowTagsEndpoint     = "https://api.bilibili.com/x/relation/tags"
	userInfosEndpoint          = "https://api.vc.bilibili.com/x/im/user_infos"
	userMedalWallEndpoint      = "https://api.live.bilibili.com/xlive/web-ucenter/user/MedalWall"
	userNameToUIDEndpoint      = "https://api.bilibili.com/x/polymer/web-dynamic/v1/name-to-uid"
	userNavStatEndpoint        = "https://api.bilibili.com/x/space/navnum"
	userRelationStatEndpoint   = "https://api.bilibili.com/x/relation/stat"
	userSpaceInfoEndpoint      = "https://api.bilibili.com/x/space/wbi/acc/info"
	userSpaceNoticeEndpoint    = "https://api.bilibili.com/x/space/notice"
	userUpStatEndpoint         = "https://api.bilibili.com/x/space/upstat"
	userUploadedVideosEndpoint = "https://api.bilibili.com/x/space/wbi/arc/search"
)

// Client provides user profile, count, and relationship read operations.
type Client struct{ client *core.Client }

// NewClient binds the user module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (u Client) AlbumCount(ctx context.Context, params AlbumCountParams) (AlbumCount, error) {
	return sendUserPayload[AlbumCount](ctx, u.client, userAlbumCountEndpoint, "user.album_count", params.EncodeQuery)
}

func (u Client) BangumiFollowList(ctx context.Context, params BangumiFollowListParams) (BangumiFollowList, error) {
	return sendUserPayload[BangumiFollowList](ctx, u.client, userBangumiFollowEndpoint, "user.bangumi_follow_list", params.EncodeQuery)
}

func (u Client) Card(ctx context.Context, params CardParams) (CardProfile, error) {
	return sendUserPayload[CardProfile](ctx, u.client, userCardEndpoint, "user.card", params.EncodeQuery)
}

func (u Client) Cards(ctx context.Context, params CardsParams) ([]BatchCard, error) {
	return sendUserPayload[[]BatchCard](ctx, u.client, userCardsEndpoint, "user.cards", params.EncodeQuery)
}

func (u Client) FollowTags(ctx context.Context) ([]FollowTag, error) {
	return sendUserNoParams[[]FollowTag](ctx, u.client, userFollowTagsEndpoint, "user.follow_tags")
}

func (u Client) Followers(ctx context.Context, params FollowersParams) (Followers, error) {
	return sendUserPayload[Followers](ctx, u.client, userFollowersEndpoint, "user.followers", params.EncodeQuery)
}

func (u Client) Followings(ctx context.Context, params FollowingsParams) (Followings, error) {
	return sendUserPayload[Followings](ctx, u.client, userFollowingsEndpoint, "user.followings", params.EncodeQuery)
}

func (u Client) Infos(ctx context.Context, params InfosParams) ([]BatchInfo, error) {
	return sendUserPayload[[]BatchInfo](ctx, u.client, userInfosEndpoint, "user.infos", params.EncodeQuery)
}

func (u Client) MedalWall(ctx context.Context, params MedalWallParams) (MedalWall, error) {
	return sendUserPayload[MedalWall](ctx, u.client, userMedalWallEndpoint, "user.medal_wall", params.EncodeQuery)
}

func (u Client) NameToUID(ctx context.Context, params NameToUIDParams) (NameToUID, error) {
	return sendUserPayload[NameToUID](ctx, u.client, userNameToUIDEndpoint, "user.name_to_uid", params.EncodeQuery)
}

func (u Client) NavStat(ctx context.Context, params NavStatParams) (NavStat, error) {
	return sendUserPayload[NavStat](ctx, u.client, userNavStatEndpoint, "user.nav_stat", params.EncodeQuery)
}

func (u Client) RelationStat(ctx context.Context, params RelationStatParams) (RelationStat, error) {
	return sendUserPayload[RelationStat](ctx, u.client, userRelationStatEndpoint, "user.relation_stat", params.EncodeQuery)
}

func (u Client) SpaceInfo(ctx context.Context, params SpaceParams) (SpaceProfile, error) {
	return sendUserWBI[SpaceProfile](ctx, u.client, userSpaceInfoEndpoint, "user.space_info", params.EncodeQuery)
}

func (u Client) SpaceNotice(ctx context.Context, params SpaceNoticeParams) (SpaceNotice, error) {
	return sendUserPayload[SpaceNotice](ctx, u.client, userSpaceNoticeEndpoint, "user.space_notice", params.EncodeQuery)
}

func (u Client) UpStat(ctx context.Context, params UpStatParams) (UpStat, error) {
	return sendUserPayload[UpStat](ctx, u.client, userUpStatEndpoint, "user.up_stat", params.EncodeQuery)
}

func (u Client) UploadedVideos(ctx context.Context, params UploadedVideosParams) (UploadedVideos, error) {
	return sendUserWBI[UploadedVideos](ctx, u.client, userUploadedVideosEndpoint, "user.uploaded_videos", params.EncodeQuery)
}

func sendUserNoParams[T any](ctx context.Context, client *core.Client, endpoint, operation string) (T, error) {
	return sendUserPayload[T](ctx, client, endpoint, operation, func() (url.Values, error) { return nil, nil })
}

func sendUserPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "user client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}

func sendUserWBI[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "user client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return zero, err
	}
	signed, err := client.WBIValues(ctx, query)
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, signed)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
