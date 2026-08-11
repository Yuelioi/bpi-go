// Package manga contains validated parameters and stable response models for
// Bilibili Manga endpoints.
package manga

import "github.com/Yuelioi/bpi-go/internal/bpierr"

// CouponsParams selects one page of unexpired account coupons.
type CouponsParams struct {
	Page     uint32 `json:"pageNum"`
	PageSize uint32 `json:"pageSize"`
}

// NewCouponsParams creates coupon pagination with the Rust contract defaults.
func NewCouponsParams() CouponsParams {
	return CouponsParams{Page: 1, PageSize: 20}
}

func (p CouponsParams) WithPage(page uint32) (CouponsParams, error) {
	if page == 0 {
		return CouponsParams{}, &bpierr.ParameterError{Field: "pageNum", Message: "page must be at least 1"}
	}
	p.Page = page
	return p, nil
}

func (p CouponsParams) WithPageSize(size uint32) (CouponsParams, error) {
	if size == 0 || size > 100 {
		return CouponsParams{}, &bpierr.ParameterError{Field: "pageSize", Message: "page size must be between 1 and 100"}
	}
	p.PageSize = size
	return p, nil
}

func (p CouponsParams) Validate() error {
	if p.Page == 0 {
		return &bpierr.ParameterError{Field: "pageNum", Message: "page must be at least 1"}
	}
	if p.PageSize == 0 || p.PageSize > 100 {
		return &bpierr.ParameterError{Field: "pageSize", Message: "page size must be between 1 and 100"}
	}
	return nil
}

type couponsRequest struct {
	Page       uint32 `json:"pageNum"`
	PageSize   uint32 `json:"pageSize"`
	NotExpired bool   `json:"notExpired"`
	TabType    int32  `json:"tabType"`
	Type       int32  `json:"type"`
}

func (p CouponsParams) RequestBody() (any, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return couponsRequest{Page: p.Page, PageSize: p.PageSize, NotExpired: true, TabType: 1, Type: 0}, nil
}
