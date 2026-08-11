package live

import (
	"context"
	"net/http"
	"net/url"

	core "github.com/Yuelioi/bpi-go/client"
)

const (
	liveAreaListEndpoint       = "https://api.live.bilibili.com/room/v1/Area/getList"
	liveRoomInfoEndpoint       = "https://api.live.bilibili.com/room/v1/Room/get_info"
	liveStreamEndpoint         = "https://api.live.bilibili.com/room/v1/Room/playUrl"
	liveRecommendEndpoint      = "https://api.live.bilibili.com/xlive/web-interface/v1/webMain/getMoreRecList"
	liveVersionEndpoint        = "https://api.live.bilibili.com/xlive/app-blink/v1/liveVersionInfo/getHomePageLiveVersion"
	liveGiftTypesEndpoint      = "https://api.live.bilibili.com/gift/v1/master/getGiftTypes"
	liveRoomGiftListEndpoint   = "https://api.live.bilibili.com/xlive/web-room/v1/giftPanel/roomGiftList"
	liveBlindGiftInfoEndpoint  = "https://api.live.bilibili.com/xlive/general-interface/v1/blindFirstWin/getInfo"
	liveDanmuInfoEndpoint      = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo"
	liveEmoticonsEndpoint      = "https://api.live.bilibili.com/xlive/web-ucenter/v2/emoticon/GetEmoticons"
	liveLotteryInfoEndpoint    = "https://api.live.bilibili.com/xlive/lottery-interface/v1/lottery/getLotteryInfoWeb"
	liveMyMedalsEndpoint       = "https://api.live.bilibili.com/xlive/app-ucenter/v1/user/GetMyMedals"
	liveFollowUpListEndpoint   = "https://api.live.bilibili.com/xlive/web-ucenter/user/following"
	liveFollowUpWebEndpoint    = "https://api.live.bilibili.com/xlive/web-ucenter/v1/xfetter/GetWebList"
	liveReplayListEndpoint     = "https://api.live.bilibili.com/xlive/app-blink/v1/anchorVideo/AnchorGetReplayList"
	liveGuardListEndpoint      = "https://api.live.bilibili.com/xlive/app-room/v2/guardTab/topListNew"
	liveSilentUsersEndpoint    = "https://api.live.bilibili.com/xlive/web-ucenter/v1/banned/GetSilentUserList"
	liveBannedUsersEndpoint    = "https://api.live.bilibili.com/xlive/app-ucenter/v2/xbanned/banned/GetBlackList"
	liveShieldKeywordsEndpoint = "https://api.live.bilibili.com/xlive/app-ucenter/v1/banned/GetShieldKeywordList"
	liveWebHeartBeatEndpoint   = "https://live-trace.bilibili.com/xlive/rdata-interface/v1/heartbeat/webHeartBeat"
)

// Client provides the promoted read-only live API surface. It is a
// lightweight view over Client and is safe to use concurrently.
type Client struct{ client *core.Client }

// NewClient binds the live module to a shared client.
func NewClient(client *core.Client) Client { return Client{client: client} }

func (l Client) AreaList(ctx context.Context) ([]ParentArea, error) {
	return sendLivePayload[[]ParentArea](ctx, l.client, liveAreaListEndpoint, "live.area_list", noLiveQuery, false)
}

func (l Client) RoomInfo(ctx context.Context, params RoomInfoParams) (RoomInfo, error) {
	return sendLivePayload[RoomInfo](ctx, l.client, liveRoomInfoEndpoint, "live.room_info", params.EncodeQuery, false)
}

func (l Client) Stream(ctx context.Context, params StreamParams) (Stream, error) {
	return sendLivePayload[Stream](ctx, l.client, liveStreamEndpoint, "live.stream", params.EncodeQuery, false)
}

func (l Client) Recommend(ctx context.Context) (Recommend, error) {
	encode := func() (url.Values, error) {
		return url.Values{"platform": {"web"}, "web_location": {"333.1007"}}, nil
	}
	return sendLivePayload[Recommend](ctx, l.client, liveRecommendEndpoint, "live.recommend", encode, false)
}

func (l Client) Version(ctx context.Context) (Version, error) {
	encode := func() (url.Values, error) { return url.Values{"system_version": {"2"}}, nil }
	return sendLivePayload[Version](ctx, l.client, liveVersionEndpoint, "live.version", encode, false)
}

func (l Client) GiftTypes(ctx context.Context) ([]GiftType, error) {
	return sendLivePayload[[]GiftType](ctx, l.client, liveGiftTypesEndpoint, "live.gift_types", noLiveQuery, false)
}

func (l Client) RoomGiftList(ctx context.Context, params RoomGiftListParams) (RoomGiftList, error) {
	return sendLivePayload[RoomGiftList](ctx, l.client, liveRoomGiftListEndpoint, "live.room_gift_list", params.EncodeQuery, false)
}

func (l Client) BlindGiftInfo(ctx context.Context, params BlindGiftInfoParams) (BlindGiftInfo, error) {
	return sendLivePayload[BlindGiftInfo](ctx, l.client, liveBlindGiftInfoEndpoint, "live.blind_gift_info", params.EncodeQuery, false)
}

func (l Client) DanmuInfo(ctx context.Context, params DanmuInfoParams) (DanmuInfo, error) {
	return sendLivePayload[DanmuInfo](ctx, l.client, liveDanmuInfoEndpoint, "live.danmu_info", params.EncodeQuery, true)
}

func (l Client) Emoticons(ctx context.Context, params EmoticonsParams) (EmoticonData, error) {
	return sendLivePayload[EmoticonData](ctx, l.client, liveEmoticonsEndpoint, "live.emoticons", params.EncodeQuery, false)
}

func (l Client) LotteryInfo(ctx context.Context, params LotteryInfoParams) (LotteryInfo, error) {
	return sendLivePayload[LotteryInfo](ctx, l.client, liveLotteryInfoEndpoint, "live.lottery_info", params.EncodeQuery, true)
}

func (l Client) MyMedals(ctx context.Context, params MyMedalsParams) (MyMedals, error) {
	return sendLivePayload[MyMedals](ctx, l.client, liveMyMedalsEndpoint, "live.my_medals", params.EncodeQuery, false)
}

func (l Client) FollowUpList(ctx context.Context, params FollowUpListParams) (FollowUpList, error) {
	return sendLivePayload[FollowUpList](ctx, l.client, liveFollowUpListEndpoint, "live.follow_up_list", params.EncodeQuery, false)
}

func (l Client) FollowUpWebList(ctx context.Context, params FollowUpWebListParams) (FollowUpWebList, error) {
	return sendLivePayload[FollowUpWebList](ctx, l.client, liveFollowUpWebEndpoint, "live.follow_up_web_list", params.EncodeQuery, false)
}

func (l Client) ReplayList(ctx context.Context, params ReplayListParams) (ReplayList, error) {
	return sendLivePayload[ReplayList](ctx, l.client, liveReplayListEndpoint, "live.replay_list", params.EncodeQuery, false)
}

func (l Client) GuardList(ctx context.Context, params GuardListParams) (GuardList, error) {
	return sendLivePayload[GuardList](ctx, l.client, liveGuardListEndpoint, "live.guard_list", params.EncodeQuery, false)
}

func (l Client) SilentUsers(ctx context.Context, params SilentUsersParams) (SilentUsers, error) {
	var zero SilentUsers
	if l.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "live client is not initialized"}
	}
	form, err := params.EncodeForm(optionalLiveCSRF(l.client))
	if err != nil {
		return zero, err
	}
	referer, err := params.Referer()
	if err != nil {
		return zero, err
	}
	request, err := core.NewFormRequest(ctx, liveSilentUsersEndpoint, nil, form)
	if err != nil {
		return zero, err
	}
	setLiveRoomHeaders(request, referer)
	return core.SendPayload[SilentUsers](ctx, l.client, request, "live.silent_users")
}

func (l Client) BannedUsers(ctx context.Context, params BannedUsersParams) (BannedUsers, error) {
	var zero BannedUsers
	if l.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "live client is not initialized"}
	}
	query, err := params.EncodeQuery(optionalLiveCSRF(l.client))
	if err != nil {
		return zero, err
	}
	referer, err := params.Referer()
	if err != nil {
		return zero, err
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, liveBannedUsersEndpoint, query)
	if err != nil {
		return zero, err
	}
	setLiveRoomHeaders(request, referer)
	return core.SendPayload[BannedUsers](ctx, l.client, request, "live.banned_users")
}

func (l Client) ShieldKeywords(ctx context.Context, params ShieldKeywordsParams) (ShieldKeywords, error) {
	var zero ShieldKeywords
	if l.client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "live client is not initialized"}
	}
	form, err := params.EncodeForm(optionalLiveCSRF(l.client))
	if err != nil {
		return zero, err
	}
	referer, err := params.Referer()
	if err != nil {
		return zero, err
	}
	request, err := core.NewFormRequest(ctx, liveShieldKeywordsEndpoint, nil, form)
	if err != nil {
		return zero, err
	}
	setLiveRoomHeaders(request, referer)
	return core.SendPayload[ShieldKeywords](ctx, l.client, request, "live.shield_keywords")
}

func (l Client) WebHeartBeat(ctx context.Context, params WebHeartBeatParams) (HeartBeat, error) {
	return sendLivePayload[HeartBeat](ctx, l.client, liveWebHeartBeatEndpoint, "live.web_heart_beat", params.EncodeQuery, false)
}

func sendLivePayload[T any](ctx context.Context, client *core.Client, endpoint, operation string, encode func() (url.Values, error), wbi bool) (T, error) {
	var zero T
	if client == nil {
		return zero, &core.ParameterError{Field: "client", Message: "live client is not initialized"}
	}
	query, err := encode()
	if err != nil {
		return zero, err
	}
	if wbi {
		query, err = client.WBIValues(ctx, query)
		if err != nil {
			return zero, err
		}
	}
	request, err := core.NewQueryRequest(ctx, http.MethodGet, endpoint, query)
	if err != nil {
		return zero, err
	}
	return core.SendPayload[T](ctx, client, request, operation)
}

func noLiveQuery() (url.Values, error) { return nil, nil }

func optionalLiveCSRF(client *core.Client) string {
	account, ok := client.Account()
	if !ok {
		return ""
	}
	return account.BiliJCT
}

func setLiveRoomHeaders(request *http.Request, referer string) {
	request.Header.Set("Referer", referer)
	request.Header.Set("Origin", "https://live.bilibili.com")
}
