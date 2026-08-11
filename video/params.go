// Package video contains validated parameters and response models for the
// Bilibili video domain.
package video

import (
	"net/url"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type identifier struct {
	aid  ids.AID
	bvid ids.BVID
}

func byAID(aid ids.AID) identifier    { return identifier{aid: aid} }
func byBVID(bvid ids.BVID) identifier { return identifier{bvid: bvid} }

func (id identifier) encode() (url.Values, error) {
	switch {
	case id.aid != 0 && id.bvid == "":
		if err := id.aid.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "aid", Message: "video ID is invalid"}
		}
		return url.Values{"aid": {id.aid.String()}}, nil
	case id.bvid != "" && id.aid == 0:
		if err := id.bvid.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "bvid", Message: "video ID is invalid"}
		}
		return url.Values{"bvid": {id.bvid.String()}}, nil
	default:
		return nil, &bpierr.ParameterError{Field: "video_id", Message: "exactly one AID or BVID is required"}
	}
}

// ViewParams configures video.view.
type ViewParams struct{ id identifier }

// ViewByAID identifies video.view by a validated AV ID.
func ViewByAID(aid ids.AID) ViewParams { return ViewParams{id: byAID(aid)} }

// ViewByBVID identifies video.view by a validated BV ID.
func ViewByBVID(bvid ids.BVID) ViewParams { return ViewParams{id: byBVID(bvid)} }

// EncodeQuery returns a new URL query containing the request parameters.
func (p ViewParams) EncodeQuery() (url.Values, error) { return p.id.encode() }

// DetailParams configures video.detail.
type DetailParams struct {
	id       identifier
	needElec *bool
}

// DetailByAID identifies video.detail by a validated AV ID.
func DetailByAID(aid ids.AID) DetailParams { return DetailParams{id: byAID(aid)} }

// DetailByBVID identifies video.detail by a validated BV ID.
func DetailByBVID(bvid ids.BVID) DetailParams { return DetailParams{id: byBVID(bvid)} }

// WithElectric returns a copy that controls inclusion of electric-charge data.
func (p DetailParams) WithElectric(enabled bool) DetailParams {
	p.needElec = &enabled
	return p
}

// EncodeQuery returns a new URL query containing the request parameters.
func (p DetailParams) EncodeQuery() (url.Values, error) {
	values, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	if p.needElec != nil {
		if *p.needElec {
			values.Set("need_elec", "1")
		} else {
			values.Set("need_elec", "0")
		}
	}
	return values, nil
}

// PageListParams configures video.pagelist.
type PageListParams struct{ id identifier }

// PageListByAID identifies video.pagelist by a validated AV ID.
func PageListByAID(aid ids.AID) PageListParams { return PageListParams{id: byAID(aid)} }

// PageListByBVID identifies video.pagelist by a validated BV ID.
func PageListByBVID(bvid ids.BVID) PageListParams { return PageListParams{id: byBVID(bvid)} }

// EncodeQuery returns a new URL query containing the request parameters.
func (p PageListParams) EncodeQuery() (url.Values, error) { return p.id.encode() }

// DescParams configures video.desc.
type DescParams struct{ id identifier }

// DescByAID identifies video.desc by a validated AV ID.
func DescByAID(aid ids.AID) DescParams { return DescParams{id: byAID(aid)} }

// DescByBVID identifies video.desc by a validated BV ID.
func DescByBVID(bvid ids.BVID) DescParams { return DescParams{id: byBVID(bvid)} }

// EncodeQuery returns a new URL query containing the request parameters.
func (p DescParams) EncodeQuery() (url.Values, error) { return p.id.encode() }
