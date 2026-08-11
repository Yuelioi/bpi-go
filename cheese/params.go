// Package cheese contains validated parameters and stable response models for
// Bilibili PUGV/course endpoints.
package cheese

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/bangumi"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type EpisodeListParams struct {
	seasonID ids.SeasonID
	pageSize uint32
	page     uint32
}

func NewEpisodeListParams(seasonID ids.SeasonID) EpisodeListParams {
	return EpisodeListParams{seasonID: seasonID}
}

func (p EpisodeListParams) WithPageSize(size uint32) (EpisodeListParams, error) {
	if size == 0 {
		return EpisodeListParams{}, &bpierr.ParameterError{Field: "ps", Message: "page size must be at least 1"}
	}
	p.pageSize = size
	return p, nil
}

func (p EpisodeListParams) WithPage(page uint32) (EpisodeListParams, error) {
	if page == 0 {
		return EpisodeListParams{}, &bpierr.ParameterError{Field: "pn", Message: "page must be at least 1"}
	}
	p.page = page
	return p, nil
}

func (p EpisodeListParams) EncodeQuery() (url.Values, error) {
	if err := p.seasonID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "season_id", Message: "season ID is invalid"}
	}
	values := url.Values{"season_id": {p.seasonID.String()}}
	if p.pageSize != 0 {
		values.Set("ps", strconv.FormatUint(uint64(p.pageSize), 10))
	}
	if p.page != 0 {
		values.Set("pn", strconv.FormatUint(uint64(p.page), 10))
	}
	return values, nil
}

type VideoQuality = bangumi.VideoQuality
type FormatFlags = bangumi.FormatFlags

const (
	Quality480P = bangumi.Quality480P
	Quality4K   = bangumi.Quality4K
	FormatDASH  = bangumi.FormatDASH
	Format4K    = bangumi.Format4K
)

type PlayURLParams struct {
	aid       ids.AID
	episodeID ids.EpisodeID
	cid       ids.CID
	quality   *VideoQuality
	flags     *FormatFlags
}

func NewPlayURLParams(aid ids.AID, episodeID ids.EpisodeID, cid ids.CID) PlayURLParams {
	return PlayURLParams{aid: aid, episodeID: episodeID, cid: cid}
}

func (p PlayURLParams) WithQuality(quality VideoQuality) PlayURLParams {
	p.quality = &quality
	return p
}

func (p PlayURLParams) WithFormatFlags(flags FormatFlags) PlayURLParams {
	p.flags = &flags
	return p
}

func (p PlayURLParams) EncodeQuery() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "avid", Message: "video ID is invalid"}
	}
	if err := p.episodeID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "ep_id", Message: "episode ID is invalid"}
	}
	if err := p.cid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID is invalid"}
	}
	values := url.Values{
		"avid":  {p.aid.String()},
		"ep_id": {p.episodeID.String()},
		"cid":   {p.cid.String()},
		"fnver": {"0"},
	}
	if p.quality != nil {
		if *p.quality == 0 {
			return nil, &bpierr.ParameterError{Field: "qn", Message: "quality must be non-zero"}
		}
		values.Set("qn", strconv.FormatUint(uint64(*p.quality), 10))
	}
	if p.flags != nil {
		if *p.flags == 0 {
			return nil, &bpierr.ParameterError{Field: "fnval", Message: "format flags must be non-zero"}
		}
		values.Set("fnval", strconv.FormatUint(uint64(*p.flags), 10))
		if *p.flags&Format4K != 0 {
			values.Set("fourk", "1")
		}
	}
	return values, nil
}
