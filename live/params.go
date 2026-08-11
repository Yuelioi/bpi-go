// Package live contains parameters and stable response models for Bilibili's
// promoted live read APIs.
package live

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// RoomInfoParams identifies a live room whose public metadata is requested.
type RoomInfoParams struct{ roomID ids.RoomID }

func NewRoomInfoParams(roomID ids.RoomID) RoomInfoParams { return RoomInfoParams{roomID: roomID} }

func (p RoomInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"room_id": {p.roomID.String()}}, nil
}

// StreamParams configures the legacy live stream URL read.
type StreamParams struct {
	roomID   ids.RoomID
	platform string
	quality  *uint32
	qn       *uint32
}

func NewStreamParams(roomID ids.RoomID) StreamParams { return StreamParams{roomID: roomID} }

func (p StreamParams) WithPlatform(platform string) (StreamParams, error) {
	value, err := nonBlank("platform", platform)
	if err != nil {
		return StreamParams{}, err
	}
	p.platform = value
	return p, nil
}

func (p StreamParams) WithQuality(quality uint32) (StreamParams, error) {
	if quality == 0 {
		return StreamParams{}, parameterError("quality", "value must be non-zero")
	}
	p.quality = &quality
	return p, nil
}

func (p StreamParams) WithQN(qn uint32) (StreamParams, error) {
	if qn == 0 {
		return StreamParams{}, parameterError("qn", "value must be non-zero")
	}
	p.qn = &qn
	return p, nil
}

func (p StreamParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"cid": {p.roomID.String()}}
	if p.platform != "" {
		values.Set("platform", p.platform)
	}
	if p.quality != nil {
		values.Set("quality", uint32String(*p.quality))
	}
	if p.qn != nil {
		values.Set("qn", uint32String(*p.qn))
	}
	return values, nil
}

// RoomGiftListParams configures a room gift-panel read.
type RoomGiftListParams struct {
	roomID       ids.RoomID
	areaParentID *uint32
	areaID       *uint32
}

func NewRoomGiftListParams(roomID ids.RoomID) RoomGiftListParams {
	return RoomGiftListParams{roomID: roomID}
}

func (p RoomGiftListParams) WithAreaParentID(areaParentID uint32) (RoomGiftListParams, error) {
	if areaParentID == 0 {
		return RoomGiftListParams{}, parameterError("area_parent_id", "value must be non-zero")
	}
	p.areaParentID = &areaParentID
	return p, nil
}

func (p RoomGiftListParams) WithAreaID(areaID uint32) (RoomGiftListParams, error) {
	if areaID == 0 {
		return RoomGiftListParams{}, parameterError("area_id", "value must be non-zero")
	}
	p.areaID = &areaID
	return p, nil
}

func (p RoomGiftListParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"room_id": {p.roomID.String()}, "platform": {"web"}}
	if p.areaParentID != nil {
		values.Set("area_parent_id", uint32String(*p.areaParentID))
	}
	if p.areaID != nil {
		values.Set("area_id", uint32String(*p.areaID))
	}
	return values, nil
}

// BlindGiftInfoParams identifies a blind-box gift.
type BlindGiftInfoParams struct{ giftID uint64 }

func NewBlindGiftInfoParams(giftID uint64) (BlindGiftInfoParams, error) {
	if giftID == 0 {
		return BlindGiftInfoParams{}, parameterError("gift_id", "value must be non-zero")
	}
	return BlindGiftInfoParams{giftID: giftID}, nil
}

func (p BlindGiftInfoParams) EncodeQuery() (url.Values, error) {
	if p.giftID == 0 {
		return nil, parameterError("gift_id", "value must be non-zero")
	}
	return url.Values{"gift_id": {strconv.FormatUint(p.giftID, 10)}}, nil
}

// DanmuInfoParams identifies a room and live danmaku connection type.
type DanmuInfoParams struct {
	roomID   ids.RoomID
	infoType uint8
}

func NewDanmuInfoParams(roomID ids.RoomID) DanmuInfoParams {
	return DanmuInfoParams{roomID: roomID}
}

func (p DanmuInfoParams) WithType(infoType uint8) DanmuInfoParams {
	p.infoType = infoType
	return p
}

func (p DanmuInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{
		"id":   {p.roomID.String()},
		"type": {strconv.FormatUint(uint64(p.infoType), 10)},
	}, nil
}

// EmoticonsParams configures a room emoticon read. The platform defaults to
// "pc", matching the promoted web contract.
type EmoticonsParams struct {
	roomID   ids.RoomID
	platform string
}

func NewEmoticonsParams(roomID ids.RoomID) EmoticonsParams {
	return EmoticonsParams{roomID: roomID, platform: "pc"}
}

func (p EmoticonsParams) WithPlatform(platform string) (EmoticonsParams, error) {
	value, err := nonBlank("platform", platform)
	if err != nil {
		return EmoticonsParams{}, err
	}
	p.platform = value
	return p, nil
}

func (p EmoticonsParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.platform) == "" {
		return nil, parameterError("platform", "value cannot be blank")
	}
	return url.Values{"room_id": {p.roomID.String()}, "platform": {p.platform}}, nil
}

// LotteryInfoParams identifies a live room whose lottery state is requested.
type LotteryInfoParams struct{ roomID ids.RoomID }

func NewLotteryInfoParams(roomID ids.RoomID) LotteryInfoParams {
	return LotteryInfoParams{roomID: roomID}
}

func (p LotteryInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"roomid": {p.roomID.String()}}, nil
}

// MyMedalsParams configures pagination for the current account's fan medals.
type MyMedalsParams struct{ page, pageSize uint32 }

func NewMyMedalsParams() MyMedalsParams { return MyMedalsParams{page: 1, pageSize: 10} }

func (p MyMedalsParams) WithPage(page uint32) (MyMedalsParams, error) {
	if page == 0 {
		return MyMedalsParams{}, parameterError("page", "value must be non-zero")
	}
	p.page = page
	return p, nil
}

func (p MyMedalsParams) WithPageSize(pageSize uint32) (MyMedalsParams, error) {
	if pageSize == 0 {
		return MyMedalsParams{}, parameterError("page_size", "value must be non-zero")
	}
	p.pageSize = pageSize
	return p, nil
}

func (p MyMedalsParams) EncodeQuery() (url.Values, error) {
	if p.page == 0 || p.pageSize == 0 {
		return nil, parameterError("pagination", "page and page_size must be non-zero")
	}
	return url.Values{"page": {uint32String(p.page)}, "page_size": {uint32String(p.pageSize)}}, nil
}

// FollowUpListParams configures the current account's followed-anchor list.
// Unset optional fields are omitted to preserve the Rust endpoint semantics.
type FollowUpListParams struct {
	page         *uint32
	pageSize     *uint32
	ignoreRecord *uint8
	hitAB        *bool
}

func NewFollowUpListParams() FollowUpListParams { return FollowUpListParams{} }

func (p FollowUpListParams) WithPage(page uint32) (FollowUpListParams, error) {
	if page == 0 {
		return FollowUpListParams{}, parameterError("page", "value must be non-zero")
	}
	p.page = &page
	return p, nil
}

func (p FollowUpListParams) WithPageSize(pageSize uint32) (FollowUpListParams, error) {
	if pageSize == 0 {
		return FollowUpListParams{}, parameterError("page_size", "value must be non-zero")
	}
	p.pageSize = &pageSize
	return p, nil
}

func (p FollowUpListParams) WithIgnoreRecord(ignoreRecord uint8) (FollowUpListParams, error) {
	if ignoreRecord > 1 {
		return FollowUpListParams{}, parameterError("ignoreRecord", "value must be 0 or 1")
	}
	p.ignoreRecord = &ignoreRecord
	return p, nil
}

func (p FollowUpListParams) WithHitAB(hitAB bool) FollowUpListParams {
	p.hitAB = &hitAB
	return p
}

func (p FollowUpListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.page != nil {
		values.Set("page", uint32String(*p.page))
	}
	if p.pageSize != nil {
		values.Set("page_size", uint32String(*p.pageSize))
	}
	if p.ignoreRecord != nil {
		values.Set("ignoreRecord", strconv.FormatUint(uint64(*p.ignoreRecord), 10))
	}
	if p.hitAB != nil {
		values.Set("hit_ab", strconv.FormatBool(*p.hitAB))
	}
	return values, nil
}

// FollowUpWebListParams configures the current live followed-anchor list.
type FollowUpWebListParams struct{ hitAB *bool }

func NewFollowUpWebListParams() FollowUpWebListParams { return FollowUpWebListParams{} }

func (p FollowUpWebListParams) WithHitAB(hitAB bool) FollowUpWebListParams {
	p.hitAB = &hitAB
	return p
}

func (p FollowUpWebListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.hitAB != nil {
		values.Set("hit_ab", strconv.FormatBool(*p.hitAB))
	}
	return values, nil
}

// ReplayListParams configures pagination for the current account's replays.
type ReplayListParams struct{ page, pageSize *uint32 }

func NewReplayListParams() ReplayListParams { return ReplayListParams{} }

func (p ReplayListParams) WithPage(page uint32) (ReplayListParams, error) {
	if page == 0 {
		return ReplayListParams{}, parameterError("page", "value must be non-zero")
	}
	p.page = &page
	return p, nil
}

func (p ReplayListParams) WithPageSize(pageSize uint32) (ReplayListParams, error) {
	if pageSize == 0 {
		return ReplayListParams{}, parameterError("page_size", "value must be non-zero")
	}
	p.pageSize = &pageSize
	return p, nil
}

func (p ReplayListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.page != nil {
		values.Set("page", uint32String(*p.page))
	}
	if p.pageSize != nil {
		values.Set("page_size", uint32String(*p.pageSize))
	}
	return values, nil
}

// GuardListParams configures a room's public guard-member list.
type GuardListParams struct {
	roomID   ids.RoomID
	anchorID ids.MID
	page     uint32
	pageSize uint32
	typeID   uint32
}

func NewGuardListParams(roomID ids.RoomID, anchorID ids.MID) GuardListParams {
	return GuardListParams{roomID: roomID, anchorID: anchorID, page: 1, pageSize: 20, typeID: 5}
}

func (p GuardListParams) WithPage(page uint32) (GuardListParams, error) {
	if page == 0 {
		return GuardListParams{}, parameterError("page", "value must be non-zero")
	}
	p.page = page
	return p, nil
}

func (p GuardListParams) WithPageSize(pageSize uint32) (GuardListParams, error) {
	if pageSize == 0 {
		return GuardListParams{}, parameterError("page_size", "value must be non-zero")
	}
	p.pageSize = pageSize
	return p, nil
}

func (p GuardListParams) WithType(typeID uint32) (GuardListParams, error) {
	if typeID == 0 {
		return GuardListParams{}, parameterError("typ", "value must be non-zero")
	}
	p.typeID = typeID
	return p, nil
}

func (p GuardListParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	if err := p.anchorID.Validate(); err != nil {
		return nil, err
	}
	if p.page == 0 || p.pageSize == 0 || p.typeID == 0 {
		return nil, parameterError("pagination", "page, page_size, and typ must be non-zero")
	}
	return url.Values{
		"roomid":    {p.roomID.String()},
		"ruid":      {p.anchorID.String()},
		"page":      {uint32String(p.page)},
		"page_size": {uint32String(p.pageSize)},
		"typ":       {uint32String(p.typeID)},
	}, nil
}

// SilentUsersParams configures a room moderation-list read.
type SilentUsersParams struct {
	roomID   ids.RoomID
	page     uint32
	pageSize uint32
}

func NewSilentUsersParams(roomID ids.RoomID) SilentUsersParams {
	return SilentUsersParams{roomID: roomID, page: 1, pageSize: 10}
}

func (p SilentUsersParams) WithPage(page uint32) (SilentUsersParams, error) {
	if page == 0 {
		return SilentUsersParams{}, parameterError("pn", "value must be non-zero")
	}
	p.page = page
	return p, nil
}

func (p SilentUsersParams) WithPageSize(pageSize uint32) (SilentUsersParams, error) {
	if pageSize == 0 {
		return SilentUsersParams{}, parameterError("ps", "value must be non-zero")
	}
	p.pageSize = pageSize
	return p, nil
}

func (p SilentUsersParams) EncodeForm(csrf string) (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	if p.page == 0 || p.pageSize == 0 {
		return nil, parameterError("pagination", "pn and ps must be non-zero")
	}
	return url.Values{
		"room_id":    {p.roomID.String()},
		"pn":         {uint32String(p.page)},
		"ps":         {uint32String(p.pageSize)},
		"csrf_token": {csrf},
		"csrf":       {csrf},
	}, nil
}

func (p SilentUsersParams) Referer() (string, error) { return roomReferer(p.roomID) }

// BannedUsersParams configures an anchor blacklist read. roomID is carried
// separately because the endpoint requires a room-specific Referer.
type BannedUsersParams struct {
	roomID   ids.RoomID
	anchorID ids.MID
	page     uint32
	pageSize uint32
}

func NewBannedUsersParams(roomID ids.RoomID, anchorID ids.MID) BannedUsersParams {
	return BannedUsersParams{roomID: roomID, anchorID: anchorID, page: 1, pageSize: 10}
}

func (p BannedUsersParams) WithPage(page uint32) (BannedUsersParams, error) {
	if page == 0 {
		return BannedUsersParams{}, parameterError("pn", "value must be non-zero")
	}
	p.page = page
	return p, nil
}

func (p BannedUsersParams) WithPageSize(pageSize uint32) (BannedUsersParams, error) {
	if pageSize == 0 {
		return BannedUsersParams{}, parameterError("ps", "value must be non-zero")
	}
	p.pageSize = pageSize
	return p, nil
}

func (p BannedUsersParams) EncodeQuery(csrf string) (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	if err := p.anchorID.Validate(); err != nil {
		return nil, err
	}
	if p.page == 0 || p.pageSize == 0 {
		return nil, parameterError("pagination", "pn and ps must be non-zero")
	}
	return url.Values{
		"anchor_id":  {p.anchorID.String()},
		"pn":         {uint32String(p.page)},
		"ps":         {uint32String(p.pageSize)},
		"mobi_app":   {"android"},
		"platform":   {"android"},
		"spmid":      {"444.8.0.0"},
		"csrf_token": {csrf},
		"csrf":       {csrf},
		"visit_id":   {""},
	}, nil
}

func (p BannedUsersParams) Referer() (string, error) { return roomReferer(p.roomID) }

// ShieldKeywordsParams identifies a room's private keyword list.
type ShieldKeywordsParams struct{ roomID ids.RoomID }

func NewShieldKeywordsParams(roomID ids.RoomID) ShieldKeywordsParams {
	return ShieldKeywordsParams{roomID: roomID}
}

func (p ShieldKeywordsParams) EncodeForm(csrf string) (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{
		"room_id":    {p.roomID.String()},
		"spmid":      {"444.8.0.0"},
		"csrf_token": {csrf},
		"csrf":       {csrf},
		"visit_id":   {""},
		"mobi_app":   {"android"},
		"platform":   {"android"},
	}, nil
}

func (p ShieldKeywordsParams) Referer() (string, error) { return roomReferer(p.roomID) }

// WebHeartBeatParams configures the public live telemetry heartbeat.
type WebHeartBeatParams struct {
	roomID       ids.RoomID
	nextInterval uint32
	platform     string
}

func NewWebHeartBeatParams(roomID ids.RoomID) WebHeartBeatParams {
	return WebHeartBeatParams{roomID: roomID, nextInterval: 60, platform: "web"}
}

func (p WebHeartBeatParams) WithNextInterval(nextInterval uint32) (WebHeartBeatParams, error) {
	if nextInterval == 0 {
		return WebHeartBeatParams{}, parameterError("next_interval", "value must be non-zero")
	}
	p.nextInterval = nextInterval
	return p, nil
}

func (p WebHeartBeatParams) WithPlatform(platform string) (WebHeartBeatParams, error) {
	value, err := nonBlank("platform", platform)
	if err != nil {
		return WebHeartBeatParams{}, err
	}
	p.platform = value
	return p, nil
}

func (p WebHeartBeatParams) EncodeQuery() (url.Values, error) {
	if err := p.roomID.Validate(); err != nil {
		return nil, err
	}
	if p.nextInterval == 0 {
		return nil, parameterError("next_interval", "value must be non-zero")
	}
	if strings.TrimSpace(p.platform) == "" {
		return nil, parameterError("platform", "value cannot be blank")
	}
	heartBeat := uint32String(p.nextInterval) + "|" + p.roomID.String() + "|1|0"
	return url.Values{
		"hb": {base64.StdEncoding.EncodeToString([]byte(heartBeat))},
		"pf": {p.platform},
	}, nil
}

func roomReferer(roomID ids.RoomID) (string, error) {
	if err := roomID.Validate(); err != nil {
		return "", err
	}
	return "https://live.bilibili.com/" + roomID.String(), nil
}

func uint32String(value uint32) string { return strconv.FormatUint(uint64(value), 10) }

func nonBlank(field, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", parameterError(field, "value cannot be blank")
	}
	return value, nil
}

func parameterError(field, message string) error {
	return &bpierr.ParameterError{Field: field, Message: message}
}
