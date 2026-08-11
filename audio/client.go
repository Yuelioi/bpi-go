package audio

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	audioInfoEndpoint             = "https://www.bilibili.com/audio/music-service-c/web/song/info"
	audioTagsEndpoint             = "https://www.bilibili.com/audio/music-service-c/web/tag/song"
	audioMembersEndpoint          = "https://www.bilibili.com/audio/music-service-c/web/member/song"
	audioLyricEndpoint            = "https://www.bilibili.com/audio/music-service-c/web/song/lyric"
	audioStatusNumberEndpoint     = "https://www.bilibili.com/audio/music-service-c/web/stat/song"
	audioCollectionStatusEndpoint = "https://www.bilibili.com/audio/music-service-c/web/collections/songs-coll"
	audioCoinCountEndpoint        = "https://www.bilibili.com/audio/music-service-c/web/coin/audio"
	audioStreamURLWebEndpoint     = "https://www.bilibili.com/audio/music-service-c/web/url"
	audioStreamURLEndpoint        = "https://api.bilibili.com/audio/music-service-c/url"
	audioCollectionsListEndpoint  = "https://www.bilibili.com/audio/music-service-c/web/collections/list"
	audioCollectionInfoEndpoint   = "https://www.bilibili.com/audio/music-service-c/web/collections/info"
	audioHotMenuEndpoint          = "https://www.bilibili.com/audio/music-service-c/web/menu/hit"
	audioRankMenuEndpoint         = "https://www.bilibili.com/audio/music-service-c/web/menu/rank"
	audioRankPeriodEndpoint       = "https://api.bilibili.com/x/copyright-music-publicity/toplist/all_period"
	audioRankDetailEndpoint       = "https://api.bilibili.com/x/copyright-music-publicity/toplist/detail"
	audioRankMusicListEndpoint    = "https://api.bilibili.com/x/copyright-music-publicity/toplist/music_list"
)

// Client provides audio-domain operations.
type Client struct{ client *core.Client }

// NewClient binds the audio module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (a Client) Info(ctx context.Context, params SongParams) (Info, error) {
	return sendAudioPayload[Info](ctx, a.client, audioInfoEndpoint, "audio.info", params.EncodeQuery)
}

func (a Client) Tags(ctx context.Context, params SongParams) ([]Tag, error) {
	return sendAudioPayload[[]Tag](ctx, a.client, audioTagsEndpoint, "audio.tags", params.EncodeQuery)
}

func (a Client) Members(ctx context.Context, params SongParams) ([]MemberGroup, error) {
	return sendAudioPayload[[]MemberGroup](ctx, a.client, audioMembersEndpoint, "audio.members", params.EncodeQuery)
}

func (a Client) Lyric(ctx context.Context, params SongParams) (string, error) {
	return sendAudioPayload[string](ctx, a.client, audioLyricEndpoint, "audio.lyric", params.EncodeQuery)
}

func (a Client) StatusNumber(ctx context.Context, params SongParams) (StatusNumber, error) {
	return sendAudioPayload[StatusNumber](ctx, a.client, audioStatusNumberEndpoint, "audio.status_number", params.EncodeQuery)
}

func (a Client) CollectionStatus(ctx context.Context, params SongParams) (bool, error) {
	return sendAudioPayload[bool](ctx, a.client, audioCollectionStatusEndpoint, "audio.collection_status", params.EncodeQuery)
}

func (a Client) CoinCount(ctx context.Context, params SongParams) (int32, error) {
	return sendAudioPayload[int32](ctx, a.client, audioCoinCountEndpoint, "audio.coin_count", params.EncodeQuery)
}

func (a Client) StreamURLWeb(ctx context.Context, params StreamURLWebParams) (StreamURLWeb, error) {
	return sendAudioPayload[StreamURLWeb](ctx, a.client, audioStreamURLWebEndpoint, "audio.stream_url_web", params.EncodeQuery)
}

func (a Client) StreamURL(ctx context.Context, params StreamURLParams) (StreamURL, error) {
	return sendAudioPayload[StreamURL](ctx, a.client, audioStreamURLEndpoint, "audio.stream_url", params.EncodeQuery)
}

func (a Client) CollectionsList(ctx context.Context, params PageParams) (Page[Collection], error) {
	return sendAudioPayload[Page[Collection]](ctx, a.client, audioCollectionsListEndpoint, "audio.collections_list", params.EncodeQuery)
}

func (a Client) CollectionInfo(ctx context.Context, params CollectionInfoParams) (*Collection, error) {
	return sendAudioOptional[Collection](ctx, a.client, audioCollectionInfoEndpoint, "audio.collection_info", params.EncodeQuery)
}

func (a Client) HotMenu(ctx context.Context, params PageParams) (Page[HotMenu], error) {
	return sendAudioPayload[Page[HotMenu]](ctx, a.client, audioHotMenuEndpoint, "audio.hot_menu", params.EncodeQuery)
}

func (a Client) RankMenu(ctx context.Context, params PageParams) (Page[RankMenu], error) {
	return sendAudioPayload[Page[RankMenu]](ctx, a.client, audioRankMenuEndpoint, "audio.rank_menu", params.EncodeQuery)
}

func (a Client) RankPeriod(ctx context.Context, params RankPeriodParams) (RankPeriods, error) {
	return sendAudioCSRF[RankPeriods](ctx, a.client, audioRankPeriodEndpoint, "audio.rank_period", params.EncodeQuery)
}

func (a Client) RankDetail(ctx context.Context, params RankListParams) (RankDetail, error) {
	return sendAudioCSRF[RankDetail](ctx, a.client, audioRankDetailEndpoint, "audio.rank_detail", params.EncodeQuery)
}

func (a Client) RankMusicList(ctx context.Context, params RankListParams) (RankMusicList, error) {
	return sendAudioCSRF[RankMusicList](ctx, a.client, audioRankMusicListEndpoint, "audio.rank_music_list", params.EncodeQuery)
}

func sendAudioPayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "audio client is not initialized"}
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

func sendAudioOptional[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func() (url.Values, error)) (*T, error) {
	if client == nil {
		return nil, &core.ParameterError{Field: "client", Message: "audio client is not initialized"}
	}
	query, err := encodeQuery()
	if err != nil {
		return nil, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return nil, err
	}
	return core.SendOptionalPayload[T](ctx, client, request, operation)
}

func sendAudioCSRF[T any](ctx context.Context, client *core.Client, endpoint, operation string, encodeQuery func(string) (url.Values, error)) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "audio client is not initialized"}
	}
	csrf, _ := client.CSRF()
	query, err := encodeQuery(csrf)
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}
