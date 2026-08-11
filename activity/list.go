package activity

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

const (
	defaultPlatformFilter = "1,3"
	defaultHTTPMode       = uint32(3)
	defaultPage           = uint32(1)
	defaultPageSize       = uint32(15)
)

// ListPage is one page returned by activity.list.
type ListPage struct {
	Items []Item `json:"list"`
	Page  int32  `json:"num"`
	Size  int32  `json:"size"`
	Total int32  `json:"total"`
}

// Item is one stable activity-list entry.
type Item struct {
	ID          int32  `json:"id"`
	State       int32  `json:"state"`
	StartTime   int64  `json:"stime"`
	EndTime     int64  `json:"etime"`
	CreateTime  int64  `json:"ctime"`
	ModifyTime  int64  `json:"mtime"`
	Name        string `json:"name"`
	H5URL       string `json:"h5_url"`
	H5Cover     string `json:"h5_cover"`
	PageName    string `json:"page_name"`
	Platform    int32  `json:"plat"`
	Description string `json:"desc"`
}

// ListParams configures activity.list. Its zero value uses Bilibili's Web
// defaults: platform 1,3; mold 0; HTTP mode 3; page 1; page size 15.
type ListParams struct {
	platformFilter string
	mold           uint32
	httpMode       uint32
	httpModeSet    bool
	page           uint32
	pageSize       uint32
}

// NewListParams returns activity-list parameters with Web defaults.
func NewListParams() ListParams { return ListParams{} }

// WithPlatformFilter returns a copy with a non-blank platform filter.
func (p ListParams) WithPlatformFilter(filter string) (ListParams, error) {
	if strings.TrimSpace(filter) == "" {
		return ListParams{}, &bpierr.ParameterError{Field: "plat", Message: "platform filter cannot be blank"}
	}
	p.platformFilter = filter
	return p, nil
}

// WithMold returns a copy with the API mold flag.
func (p ListParams) WithMold(mold uint32) ListParams {
	p.mold = mold
	return p
}

// WithHTTPMode returns a copy with the API HTTP mode flag.
func (p ListParams) WithHTTPMode(mode uint32) ListParams {
	p.httpMode = mode
	p.httpModeSet = true
	return p
}

// WithPage returns a copy with a non-zero page number.
func (p ListParams) WithPage(page uint32) (ListParams, error) {
	if page == 0 {
		return ListParams{}, &bpierr.ParameterError{Field: "pn", Message: "page must be non-zero"}
	}
	p.page = page
	return p, nil
}

// WithPageSize returns a copy with a non-zero page size.
func (p ListParams) WithPageSize(size uint32) (ListParams, error) {
	if size == 0 {
		return ListParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be non-zero"}
	}
	p.pageSize = size
	return p, nil
}

// EncodeQuery returns a new URL query containing effective defaults and all
// configured values.
func (p ListParams) EncodeQuery() (url.Values, error) {
	platform := p.platformFilter
	if platform == "" {
		platform = defaultPlatformFilter
	}
	if strings.TrimSpace(platform) == "" {
		return nil, &bpierr.ParameterError{Field: "plat", Message: "platform filter cannot be blank"}
	}
	httpMode := defaultHTTPMode
	if p.httpModeSet {
		httpMode = p.httpMode
	}
	page := p.page
	if page == 0 {
		page = defaultPage
	}
	pageSize := p.pageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return url.Values{
		"plat": {platform},
		"mold": {strconv.FormatUint(uint64(p.mold), 10)},
		"http": {strconv.FormatUint(uint64(httpMode), 10)},
		"pn":   {strconv.FormatUint(uint64(page), 10)},
		"ps":   {strconv.FormatUint(uint64(pageSize), 10)},
	}, nil
}
