// Package audio contains validated parameters and stable response models for
// Bilibili audio endpoints.
package audio

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type SongParams struct{ sid ids.AudioID }

func NewSongParams(sid ids.AudioID) SongParams { return SongParams{sid: sid} }

func (p SongParams) EncodeQuery() (url.Values, error) {
	if err := p.sid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "sid", Message: "audio ID is invalid"}
	}
	return url.Values{"sid": {p.sid.String()}}, nil
}

type PageParams struct {
	page uint32
	size uint32
}

func NewPageParams(page, size uint32) (PageParams, error) {
	if page == 0 {
		return PageParams{}, &bpierr.ParameterError{Field: "pn", Message: "page must be at least 1"}
	}
	if size == 0 {
		return PageParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be at least 1"}
	}
	return PageParams{page: page, size: size}, nil
}

func (p PageParams) EncodeQuery() (url.Values, error) {
	if p.page == 0 {
		return nil, &bpierr.ParameterError{Field: "pn", Message: "page must be at least 1"}
	}
	if p.size == 0 {
		return nil, &bpierr.ParameterError{Field: "ps", Message: "page size must be at least 1"}
	}
	return url.Values{
		"pn": {strconv.FormatUint(uint64(p.page), 10)},
		"ps": {strconv.FormatUint(uint64(p.size), 10)},
	}, nil
}

type CollectionInfoParams struct{ id uint64 }

func NewCollectionInfoParams(id uint64) (CollectionInfoParams, error) {
	if id == 0 {
		return CollectionInfoParams{}, &bpierr.ParameterError{Field: "sid", Message: "collection ID must be non-zero"}
	}
	return CollectionInfoParams{id: id}, nil
}

func (p CollectionInfoParams) EncodeQuery() (url.Values, error) {
	if p.id == 0 {
		return nil, &bpierr.ParameterError{Field: "sid", Message: "collection ID must be non-zero"}
	}
	return url.Values{"sid": {strconv.FormatUint(p.id, 10)}}, nil
}

type Quality uint32

const (
	QualitySmooth Quality = iota
	QualityStandard
	QualityHigh
	QualityLossless
)

func (quality Quality) validate() error {
	if quality > QualityLossless {
		return &bpierr.ParameterError{Field: "quality", Message: "quality must be between 0 and 3"}
	}
	return nil
}

type StreamURLWebParams struct {
	sid       ids.AudioID
	quality   Quality
	privilege uint32
}

func NewStreamURLWebParams(sid ids.AudioID) StreamURLWebParams {
	return StreamURLWebParams{sid: sid, quality: QualityHigh, privilege: 2}
}

func (p StreamURLWebParams) WithQuality(quality Quality) StreamURLWebParams {
	p.quality = quality
	return p
}

func (p StreamURLWebParams) WithPrivilege(privilege uint32) (StreamURLWebParams, error) {
	if privilege == 0 {
		return StreamURLWebParams{}, &bpierr.ParameterError{Field: "privilege", Message: "privilege must be non-zero"}
	}
	p.privilege = privilege
	return p, nil
}

func (p StreamURLWebParams) EncodeQuery() (url.Values, error) {
	if err := p.sid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "sid", Message: "audio ID is invalid"}
	}
	if err := p.quality.validate(); err != nil {
		return nil, err
	}
	if p.privilege == 0 {
		return nil, &bpierr.ParameterError{Field: "privilege", Message: "privilege must be non-zero"}
	}
	return url.Values{
		"sid":       {p.sid.String()},
		"quality":   {strconv.FormatUint(uint64(p.quality), 10)},
		"privilege": {strconv.FormatUint(uint64(p.privilege), 10)},
	}, nil
}

type StreamURLParams struct {
	sid       ids.AudioID
	quality   Quality
	privilege uint32
	mid       uint64
	platform  string
}

func NewStreamURLParams(sid ids.AudioID, quality Quality) StreamURLParams {
	return StreamURLParams{sid: sid, quality: quality, privilege: 2, mid: 2, platform: "android"}
}

func (p StreamURLParams) WithPrivilege(privilege uint32) (StreamURLParams, error) {
	if privilege == 0 {
		return StreamURLParams{}, &bpierr.ParameterError{Field: "privilege", Message: "privilege must be non-zero"}
	}
	p.privilege = privilege
	return p, nil
}

func (p StreamURLParams) WithMID(mid uint64) (StreamURLParams, error) {
	if mid == 0 {
		return StreamURLParams{}, &bpierr.ParameterError{Field: "mid", Message: "member ID must be non-zero"}
	}
	p.mid = mid
	return p, nil
}

func (p StreamURLParams) WithPlatform(platform string) (StreamURLParams, error) {
	platform = strings.TrimSpace(platform)
	if platform == "" {
		return StreamURLParams{}, &bpierr.ParameterError{Field: "platform", Message: "platform cannot be blank"}
	}
	p.platform = platform
	return p, nil
}

func (p StreamURLParams) EncodeQuery() (url.Values, error) {
	if err := p.sid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "songid", Message: "audio ID is invalid"}
	}
	if err := p.quality.validate(); err != nil {
		return nil, err
	}
	if p.privilege == 0 || p.mid == 0 || strings.TrimSpace(p.platform) == "" {
		return nil, &bpierr.ParameterError{Field: "stream", Message: "privilege, member ID, and platform are required"}
	}
	return url.Values{
		"songid":    {p.sid.String()},
		"quality":   {strconv.FormatUint(uint64(p.quality), 10)},
		"privilege": {strconv.FormatUint(uint64(p.privilege), 10)},
		"mid":       {strconv.FormatUint(p.mid, 10)},
		"platform":  {p.platform},
	}, nil
}

type RankListType uint32

const (
	RankHot      RankListType = 1
	RankOriginal RankListType = 2
)

func NewCustomRankListType(value uint32) (RankListType, error) {
	if value == 0 {
		return 0, &bpierr.ParameterError{Field: "list_type", Message: "list type must be non-zero"}
	}
	return RankListType(value), nil
}

type RankPeriodParams struct{ listType RankListType }

func NewRankPeriodParams(listType RankListType) RankPeriodParams {
	return RankPeriodParams{listType: listType}
}

func (p RankPeriodParams) EncodeQuery(csrf string) (url.Values, error) {
	if p.listType == 0 {
		return nil, &bpierr.ParameterError{Field: "list_type", Message: "list type must be non-zero"}
	}
	return url.Values{
		"list_type": {strconv.FormatUint(uint64(p.listType), 10)},
		"csrf":      {csrf},
	}, nil
}

type RankListParams struct{ id uint64 }

func NewRankListParams(id uint64) (RankListParams, error) {
	if id == 0 {
		return RankListParams{}, &bpierr.ParameterError{Field: "list_id", Message: "list ID must be non-zero"}
	}
	return RankListParams{id: id}, nil
}

func (p RankListParams) EncodeQuery(csrf string) (url.Values, error) {
	if p.id == 0 {
		return nil, &bpierr.ParameterError{Field: "list_id", Message: "list ID must be non-zero"}
	}
	return url.Values{
		"list_id": {strconv.FormatUint(p.id, 10)},
		"csrf":    {csrf},
	}, nil
}
