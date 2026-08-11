// Package creativecenter contains parameters and stable response models for
// Bilibili's authenticated creator-center read APIs.
package creativecenter

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type SeasonListOrder string

const (
	SeasonOrderCreated SeasonListOrder = "ctime"
	SeasonOrderUpdated SeasonListOrder = "mtime"
)

type SortOrder string

const (
	SortAscending  SortOrder = "asc"
	SortDescending SortOrder = "desc"
)

type SeasonListParams struct {
	page, pageSize uint32
	order          SeasonListOrder
	sort           SortOrder
}

func NewSeasonListParams(page, pageSize uint32) (SeasonListParams, error) {
	if err := validatePage(page, pageSize); err != nil {
		return SeasonListParams{}, err
	}
	return SeasonListParams{page: page, pageSize: pageSize}, nil
}

func (p SeasonListParams) WithOrder(order SeasonListOrder) SeasonListParams {
	p.order = order
	return p
}
func (p SeasonListParams) WithSort(sort SortOrder) SeasonListParams { p.sort = sort; return p }

func (p SeasonListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{"pn": {uint32String(p.page)}, "ps": {uint32String(p.pageSize)}}
	if p.order != "" {
		values.Set("order", string(p.order))
	}
	if p.sort != "" {
		values.Set("sort", string(p.sort))
	}
	return values, nil
}

type SeasonInfoParams struct{ seasonID ids.SeasonID }

func NewSeasonInfoParams(seasonID ids.SeasonID) SeasonInfoParams {
	return SeasonInfoParams{seasonID: seasonID}
}
func (p SeasonInfoParams) EncodeQuery() (url.Values, error) {
	if err := p.seasonID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"id": {p.seasonID.String()}}, nil
}

type SeasonByAIDParams struct{ aid ids.AID }

func NewSeasonByAIDParams(aid ids.AID) SeasonByAIDParams { return SeasonByAIDParams{aid: aid} }
func (p SeasonByAIDParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"id": {p.aid.String()}}, nil
}

type SectionParams struct{ seasonID ids.SeasonID }

func NewSectionParams(seasonID ids.SeasonID) SectionParams { return SectionParams{seasonID: seasonID} }
func (p SectionParams) EncodeQuery() (url.Values, error) {
	if err := p.seasonID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"id": {p.seasonID.String()}}, nil
}

type ArchivesListParams struct{ page, pageSize uint32 }

func NewArchivesListParams(page, pageSize uint32) (ArchivesListParams, error) {
	if err := validatePage(page, pageSize); err != nil {
		return ArchivesListParams{}, err
	}
	return ArchivesListParams{page: page, pageSize: pageSize}, nil
}
func (p ArchivesListParams) EncodeQuery() (url.Values, error) {
	return url.Values{"pn": {uint32String(p.page)}, "ps": {uint32String(p.pageSize)}}, nil
}

type ArchiveVideosParams struct{ aid ids.AID }

func NewArchiveVideosParams(aid ids.AID) ArchiveVideosParams { return ArchiveVideosParams{aid: aid} }
func (p ArchiveVideosParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"aid": {p.aid.String()}}, nil
}

type ArchiveCompareParams struct {
	timestamp *uint64
	size      *uint32
}

func NewArchiveCompareParams() ArchiveCompareParams { return ArchiveCompareParams{} }
func (p ArchiveCompareParams) WithTimestamp(timestamp uint64) (ArchiveCompareParams, error) {
	if timestamp == 0 {
		return ArchiveCompareParams{}, parameterError("t", "value must be non-zero")
	}
	p.timestamp = &timestamp
	return p, nil
}
func (p ArchiveCompareParams) WithSize(size uint32) (ArchiveCompareParams, error) {
	if size == 0 {
		return ArchiveCompareParams{}, parameterError("size", "value must be non-zero")
	}
	p.size = &size
	return p, nil
}
func (p ArchiveCompareParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.timestamp != nil {
		values.Set("t", strconv.FormatUint(*p.timestamp, 10))
	}
	if p.size != nil {
		values.Set("size", uint32String(*p.size))
	}
	return values, nil
}

type VideoTrendMetric uint8

const (
	VideoTrendPlay VideoTrendMetric = iota + 1
	VideoTrendDanmaku
	VideoTrendReply
	VideoTrendShare
	VideoTrendCoin
	VideoTrendFavorite
	VideoTrendCharge
	VideoTrendLike
)

type VideoTrendParams struct{ metric VideoTrendMetric }

func NewVideoTrendParams(metric VideoTrendMetric) VideoTrendParams {
	return VideoTrendParams{metric: metric}
}
func (p VideoTrendParams) EncodeQuery() (url.Values, error) {
	if p.metric < VideoTrendPlay || p.metric > VideoTrendLike {
		return nil, parameterError("type", "unsupported video trend metric")
	}
	return url.Values{"type": {strconv.FormatUint(uint64(p.metric), 10)}}, nil
}

type ArticleTrendMetric uint8

const (
	ArticleTrendRead ArticleTrendMetric = iota + 1
	ArticleTrendReply
	ArticleTrendShare
	ArticleTrendCoin
	ArticleTrendFavorite
	ArticleTrendLike
)

type ArticleTrendParams struct{ metric ArticleTrendMetric }

func NewArticleTrendParams(metric ArticleTrendMetric) ArticleTrendParams {
	return ArticleTrendParams{metric: metric}
}
func (p ArticleTrendParams) EncodeQuery() (url.Values, error) {
	if p.metric < ArticleTrendRead || p.metric > ArticleTrendLike {
		return nil, parameterError("type", "unsupported article trend metric")
	}
	return url.Values{"type": {strconv.FormatUint(uint64(p.metric), 10)}}, nil
}

func validatePage(page, pageSize uint32) error {
	if page == 0 {
		return parameterError("pn", "value must be non-zero")
	}
	if pageSize == 0 {
		return parameterError("ps", "value must be non-zero")
	}
	return nil
}
func uint32String(value uint32) string { return strconv.FormatUint(uint64(value), 10) }
func parameterError(field, message string) error {
	return &bpierr.ParameterError{Field: field, Message: message}
}
