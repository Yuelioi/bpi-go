// Package user contains validated parameters and stable response models for
// Bilibili user-domain endpoints.
package user

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func validateMID(field string, mid ids.MID) error {
	if err := mid.Validate(); err != nil {
		return &bpierr.ParameterError{Field: field, Message: "member ID is invalid"}
	}
	return nil
}

func singleMIDQuery(field string, mid ids.MID) (url.Values, error) {
	if err := validateMID(field, mid); err != nil {
		return nil, err
	}
	return url.Values{field: {mid.String()}}, nil
}

// CardPhoto controls whether user.card includes the space header image.
type CardPhoto uint8

const (
	CardPhotoInclude CardPhoto = iota + 1
	CardPhotoExclude
)

// CardParams configures user.card.
type CardParams struct {
	mid   ids.MID
	photo CardPhoto
}

func NewCardParams(mid ids.MID) CardParams { return CardParams{mid: mid} }

func (p CardParams) WithPhoto(photo CardPhoto) CardParams {
	p.photo = photo
	return p
}

func (p CardParams) EncodeQuery() (url.Values, error) {
	values, err := singleMIDQuery("mid", p.mid)
	if err != nil {
		return nil, err
	}
	switch p.photo {
	case 0:
	case CardPhotoInclude:
		values.Set("photo", "true")
	case CardPhotoExclude:
		values.Set("photo", "false")
	default:
		return nil, &bpierr.ParameterError{Field: "photo", Message: "unsupported photo option"}
	}
	return values, nil
}

// CardsParams configures user.cards.
type CardsParams struct{ mids []ids.MID }

func NewCardsParams(mids ...ids.MID) (CardsParams, error) {
	validated, err := validateMIDs(mids)
	if err != nil {
		return CardsParams{}, err
	}
	return CardsParams{mids: validated}, nil
}

func (p CardsParams) EncodeQuery() (url.Values, error) {
	mids, err := encodeMIDs(p.mids)
	if err != nil {
		return nil, err
	}
	return url.Values{"uids": {mids}}, nil
}

// InfosParams configures user.infos.
type InfosParams struct{ mids []ids.MID }

func NewInfosParams(mids ...ids.MID) (InfosParams, error) {
	validated, err := validateMIDs(mids)
	if err != nil {
		return InfosParams{}, err
	}
	return InfosParams{mids: validated}, nil
}

func (p InfosParams) EncodeQuery() (url.Values, error) {
	mids, err := encodeMIDs(p.mids)
	if err != nil {
		return nil, err
	}
	return url.Values{"uids": {mids}}, nil
}

func validateMIDs(mids []ids.MID) ([]ids.MID, error) {
	if len(mids) == 0 {
		return nil, &bpierr.ParameterError{Field: "uids", Message: "at least one member ID is required"}
	}
	result := append([]ids.MID(nil), mids...)
	for _, mid := range result {
		if err := validateMID("uids", mid); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func encodeMIDs(mids []ids.MID) (string, error) {
	validated, err := validateMIDs(mids)
	if err != nil {
		return "", err
	}
	values := make([]string, len(validated))
	for index, mid := range validated {
		values[index] = mid.String()
	}
	return strings.Join(values, ","), nil
}

type SpaceParams struct{ mid ids.MID }

func NewSpaceParams(mid ids.MID) SpaceParams { return SpaceParams{mid: mid} }
func (p SpaceParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("mid", p.mid)
}

type SpaceNoticeParams struct{ mid ids.MID }

func NewSpaceNoticeParams(mid ids.MID) SpaceNoticeParams { return SpaceNoticeParams{mid: mid} }
func (p SpaceNoticeParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("mid", p.mid)
}

// BangumiFollowKind selects followed animation or cinema seasons.
type BangumiFollowKind uint8

const (
	BangumiFollow BangumiFollowKind = iota + 1
	CinemaFollow
)

type BangumiFollowListParams struct {
	mid      ids.MID
	kind     BangumiFollowKind
	page     uint32
	pageSize uint32
}

func NewBangumiFollowListParams(mid ids.MID) BangumiFollowListParams {
	return BangumiFollowListParams{mid: mid, kind: BangumiFollow, page: 1, pageSize: 15}
}

func (p BangumiFollowListParams) WithKind(kind BangumiFollowKind) BangumiFollowListParams {
	p.kind = kind
	return p
}

func (p BangumiFollowListParams) WithPage(page uint32) (BangumiFollowListParams, error) {
	if page == 0 {
		return BangumiFollowListParams{}, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	p.page = page
	return p, nil
}

func (p BangumiFollowListParams) WithPageSize(size uint32) (BangumiFollowListParams, error) {
	if size < 1 || size > 30 {
		return BangumiFollowListParams{}, &bpierr.ParameterError{Field: "page_size", Message: "page size must be between 1 and 30"}
	}
	p.pageSize = size
	return p, nil
}

func (p BangumiFollowListParams) EncodeQuery() (url.Values, error) {
	if err := validateMID("vmid", p.mid); err != nil {
		return nil, err
	}
	if p.page == 0 {
		return nil, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	if p.pageSize < 1 || p.pageSize > 30 {
		return nil, &bpierr.ParameterError{Field: "page_size", Message: "page size must be between 1 and 30"}
	}
	if p.kind != BangumiFollow && p.kind != CinemaFollow {
		return nil, &bpierr.ParameterError{Field: "kind", Message: "unsupported follow kind"}
	}
	return url.Values{
		"vmid": {p.mid.String()},
		"type": {strconv.FormatUint(uint64(p.kind), 10)},
		"pn":   {strconv.FormatUint(uint64(p.page), 10)},
		"ps":   {strconv.FormatUint(uint64(p.pageSize), 10)},
	}, nil
}

type RelationStatParams struct{ mid ids.MID }

func NewRelationStatParams(mid ids.MID) RelationStatParams { return RelationStatParams{mid: mid} }
func (p RelationStatParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("vmid", p.mid)
}

type FollowingsParams struct {
	mid       ids.MID
	orderType string
	pageSize  uint32
	page      uint32
}

func NewFollowingsParams(mid ids.MID) FollowingsParams { return FollowingsParams{mid: mid} }

func (p FollowingsParams) WithOrderType(orderType string) (FollowingsParams, error) {
	orderType = strings.TrimSpace(orderType)
	if orderType == "" {
		return FollowingsParams{}, &bpierr.ParameterError{Field: "order_type", Message: "order type cannot be blank"}
	}
	p.orderType = orderType
	return p, nil
}

func (p FollowingsParams) WithPageSize(size uint32) (FollowingsParams, error) {
	if size == 0 {
		return FollowingsParams{}, &bpierr.ParameterError{Field: "page_size", Message: "page size must be at least 1"}
	}
	p.pageSize = size
	return p, nil
}

func (p FollowingsParams) WithPage(page uint32) (FollowingsParams, error) {
	if page == 0 {
		return FollowingsParams{}, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	p.page = page
	return p, nil
}

func (p FollowingsParams) EncodeQuery() (url.Values, error) {
	values, err := singleMIDQuery("vmid", p.mid)
	if err != nil {
		return nil, err
	}
	if p.orderType != "" {
		values.Set("order_type", p.orderType)
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	return values, nil
}

type FollowersParams struct {
	mid          ids.MID
	pageSize     uint32
	page         uint32
	offset       string
	lastAccessTS uint64
	from         string
}

func NewFollowersParams(mid ids.MID) FollowersParams { return FollowersParams{mid: mid} }

func (p FollowersParams) WithPageSize(size uint32) (FollowersParams, error) {
	if size == 0 {
		return FollowersParams{}, &bpierr.ParameterError{Field: "page_size", Message: "page size must be at least 1"}
	}
	p.pageSize = size
	return p, nil
}

func (p FollowersParams) WithPage(page uint32) (FollowersParams, error) {
	if page == 0 {
		return FollowersParams{}, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	p.page = page
	return p, nil
}

func (p FollowersParams) WithOffset(offset string) (FollowersParams, error) {
	offset = strings.TrimSpace(offset)
	if offset == "" {
		return FollowersParams{}, &bpierr.ParameterError{Field: "offset", Message: "offset cannot be blank"}
	}
	p.offset = offset
	return p, nil
}

func (p FollowersParams) WithLastAccessTimestamp(timestamp uint64) FollowersParams {
	p.lastAccessTS = timestamp
	return p
}

func (p FollowersParams) WithFrom(from string) (FollowersParams, error) {
	from = strings.TrimSpace(from)
	if from == "" {
		return FollowersParams{}, &bpierr.ParameterError{Field: "from", Message: "source cannot be blank"}
	}
	p.from = from
	return p, nil
}

func (p FollowersParams) EncodeQuery() (url.Values, error) {
	values, err := singleMIDQuery("vmid", p.mid)
	if err != nil {
		return nil, err
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	if p.offset != "" {
		values.Set("offset", p.offset)
	}
	if p.lastAccessTS != 0 {
		values.Set("last_access_ts", strconv.FormatUint(p.lastAccessTS, 10))
	}
	if p.from != "" {
		values.Set("from", p.from)
	}
	return values, nil
}

type MedalWallParams struct{ targetID ids.MID }

func NewMedalWallParams(targetID ids.MID) MedalWallParams {
	return MedalWallParams{targetID: targetID}
}
func (p MedalWallParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("target_id", p.targetID)
}

type UpStatParams struct{ mid ids.MID }

func NewUpStatParams(mid ids.MID) UpStatParams { return UpStatParams{mid: mid} }
func (p UpStatParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("mid", p.mid)
}

type NavStatParams struct{ mid ids.MID }

func NewNavStatParams(mid ids.MID) NavStatParams { return NavStatParams{mid: mid} }
func (p NavStatParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("mid", p.mid)
}

type AlbumCountParams struct{ mid ids.MID }

func NewAlbumCountParams(mid ids.MID) AlbumCountParams { return AlbumCountParams{mid: mid} }
func (p AlbumCountParams) EncodeQuery() (url.Values, error) {
	return singleMIDQuery("uid", p.mid)
}

type NameToUIDParams struct{ names []string }

func NewNameToUIDParams(names ...string) (NameToUIDParams, error) {
	if len(names) == 0 {
		return NameToUIDParams{}, &bpierr.ParameterError{Field: "names", Message: "at least one name is required"}
	}
	result := make([]string, len(names))
	for index, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			return NameToUIDParams{}, &bpierr.ParameterError{Field: "names", Message: "names cannot contain blank values"}
		}
		result[index] = name
	}
	return NameToUIDParams{names: result}, nil
}

func (p NameToUIDParams) EncodeQuery() (url.Values, error) {
	if len(p.names) == 0 {
		return nil, &bpierr.ParameterError{Field: "names", Message: "at least one name is required"}
	}
	for _, name := range p.names {
		if strings.TrimSpace(name) == "" {
			return nil, &bpierr.ParameterError{Field: "names", Message: "names cannot contain blank values"}
		}
	}
	return url.Values{"names": {strings.Join(p.names, ",")}}, nil
}

type UploadedVideoOrder string

const (
	UploadedVideoPubdate UploadedVideoOrder = "pubdate"
	UploadedVideoClick   UploadedVideoOrder = "click"
	UploadedVideoStow    UploadedVideoOrder = "stow"
)

func ParseUploadedVideoOrder(value string) (UploadedVideoOrder, error) {
	order := UploadedVideoOrder(strings.TrimSpace(value))
	if err := order.validate(); err != nil {
		return "", err
	}
	return order, nil
}

func (order UploadedVideoOrder) validate() error {
	switch order {
	case UploadedVideoPubdate, UploadedVideoClick, UploadedVideoStow:
		return nil
	default:
		return &bpierr.ParameterError{Field: "order", Message: "order must be pubdate, click, or stow"}
	}
}

type UploadedVideosParams struct {
	mid      ids.MID
	order    UploadedVideoOrder
	tid      uint64
	keyword  string
	page     uint32
	pageSize uint32
}

func NewUploadedVideosParams(mid ids.MID) UploadedVideosParams {
	return UploadedVideosParams{mid: mid, order: UploadedVideoPubdate, page: 1, pageSize: 30}
}

func (p UploadedVideosParams) WithOrder(order UploadedVideoOrder) UploadedVideosParams {
	p.order = order
	return p
}

func (p UploadedVideosParams) WithTID(tid uint64) UploadedVideosParams {
	p.tid = tid
	return p
}

func (p UploadedVideosParams) WithKeyword(keyword string) UploadedVideosParams {
	p.keyword = strings.TrimSpace(keyword)
	return p
}

func (p UploadedVideosParams) WithPage(page uint32) (UploadedVideosParams, error) {
	if page == 0 {
		return UploadedVideosParams{}, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	p.page = page
	return p, nil
}

func (p UploadedVideosParams) WithPageSize(size uint32) (UploadedVideosParams, error) {
	if size == 0 {
		return UploadedVideosParams{}, &bpierr.ParameterError{Field: "page_size", Message: "page size must be at least 1"}
	}
	p.pageSize = size
	return p, nil
}

func (p UploadedVideosParams) EncodeQuery() (url.Values, error) {
	if err := validateMID("mid", p.mid); err != nil {
		return nil, err
	}
	if err := p.order.validate(); err != nil {
		return nil, err
	}
	if p.page == 0 {
		return nil, &bpierr.ParameterError{Field: "page", Message: "page number must be at least 1"}
	}
	if p.pageSize == 0 {
		return nil, &bpierr.ParameterError{Field: "page_size", Message: "page size must be at least 1"}
	}
	values := url.Values{
		"mid":   {p.mid.String()},
		"order": {string(p.order)},
		"tid":   {strconv.FormatUint(p.tid, 10)},
		"pn":    {strconv.FormatUint(uint64(p.page), 10)},
		"ps":    {strconv.FormatUint(uint64(p.pageSize), 10)},
	}
	if p.keyword != "" {
		values.Set("keyword", p.keyword)
	}
	return values, nil
}
