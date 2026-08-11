// Package comment contains validated parameters and response models for the
// Bilibili comment domain.
package comment

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// Target identifies one Bilibili comment area.
type Target struct {
	typeID int32
	oid    int64
}

// NewTarget creates a positive comment type/object pair.
func NewTarget(typeID int32, oid int64) (Target, error) {
	if typeID <= 0 {
		return Target{}, &bpierr.ParameterError{Field: "type", Message: "comment type must be greater than zero"}
	}
	if oid <= 0 {
		return Target{}, &bpierr.ParameterError{Field: "oid", Message: "comment object ID must be greater than zero"}
	}
	return Target{typeID: typeID, oid: oid}, nil
}

func (t Target) validate() error {
	if t.typeID <= 0 {
		return &bpierr.ParameterError{Field: "type", Message: "comment type must be greater than zero"}
	}
	if t.oid <= 0 {
		return &bpierr.ParameterError{Field: "oid", Message: "comment object ID must be greater than zero"}
	}
	return nil
}

// HotParams configures comment.read.hot.
type HotParams struct {
	target   Target
	root     int64
	page     uint32
	pageSize uint32
}

// NewHotParams creates hot-comment parameters for a positive root reply ID.
func NewHotParams(target Target, root int64) (HotParams, error) {
	if err := target.validate(); err != nil {
		return HotParams{}, err
	}
	if root <= 0 {
		return HotParams{}, &bpierr.ParameterError{Field: "root", Message: "root reply ID must be greater than zero"}
	}
	return HotParams{target: target, root: root}, nil
}

// WithPage returns a copy with a non-zero page.
func (p HotParams) WithPage(page uint32) (HotParams, error) {
	if page == 0 {
		return HotParams{}, &bpierr.ParameterError{Field: "pn", Message: "page must be greater than zero"}
	}
	p.page = page
	return p, nil
}

// WithPageSize returns a copy with a non-zero page size.
func (p HotParams) WithPageSize(size uint32) (HotParams, error) {
	if size == 0 {
		return HotParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be greater than zero"}
	}
	p.pageSize = size
	return p, nil
}

// EncodeQuery returns a new URL query containing the request parameters.
func (p HotParams) EncodeQuery() (url.Values, error) {
	if err := p.target.validate(); err != nil {
		return nil, err
	}
	if p.root <= 0 {
		return nil, &bpierr.ParameterError{Field: "root", Message: "root reply ID must be greater than zero"}
	}
	values := url.Values{
		"type": {strconv.FormatInt(int64(p.target.typeID), 10)},
		"oid":  {strconv.FormatInt(p.target.oid, 10)},
		"root": {strconv.FormatInt(p.root, 10)},
	}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	return values, nil
}

// Hot is a hot-comment payload when the endpoint returns one. Individual
// reply bodies remain raw until promoted fixtures establish a stable subset.
type Hot struct {
	Page    HotPage           `json:"page"`
	Replies []json.RawMessage `json:"replies"`
}

// HotPage contains hot-comment pagination counters.
type HotPage struct {
	AllCount int64 `json:"acount"`
	Count    int64 `json:"count"`
	Page     int32 `json:"num"`
	Size     int32 `json:"size"`
}
