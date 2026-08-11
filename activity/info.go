// Package activity contains validated parameters and response models for the
// Bilibili activity domain.
package activity

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// Info is the stable activity-subject payload returned by activity.info.
type Info struct {
	ID          uint64 `json:"id"`
	StartTime   int64  `json:"stime"`
	EndTime     int64  `json:"etime"`
	CreateTime  int64  `json:"ctime"`
	ModifyTime  int64  `json:"mtime"`
	Name        string `json:"name"`
	URL         string `json:"act_url"`
	Cover       string `json:"cover"`
	Description string `json:"dic"`
	H5Cover     string `json:"h5_cover"`
	AndroidURL  string `json:"android_url"`
	IOSURL      string `json:"ios_url"`
	ChildSIDs   string `json:"child_sids"`
	LID         *int64 `json:"lid"`
}

// InfoParams identifies an activity subject and an optional source video.
// Construct it with NewInfoParams.
type InfoParams struct {
	sid  uint64
	bvid ids.BVID
}

// NewInfoParams creates parameters for a non-zero activity subject ID.
func NewInfoParams(sid uint64) (InfoParams, error) {
	if sid == 0 {
		return InfoParams{}, &bpierr.ParameterError{Field: "sid", Message: "activity subject ID must be non-zero"}
	}
	return InfoParams{sid: sid}, nil
}

// WithBVID returns a copy that includes a validated source video ID.
func (p InfoParams) WithBVID(bvid ids.BVID) (InfoParams, error) {
	if err := bvid.Validate(); err != nil {
		return InfoParams{}, &bpierr.ParameterError{Field: "bvid", Message: "source video ID is invalid"}
	}
	p.bvid = bvid
	return p, nil
}

// Validate checks that the parameter value is ready for a request.
func (p InfoParams) Validate() error {
	if p.sid == 0 {
		return &bpierr.ParameterError{Field: "sid", Message: "activity subject ID must be non-zero"}
	}
	if p.bvid != "" {
		if err := p.bvid.Validate(); err != nil {
			return &bpierr.ParameterError{Field: "bvid", Message: "source video ID is invalid"}
		}
	}
	return nil
}

// EncodeQuery returns a new URL query containing the request parameters.
func (p InfoParams) EncodeQuery() (url.Values, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"sid": {strconv.FormatUint(p.sid, 10)}}
	if p.bvid != "" {
		values.Set("bvid", p.bvid.String())
	}
	return values, nil
}
