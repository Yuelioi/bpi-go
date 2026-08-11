// Package videoranking contains validated parameters and response models for
// Bilibili video-ranking endpoints.
package videoranking

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type pagination struct {
	page     uint32
	pageSize uint32
}

func (p pagination) values() url.Values {
	values := url.Values{}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	return values
}

func validatePositive(field string, value uint64) error {
	if value == 0 {
		return &bpierr.ParameterError{Field: field, Message: "value must be non-zero"}
	}
	return nil
}

type PopularListParams struct{ pagination pagination }

func NewPopularListParams() PopularListParams { return PopularListParams{} }
func (p PopularListParams) WithPage(page uint32) (PopularListParams, error) {
	if err := validatePositive("pn", uint64(page)); err != nil {
		return PopularListParams{}, err
	}
	p.pagination.page = page
	return p, nil
}
func (p PopularListParams) WithPageSize(size uint32) (PopularListParams, error) {
	if err := validatePositive("ps", uint64(size)); err != nil {
		return PopularListParams{}, err
	}
	p.pagination.pageSize = size
	return p, nil
}
func (p PopularListParams) EncodeQuery() (url.Values, error) { return p.pagination.values(), nil }

type PopularSeriesParams struct{ number uint32 }

func NewPopularSeriesParams(number uint32) (PopularSeriesParams, error) {
	if err := validatePositive("number", uint64(number)); err != nil {
		return PopularSeriesParams{}, err
	}
	return PopularSeriesParams{number: number}, nil
}
func (p PopularSeriesParams) EncodeQuery() (url.Values, error) {
	if err := validatePositive("number", uint64(p.number)); err != nil {
		return nil, err
	}
	return url.Values{"number": {strconv.FormatUint(uint64(p.number), 10)}}, nil
}

type RankingType string

const (
	RankingAll    RankingType = "all"
	RankingRookie RankingType = "rookie"
	RankingOrigin RankingType = "origin"
)

func (typ RankingType) validate() error {
	switch typ {
	case RankingAll, RankingRookie, RankingOrigin:
		return nil
	default:
		return &bpierr.ParameterError{Field: "type", Message: "value must be all, rookie, or origin"}
	}
}

type RankingListParams struct {
	regionID uint32
	typeID   RankingType
}

func NewRankingListParams() RankingListParams { return RankingListParams{} }
func (p RankingListParams) WithRegionID(regionID uint32) (RankingListParams, error) {
	if err := validatePositive("rid", uint64(regionID)); err != nil {
		return RankingListParams{}, err
	}
	p.regionID = regionID
	return p, nil
}
func (p RankingListParams) WithType(typ RankingType) (RankingListParams, error) {
	if err := typ.validate(); err != nil {
		return RankingListParams{}, err
	}
	p.typeID = typ
	return p, nil
}
func (p RankingListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.regionID != 0 {
		values.Set("rid", strconv.FormatUint(uint64(p.regionID), 10))
	}
	if p.typeID != "" {
		if err := p.typeID.validate(); err != nil {
			return nil, err
		}
		values.Set("type", string(p.typeID))
	}
	return values, nil
}

type RegionDynamicParams struct {
	regionID   uint32
	pagination pagination
}

func NewRegionDynamicParams(regionID uint32) (RegionDynamicParams, error) {
	if err := validatePositive("rid", uint64(regionID)); err != nil {
		return RegionDynamicParams{}, err
	}
	return RegionDynamicParams{regionID: regionID}, nil
}
func (p RegionDynamicParams) WithPage(page uint32) (RegionDynamicParams, error) {
	if err := validatePositive("pn", uint64(page)); err != nil {
		return RegionDynamicParams{}, err
	}
	p.pagination.page = page
	return p, nil
}
func (p RegionDynamicParams) WithPageSize(size uint32) (RegionDynamicParams, error) {
	if err := validatePositive("ps", uint64(size)); err != nil {
		return RegionDynamicParams{}, err
	}
	p.pagination.pageSize = size
	return p, nil
}
func (p RegionDynamicParams) EncodeQuery() (url.Values, error) {
	if err := validatePositive("rid", uint64(p.regionID)); err != nil {
		return nil, err
	}
	values := p.pagination.values()
	values.Set("rid", strconv.FormatUint(uint64(p.regionID), 10))
	return values, nil
}

type RegionTagDynamicParams struct {
	regionID   uint32
	tagID      uint64
	pagination pagination
}

func NewRegionTagDynamicParams(regionID uint32, tagID uint64) (RegionTagDynamicParams, error) {
	if err := validatePositive("rid", uint64(regionID)); err != nil {
		return RegionTagDynamicParams{}, err
	}
	if err := validatePositive("tag_id", tagID); err != nil {
		return RegionTagDynamicParams{}, err
	}
	return RegionTagDynamicParams{regionID: regionID, tagID: tagID}, nil
}
func (p RegionTagDynamicParams) WithPage(page uint32) (RegionTagDynamicParams, error) {
	if err := validatePositive("pn", uint64(page)); err != nil {
		return RegionTagDynamicParams{}, err
	}
	p.pagination.page = page
	return p, nil
}
func (p RegionTagDynamicParams) WithPageSize(size uint32) (RegionTagDynamicParams, error) {
	if err := validatePositive("ps", uint64(size)); err != nil {
		return RegionTagDynamicParams{}, err
	}
	p.pagination.pageSize = size
	return p, nil
}
func (p RegionTagDynamicParams) EncodeQuery() (url.Values, error) {
	if err := validatePositive("rid", uint64(p.regionID)); err != nil {
		return nil, err
	}
	if err := validatePositive("tag_id", p.tagID); err != nil {
		return nil, err
	}
	values := p.pagination.values()
	values.Set("rid", strconv.FormatUint(uint64(p.regionID), 10))
	values.Set("tag_id", strconv.FormatUint(p.tagID, 10))
	return values, nil
}

type RegionNewListParams struct {
	regionID   uint32
	typeID     uint32
	pagination pagination
}

func NewRegionNewListParams(regionID uint32) (RegionNewListParams, error) {
	if err := validatePositive("rid", uint64(regionID)); err != nil {
		return RegionNewListParams{}, err
	}
	return RegionNewListParams{regionID: regionID}, nil
}
func (p RegionNewListParams) WithPage(page uint32) (RegionNewListParams, error) {
	if err := validatePositive("pn", uint64(page)); err != nil {
		return RegionNewListParams{}, err
	}
	p.pagination.page = page
	return p, nil
}
func (p RegionNewListParams) WithPageSize(size uint32) (RegionNewListParams, error) {
	if err := validatePositive("ps", uint64(size)); err != nil {
		return RegionNewListParams{}, err
	}
	p.pagination.pageSize = size
	return p, nil
}
func (p RegionNewListParams) WithType(typ uint32) (RegionNewListParams, error) {
	if err := validatePositive("type", uint64(typ)); err != nil {
		return RegionNewListParams{}, err
	}
	p.typeID = typ
	return p, nil
}
func (p RegionNewListParams) EncodeQuery() (url.Values, error) {
	if err := validatePositive("rid", uint64(p.regionID)); err != nil {
		return nil, err
	}
	values := p.pagination.values()
	values.Set("rid", strconv.FormatUint(uint64(p.regionID), 10))
	if p.typeID != 0 {
		values.Set("type", strconv.FormatUint(uint64(p.typeID), 10))
	}
	return values, nil
}

type NewListRankOrder string

const (
	NewListRankClick   NewListRankOrder = "click"
	NewListRankScores  NewListRankOrder = "scores"
	NewListRankPublish NewListRankOrder = "pubdate"
)

func (order NewListRankOrder) validate() error {
	switch order {
	case NewListRankClick, NewListRankScores, NewListRankPublish:
		return nil
	default:
		return &bpierr.ParameterError{Field: "order", Message: "value must be click, scores, or pubdate"}
	}
}

type RegionNewListRankParams struct {
	categoryID uint32
	order      NewListRankOrder
	page       uint32
	pageSize   uint32
	from       string
	to         string
}

func NewRegionNewListRankParams(categoryID, pageSize uint32, from, to string) (RegionNewListRankParams, error) {
	if err := validatePositive("cate_id", uint64(categoryID)); err != nil {
		return RegionNewListRankParams{}, err
	}
	if err := validatePositive("pagesize", uint64(pageSize)); err != nil {
		return RegionNewListRankParams{}, err
	}
	if strings.TrimSpace(from) == "" {
		return RegionNewListRankParams{}, &bpierr.ParameterError{Field: "time_from", Message: "value cannot be blank"}
	}
	if strings.TrimSpace(to) == "" {
		return RegionNewListRankParams{}, &bpierr.ParameterError{Field: "time_to", Message: "value cannot be blank"}
	}
	return RegionNewListRankParams{categoryID: categoryID, pageSize: pageSize, from: from, to: to}, nil
}
func (p RegionNewListRankParams) WithOrder(order NewListRankOrder) (RegionNewListRankParams, error) {
	if err := order.validate(); err != nil {
		return RegionNewListRankParams{}, err
	}
	p.order = order
	return p, nil
}
func (p RegionNewListRankParams) WithPage(page uint32) (RegionNewListRankParams, error) {
	if err := validatePositive("page", uint64(page)); err != nil {
		return RegionNewListRankParams{}, err
	}
	p.page = page
	return p, nil
}
func (p RegionNewListRankParams) EncodeQuery() (url.Values, error) {
	if err := validatePositive("cate_id", uint64(p.categoryID)); err != nil {
		return nil, err
	}
	if err := validatePositive("pagesize", uint64(p.pageSize)); err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.from) == "" {
		return nil, &bpierr.ParameterError{Field: "time_from", Message: "value cannot be blank"}
	}
	if strings.TrimSpace(p.to) == "" {
		return nil, &bpierr.ParameterError{Field: "time_to", Message: "value cannot be blank"}
	}
	values := url.Values{
		"search_type": {"video"}, "view_type": {"hot_rank"},
		"cate_id":   {strconv.FormatUint(uint64(p.categoryID), 10)},
		"pagesize":  {strconv.FormatUint(uint64(p.pageSize), 10)},
		"time_from": {p.from}, "time_to": {p.to},
	}
	if p.order != "" {
		if err := p.order.validate(); err != nil {
			return nil, err
		}
		values.Set("order", string(p.order))
	}
	if p.page != 0 {
		values.Set("page", strconv.FormatUint(uint64(p.page), 10))
	}
	return values, nil
}

type PopularList struct {
	Items  []json.RawMessage `json:"list"`
	NoMore bool              `json:"no_more"`
}

type PopularSeriesItem struct {
	Number  uint32 `json:"number"`
	Subject string `json:"subject"`
	Status  uint8  `json:"status"`
	Name    string `json:"name"`
}

type PopularSeriesList struct {
	Items []PopularSeriesItem `json:"list"`
}

type PopularSeriesConfig struct {
	ID            uint64 `json:"id"`
	Type          string `json:"type"`
	Number        uint32 `json:"number"`
	Subject       string `json:"subject"`
	StartTime     uint64 `json:"stime"`
	EndTime       uint64 `json:"etime"`
	Status        uint8  `json:"status"`
	Name          string `json:"name"`
	Label         string `json:"label"`
	Hint          string `json:"hint"`
	Color         uint32 `json:"color"`
	Cover         string `json:"cover"`
	ShareTitle    string `json:"share_title"`
	ShareSubtitle string `json:"share_subtitle"`
	MediaID       uint64 `json:"media_id"`
}

type PopularSeries struct {
	Config   PopularSeriesConfig `json:"config"`
	Reminder *string             `json:"reminder"`
	Items    []json.RawMessage   `json:"list"`
}

type PreciousVideos struct {
	Title   string            `json:"title"`
	MediaID uint64            `json:"media_id"`
	Explain string            `json:"explain"`
	Items   []json.RawMessage `json:"list"`
}

type RankingList struct {
	Note  string            `json:"note"`
	Items []json.RawMessage `json:"list"`
}

type RegionPage struct {
	Count uint32 `json:"count"`
	Page  uint32 `json:"num"`
	Size  uint32 `json:"size"`
}

type RegionArchives struct {
	Archives []json.RawMessage `json:"archives"`
	Page     RegionPage        `json:"page"`
}

type NewListRankItem struct {
	PublishDate  string  `json:"pubdate"`
	Picture      string  `json:"pic"`
	Tag          string  `json:"tag"`
	Duration     uint32  `json:"duration"`
	ID           uint64  `json:"id"`
	RankScore    *uint64 `json:"rank_score"`
	BadgePay     bool    `json:"badgepay"`
	SendDate     *uint64 `json:"senddate"`
	Author       string  `json:"author"`
	Review       uint64  `json:"review"`
	MID          uint64  `json:"mid"`
	IsUnionVideo uint8   `json:"is_union_video"`
	RankIndex    *uint64 `json:"rank_index"`
	Type         string  `json:"type"`
	Play         string  `json:"play"`
	Danmaku      uint64  `json:"video_review"`
	IsPay        uint8   `json:"is_pay"`
	Favorites    uint64  `json:"favorites"`
	ArchiveURL   string  `json:"arcurl"`
	BVID         string  `json:"bvid"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
}

type NewListRank struct {
	Result     *[]NewListRankItem `json:"result"`
	NumResults uint32             `json:"numResults"`
	Page       uint32             `json:"page"`
	PageSize   uint32             `json:"pagesize"`
	Message    string             `json:"msg"`
}
