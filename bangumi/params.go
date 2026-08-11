package bangumi

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type InfoParams struct{ mediaID ids.MediaID }

func NewInfoParams(mediaID ids.MediaID) InfoParams { return InfoParams{mediaID: mediaID} }

func (p InfoParams) EncodeQuery() (url.Values, error) {
	if err := p.mediaID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "media_id", Message: "media ID is invalid"}
	}
	return url.Values{"media_id": {p.mediaID.String()}}, nil
}

type DetailParams struct {
	seasonID  ids.SeasonID
	episodeID ids.EpisodeID
}

func DetailBySeasonID(id ids.SeasonID) DetailParams   { return DetailParams{seasonID: id} }
func DetailByEpisodeID(id ids.EpisodeID) DetailParams { return DetailParams{episodeID: id} }

func (p DetailParams) EncodeQuery() (url.Values, error) {
	switch {
	case p.seasonID != 0 && p.episodeID == 0:
		if err := p.seasonID.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "season_id", Message: "season ID is invalid"}
		}
		return url.Values{"season_id": {p.seasonID.String()}}, nil
	case p.episodeID != 0 && p.seasonID == 0:
		if err := p.episodeID.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "ep_id", Message: "episode ID is invalid"}
		}
		return url.Values{"ep_id": {p.episodeID.String()}}, nil
	default:
		return nil, &bpierr.ParameterError{Field: "detail_id", Message: "exactly one season or episode ID is required"}
	}
}

type SectionsParams struct{ seasonID ids.SeasonID }

func NewSectionsParams(seasonID ids.SeasonID) SectionsParams {
	return SectionsParams{seasonID: seasonID}
}

func (p SectionsParams) EncodeQuery() (url.Values, error) {
	if err := p.seasonID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "season_id", Message: "season ID is invalid"}
	}
	return url.Values{"season_id": {p.seasonID.String()}}, nil
}

type VideoQuality uint64

const (
	Quality360P  VideoQuality = 16
	Quality480P  VideoQuality = 32
	Quality720P  VideoQuality = 64
	Quality1080P VideoQuality = 80
	Quality4K    VideoQuality = 120
)

type FormatFlags uint64

const (
	FormatDASH  FormatFlags = 16
	FormatHDR   FormatFlags = 64
	Format4K    FormatFlags = 128
	FormatDolby FormatFlags = 256
	Format8K    FormatFlags = 1024
	FormatAV1   FormatFlags = 2048
)

type PlayURLParams struct {
	episodeID ids.EpisodeID
	cid       ids.CID
	quality   *VideoQuality
	flags     *FormatFlags
}

func PlayURLByEpisodeID(id ids.EpisodeID) PlayURLParams { return PlayURLParams{episodeID: id} }
func PlayURLByCID(id ids.CID) PlayURLParams             { return PlayURLParams{cid: id} }

func (p PlayURLParams) WithQuality(quality VideoQuality) PlayURLParams {
	p.quality = &quality
	return p
}

func (p PlayURLParams) WithFormatFlags(flags FormatFlags) PlayURLParams {
	p.flags = &flags
	return p
}

func (p PlayURLParams) EncodeQuery() (url.Values, error) {
	values := url.Values{"fnver": {"0"}}
	switch {
	case p.episodeID != 0 && p.cid == 0:
		if err := p.episodeID.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "ep_id", Message: "episode ID is invalid"}
		}
		values.Set("ep_id", p.episodeID.String())
	case p.cid != 0 && p.episodeID == 0:
		if err := p.cid.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID is invalid"}
		}
		values.Set("cid", p.cid.String())
	default:
		return nil, &bpierr.ParameterError{Field: "play_id", Message: "exactly one episode or content ID is required"}
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
