package comment

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type Sort uint8

const (
	SortByTime Sort = iota
	SortByLike
	SortByReplies
)

type ListParams struct {
	target   Target
	page     uint32
	pageSize uint32
	sort     *Sort
	noHot    *bool
}

func NewListParams(target Target) ListParams { return ListParams{target: target} }

func (p ListParams) WithPage(page uint32) (ListParams, error) {
	if page == 0 {
		return ListParams{}, &bpierr.ParameterError{Field: "pn", Message: "page must be greater than zero"}
	}
	p.page = page
	return p, nil
}

func (p ListParams) WithPageSize(size uint32) (ListParams, error) {
	if size == 0 || size > 20 {
		return ListParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be between 1 and 20"}
	}
	p.pageSize = size
	return p, nil
}

func (p ListParams) WithSort(sort Sort) (ListParams, error) {
	if sort > SortByReplies {
		return ListParams{}, &bpierr.ParameterError{Field: "sort", Message: "sort is not supported"}
	}
	p.sort = &sort
	return p, nil
}

func (p ListParams) WithoutHot(noHot bool) ListParams {
	p.noHot = &noHot
	return p
}

func (p ListParams) EncodeQuery() (url.Values, error) {
	values, err := targetQuery(p.target)
	if err != nil {
		return nil, err
	}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	if p.sort != nil {
		if *p.sort > SortByReplies {
			return nil, &bpierr.ParameterError{Field: "sort", Message: "sort is not supported"}
		}
		values.Set("sort", strconv.FormatUint(uint64(*p.sort), 10))
	}
	if p.noHot != nil {
		if *p.noHot {
			values.Set("nohot", "1")
		} else {
			values.Set("nohot", "0")
		}
	}
	return values, nil
}

type RepliesParams struct {
	target   Target
	root     int64
	page     uint32
	pageSize uint32
}

func NewRepliesParams(target Target, root int64) (RepliesParams, error) {
	if err := target.validate(); err != nil {
		return RepliesParams{}, err
	}
	if root <= 0 {
		return RepliesParams{}, &bpierr.ParameterError{Field: "root", Message: "root reply ID must be greater than zero"}
	}
	return RepliesParams{target: target, root: root}, nil
}

func (p RepliesParams) WithPage(page uint32) (RepliesParams, error) {
	if page == 0 {
		return RepliesParams{}, &bpierr.ParameterError{Field: "pn", Message: "page must be greater than zero"}
	}
	p.page = page
	return p, nil
}

func (p RepliesParams) WithPageSize(size uint32) (RepliesParams, error) {
	if size == 0 {
		return RepliesParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be greater than zero"}
	}
	p.pageSize = size
	return p, nil
}

func (p RepliesParams) EncodeQuery() (url.Values, error) {
	values, err := targetQuery(p.target)
	if err != nil {
		return nil, err
	}
	if p.root <= 0 {
		return nil, &bpierr.ParameterError{Field: "root", Message: "root reply ID must be greater than zero"}
	}
	values.Set("root", strconv.FormatInt(p.root, 10))
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	return values, nil
}

type CountParams struct{ target Target }

func NewCountParams(target Target) CountParams { return CountParams{target: target} }

func (p CountParams) EncodeQuery() (url.Values, error) { return targetQuery(p.target) }

func targetQuery(target Target) (url.Values, error) {
	if err := target.validate(); err != nil {
		return nil, err
	}
	return url.Values{
		"type": {strconv.FormatInt(int64(target.typeID), 10)},
		"oid":  {strconv.FormatInt(target.oid, 10)},
	}, nil
}

type Count struct {
	Count int64 `json:"count"`
}

type List struct {
	Page       *Page           `json:"page"`
	Replies    []Reply         `json:"replies"`
	TopReplies []Reply         `json:"top_replies"`
	Root       *Reply          `json:"root"`
	Config     json.RawMessage `json:"config"`
	Control    json.RawMessage `json:"control"`
	Cursor     json.RawMessage `json:"cursor"`
	Top        json.RawMessage `json:"top"`
	Upper      json.RawMessage `json:"upper"`
}

type Page struct {
	Page     uint64  `json:"num"`
	Size     uint64  `json:"size"`
	Count    uint64  `json:"count"`
	AllCount *uint64 `json:"acount"`
}

type Reply struct {
	RPID       int64   `json:"rpid"`
	OID        int64   `json:"oid"`
	Type       int64   `json:"type"`
	MID        int64   `json:"mid"`
	Root       int64   `json:"root"`
	Parent     int64   `json:"parent"`
	Count      int64   `json:"count"`
	ReplyCount int64   `json:"rcount"`
	CreatedAt  int64   `json:"ctime"`
	Likes      int64   `json:"like"`
	Member     Member  `json:"member"`
	Content    Content `json:"content"`
	Replies    []Reply `json:"replies"`
}

type Member struct {
	MID    string `json:"mid"`
	Name   string `json:"uname"`
	Sex    string `json:"sex"`
	Sign   string `json:"sign"`
	Avatar string `json:"avatar"`
}

type Content struct {
	Message string `json:"message"`
}
