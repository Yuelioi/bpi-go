// Package fav contains parameters and stable response models for Bilibili
// favorite-folder read endpoints.
package fav

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type FolderInfoParams struct{ mediaID ids.MediaID }

func NewFolderInfoParams(mediaID ids.MediaID) FolderInfoParams {
	return FolderInfoParams{mediaID: mediaID}
}

func (p FolderInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.mediaID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"media_id": {p.mediaID.String()}}, nil
}

type CreatedListParams struct {
	upMID       ids.MID
	typeID      *uint8
	resourceID  *uint64
	webLocation string
}

func NewCreatedListParams(upMID ids.MID) CreatedListParams {
	return CreatedListParams{upMID: upMID, webLocation: "333.1387"}
}

func (p CreatedListParams) WithType(typeID uint8) CreatedListParams {
	p.typeID = &typeID
	return p
}

func (p CreatedListParams) WithResourceID(resourceID uint64) (CreatedListParams, error) {
	if resourceID == 0 {
		return CreatedListParams{}, parameterError("rid", "value must be non-zero")
	}
	p.resourceID = &resourceID
	return p, nil
}

func (p CreatedListParams) WithWebLocation(webLocation string) (CreatedListParams, error) {
	value, err := nonBlank("web_location", webLocation)
	if err != nil {
		return CreatedListParams{}, err
	}
	p.webLocation = value
	return p, nil
}

func (p CreatedListParams) EncodeQuery() (url.Values, error) {
	if err := p.upMID.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"up_mid": {p.upMID.String()}, "web_location": {p.webLocation}}
	if p.typeID != nil {
		values.Set("type", strconv.FormatUint(uint64(*p.typeID), 10))
	}
	if p.resourceID != nil {
		values.Set("rid", strconv.FormatUint(*p.resourceID, 10))
	}
	return values, nil
}

type CollectedListParams struct {
	upMID    ids.MID
	page     uint32
	pageSize uint32
	platform string
}

func NewCollectedListParams(upMID ids.MID) CollectedListParams {
	return CollectedListParams{upMID: upMID, page: 1, pageSize: 20, platform: "web"}
}

func (p CollectedListParams) WithPage(page uint32) (CollectedListParams, error) {
	if page == 0 {
		return CollectedListParams{}, parameterError("pn", "value must be non-zero")
	}
	p.page = page
	return p, nil
}

func (p CollectedListParams) WithPageSize(pageSize uint32) (CollectedListParams, error) {
	if pageSize == 0 {
		return CollectedListParams{}, parameterError("ps", "value must be non-zero")
	}
	p.pageSize = pageSize
	return p, nil
}

func (p CollectedListParams) EncodeQuery() (url.Values, error) {
	if err := p.upMID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{
		"up_mid":   {p.upMID.String()},
		"pn":       {strconv.FormatUint(uint64(p.page), 10)},
		"ps":       {strconv.FormatUint(uint64(p.pageSize), 10)},
		"platform": {p.platform},
	}, nil
}

type ResourceInfosParams struct {
	resources string
	platform  string
}

func NewResourceInfosParams(resources string) (ResourceInfosParams, error) {
	value, err := nonBlank("resources", resources)
	if err != nil {
		return ResourceInfosParams{}, err
	}
	return ResourceInfosParams{resources: value, platform: "web"}, nil
}

func (p ResourceInfosParams) EncodeQuery() (url.Values, error) {
	return url.Values{"resources": {p.resources}, "platform": {p.platform}}, nil
}

type ResourceIDsParams struct{ mediaID ids.MediaID }

func NewResourceIDsParams(mediaID ids.MediaID) ResourceIDsParams {
	return ResourceIDsParams{mediaID: mediaID}
}

func (p ResourceIDsParams) EncodeQuery() (url.Values, error) {
	if err := p.mediaID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"media_id": {p.mediaID.String()}, "platform": {"web"}}, nil
}

type ListDetailParams struct {
	mediaID  ids.MediaID
	tid      *uint32
	keyword  string
	order    string
	typeID   *uint8
	pageSize uint32
	page     *uint32
}

func NewListDetailParams(mediaID ids.MediaID) ListDetailParams {
	return ListDetailParams{mediaID: mediaID, pageSize: 20}
}

func (p ListDetailParams) WithTID(tid uint32) ListDetailParams { p.tid = &tid; return p }

func (p ListDetailParams) WithKeyword(keyword string) (ListDetailParams, error) {
	value, err := nonBlank("keyword", keyword)
	if err != nil {
		return ListDetailParams{}, err
	}
	p.keyword = value
	return p, nil
}

func (p ListDetailParams) WithOrder(order string) (ListDetailParams, error) {
	value, err := nonBlank("order", order)
	if err != nil {
		return ListDetailParams{}, err
	}
	p.order = value
	return p, nil
}

func (p ListDetailParams) WithContentType(typeID uint8) ListDetailParams {
	p.typeID = &typeID
	return p
}

func (p ListDetailParams) WithPageSize(pageSize uint32) (ListDetailParams, error) {
	if pageSize == 0 {
		return ListDetailParams{}, parameterError("ps", "page size must be non-zero")
	}
	p.pageSize = pageSize
	return p, nil
}

func (p ListDetailParams) WithPage(page uint32) (ListDetailParams, error) {
	if page == 0 {
		return ListDetailParams{}, parameterError("pn", "page number must be non-zero")
	}
	p.page = &page
	return p, nil
}

func (p ListDetailParams) EncodeQuery() (url.Values, error) {
	if err := p.mediaID.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{
		"media_id": {p.mediaID.String()},
		"platform": {"web"},
		"ps":       {strconv.FormatUint(uint64(p.pageSize), 10)},
	}
	if p.tid != nil {
		values.Set("tid", strconv.FormatUint(uint64(*p.tid), 10))
	}
	if p.keyword != "" {
		values.Set("keyword", p.keyword)
	}
	if p.order != "" {
		values.Set("order", p.order)
	}
	if p.typeID != nil {
		values.Set("type", strconv.FormatUint(uint64(*p.typeID), 10))
	}
	if p.page != nil {
		values.Set("pn", strconv.FormatUint(uint64(*p.page), 10))
	}
	return values, nil
}

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

type FolderUpper struct {
	MID       uint64 `json:"mid"`
	Name      string `json:"name"`
	Face      string `json:"face"`
	Followed  bool   `json:"followed"`
	VIPType   uint8  `json:"vip_type"`
	VIPStatus uint8  `json:"vip_statue"`
}

type CountInfo struct {
	Collect  uint64  `json:"collect"`
	Play     uint64  `json:"play"`
	ThumbUp  *uint64 `json:"thumb_up"`
	Share    *uint64 `json:"share"`
	Danmaku  *uint64 `json:"danmaku"`
	ViewText *string `json:"view_text_1"`
}

type FolderInfo struct {
	ID         uint64      `json:"id"`
	FID        uint64      `json:"fid"`
	MID        uint64      `json:"mid"`
	Attribute  uint32      `json:"attr"`
	Title      string      `json:"title"`
	Cover      string      `json:"cover"`
	Upper      FolderUpper `json:"upper"`
	CoverType  uint8       `json:"cover_type"`
	CountInfo  CountInfo   `json:"cnt_info"`
	Type       uint32      `json:"type"`
	Intro      string      `json:"intro"`
	CreatedAt  uint64      `json:"ctime"`
	ModifiedAt uint64      `json:"mtime"`
	State      uint8       `json:"state"`
	FavState   uint8       `json:"fav_state"`
	LikeState  uint8       `json:"like_state"`
	MediaCount uint32      `json:"media_count"`
}

type CreatedFolder struct {
	ID         uint64 `json:"id"`
	FID        uint64 `json:"fid"`
	MID        uint64 `json:"mid"`
	Attribute  uint32 `json:"attr"`
	Title      string `json:"title"`
	FavState   uint8  `json:"fav_state"`
	MediaCount uint32 `json:"media_count"`
}

type CreatedList struct {
	Count uint32          `json:"count"`
	List  []CreatedFolder `json:"list"`
}

type CollectedUpper struct {
	MID  uint64 `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

type CollectedFolder struct {
	ID         uint64         `json:"id"`
	FID        uint64         `json:"fid"`
	MID        uint64         `json:"mid"`
	Attribute  uint32         `json:"attr"`
	Title      string         `json:"title"`
	Cover      string         `json:"cover"`
	Upper      CollectedUpper `json:"upper"`
	CoverType  uint8          `json:"cover_type"`
	Intro      string         `json:"intro"`
	CreatedAt  uint64         `json:"ctime"`
	ModifiedAt uint64         `json:"mtime"`
	State      uint8          `json:"state"`
	FavState   uint8          `json:"fav_state"`
	MediaCount uint32         `json:"media_count"`
}

type CollectedList struct {
	Count uint32            `json:"count"`
	List  []CollectedFolder `json:"list"`
}

type ResourceUpper struct {
	MID  uint64 `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

type ResourceInfo struct {
	ID          uint64          `json:"id"`
	Type        uint8           `json:"type"`
	Title       string          `json:"title"`
	Cover       string          `json:"cover"`
	Intro       string          `json:"intro"`
	Page        *uint32         `json:"page"`
	Duration    uint32          `json:"duration"`
	Upper       ResourceUpper   `json:"upper"`
	Attribute   uint8           `json:"attr"`
	CountInfo   CountInfo       `json:"cnt_info"`
	Link        string          `json:"link"`
	CreatedAt   uint64          `json:"ctime"`
	PublishedAt uint64          `json:"pubtime"`
	FavoritedAt uint64          `json:"fav_time"`
	BVID        *string         `json:"bvid"`
	BVIDLegacy  *string         `json:"bv_id"`
	Season      json.RawMessage `json:"season"`
}

type ListDetail struct {
	Info    FolderInfo     `json:"info"`
	Medias  []ResourceInfo `json:"medias"`
	HasMore bool           `json:"has_more"`
	TTL     uint64         `json:"ttl"`
}

type ResourceID struct {
	ID         uint64  `json:"id"`
	Type       uint8   `json:"type"`
	BVID       *string `json:"bvid"`
	BVIDLegacy *string `json:"bv_id"`
}
