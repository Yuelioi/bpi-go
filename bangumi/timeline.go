// Package bangumi contains validated parameters and response models for the
// Bilibili bangumi domain.
package bangumi

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// TimelineType identifies the catalogue represented by a bangumi timeline.
type TimelineType int32

const (
	// TimelineAnime selects anime/bangumi releases.
	TimelineAnime TimelineType = 1
	// TimelineMovie selects movie releases.
	TimelineMovie TimelineType = 3
	// TimelineChineseAnimation selects Chinese animation releases.
	TimelineChineseAnimation TimelineType = 4
)

func (typ TimelineType) validate() error {
	switch typ {
	case TimelineAnime, TimelineMovie, TimelineChineseAnimation:
		return nil
	default:
		return &bpierr.ParameterError{Field: "types", Message: "timeline type must be anime, movie, or Chinese animation"}
	}
}

// TimelineParams configures bangumi.timeline.
type TimelineParams struct {
	typeID TimelineType
	before int32
	after  int32
}

// NewTimelineParams creates timeline parameters with day offsets in [0, 7].
func NewTimelineParams(typeID TimelineType, before, after int32) (TimelineParams, error) {
	if err := typeID.validate(); err != nil {
		return TimelineParams{}, err
	}
	if before < 0 || before > 7 {
		return TimelineParams{}, &bpierr.ParameterError{Field: "before", Message: "value must be between 0 and 7"}
	}
	if after < 0 || after > 7 {
		return TimelineParams{}, &bpierr.ParameterError{Field: "after", Message: "value must be between 0 and 7"}
	}
	return TimelineParams{typeID: typeID, before: before, after: after}, nil
}

// EncodeQuery returns a new URL query containing the request parameters.
func (p TimelineParams) EncodeQuery() (url.Values, error) {
	if err := p.typeID.validate(); err != nil {
		return nil, err
	}
	if p.before < 0 || p.before > 7 {
		return nil, &bpierr.ParameterError{Field: "before", Message: "value must be between 0 and 7"}
	}
	if p.after < 0 || p.after > 7 {
		return nil, &bpierr.ParameterError{Field: "after", Message: "value must be between 0 and 7"}
	}
	return url.Values{
		"types":  {strconv.FormatInt(int64(p.typeID), 10)},
		"before": {strconv.FormatInt(int64(p.before), 10)},
		"after":  {strconv.FormatInt(int64(p.after), 10)},
	}, nil
}

// TimelineDay is one day in a bangumi release timeline.
type TimelineDay struct {
	Date      string            `json:"date"`
	DateUnix  int64             `json:"date_ts"`
	DayOfWeek int32             `json:"day_of_week"`
	Episodes  []TimelineEpisode `json:"episodes"`
	IsToday   int32             `json:"is_today"`
}

// TimelineEpisode is one scheduled episode.
type TimelineEpisode struct {
	Cover        string `json:"cover"`
	Delay        int32  `json:"delay"`
	DelayID      int64  `json:"delay_id"`
	DelayIndex   string `json:"delay_index"`
	DelayReason  string `json:"delay_reason"`
	EpisodeCover string `json:"ep_cover"`
	EpisodeID    int64  `json:"episode_id"`
	PublishIndex string `json:"pub_index"`
	PublishTime  string `json:"pub_time"`
	PublishUnix  int64  `json:"pub_ts"`
	Published    int32  `json:"published"`
	Follows      string `json:"follows"`
	Plays        string `json:"plays"`
	SeasonID     int64  `json:"season_id"`
	SquareCover  string `json:"square_cover"`
	Title        string `json:"title"`
}
