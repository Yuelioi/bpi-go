// Package historytoview contains private account-history and watch-later read
// parameters and response models.
package historytoview

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type Business string

const (
	BusinessArchive     Business = "archive"
	BusinessPGC         Business = "pgc"
	BusinessLive        Business = "live"
	BusinessArticleList Business = "article-list"
	BusinessArticle     Business = "article"
)

type ListType string

const (
	ListAll     ListType = "all"
	ListArchive ListType = "archive"
	ListLive    ListType = "live"
	ListArticle ListType = "article"
)

type ListParams struct {
	max      *uint64
	business string
	viewAt   *uint64
	typeName string
	pageSize uint32
}

func NewListParams() ListParams { return ListParams{} }

func (p ListParams) WithMax(max uint64) ListParams { p.max = &max; return p }
func (p ListParams) WithBusiness(business Business) ListParams {
	p.business = string(business)
	return p
}
func (p ListParams) WithViewAt(viewAt uint64) ListParams   { p.viewAt = &viewAt; return p }
func (p ListParams) WithType(listType ListType) ListParams { p.typeName = string(listType); return p }

func (p ListParams) WithPageSize(size uint32) (ListParams, error) {
	if size == 0 {
		return ListParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be non-zero"}
	}
	p.pageSize = size
	return p, nil
}

func (p ListParams) WithRawBusiness(value string) (ListParams, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return ListParams{}, &bpierr.ParameterError{Field: "business", Message: "value cannot be blank"}
	}
	p.business = value
	return p, nil
}

func (p ListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.max != nil {
		values.Set("max", strconv.FormatUint(*p.max, 10))
	}
	if p.business != "" {
		values.Set("business", p.business)
	}
	if p.viewAt != nil {
		values.Set("view_at", strconv.FormatUint(*p.viewAt, 10))
	}
	if p.typeName != "" {
		values.Set("type", p.typeName)
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	return values, nil
}

type Cursor struct {
	Max      uint64 `json:"max"`
	ViewAt   uint64 `json:"view_at"`
	Business string `json:"business"`
	PageSize uint32 `json:"ps"`
}

type Tab struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type HistoryDetail struct {
	OID       uint64  `json:"oid"`
	EpisodeID *uint64 `json:"epid"`
	BVID      *string `json:"bvid"`
	Page      *uint32 `json:"page"`
	CID       *uint64 `json:"cid"`
	Part      *string `json:"part"`
	Business  string  `json:"business"`
	Device    uint32  `json:"dt"`
}

type HistoryItem struct {
	Title      string        `json:"title"`
	History    HistoryDetail `json:"history"`
	ViewAt     uint64        `json:"view_at"`
	Progress   int32         `json:"progress"`
	Duration   *uint32       `json:"duration"`
	Total      *int32        `json:"total"`
	Favorite   uint8         `json:"is_fav"`
	KID        uint64        `json:"kid"`
	AuthorName *string       `json:"author_name"`
}

type HistoryList struct {
	Cursor Cursor        `json:"cursor"`
	Tabs   []Tab         `json:"tab"`
	Items  []HistoryItem `json:"list"`
}

type ToViewStat struct {
	AID      uint64 `json:"aid"`
	View     uint64 `json:"view"`
	Danmaku  uint64 `json:"danmaku"`
	Reply    uint64 `json:"reply"`
	Favorite uint64 `json:"favorite"`
	Coin     uint64 `json:"coin"`
	Share    uint64 `json:"share"`
	Like     uint64 `json:"like"`
	VT       int64  `json:"vt"`
	VV       int64  `json:"vv"`
}

type ToViewItem struct {
	AID      uint64     `json:"aid"`
	BVID     string     `json:"bvid"`
	Title    string     `json:"title"`
	Picture  string     `json:"pic"`
	Duration uint32     `json:"duration"`
	CID      uint64     `json:"cid"`
	Progress int32      `json:"progress"`
	AddedAt  uint64     `json:"add_at"`
	Stat     ToViewStat `json:"stat"`
}

type ToViewList struct {
	Count uint32       `json:"count"`
	Items []ToViewItem `json:"list"`
}
