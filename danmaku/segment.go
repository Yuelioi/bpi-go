// Package danmaku contains validated parameters for Bilibili danmaku JSON,
// XML, and protobuf endpoints.
package danmaku

import (
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// SegmentParams configures a real-time protobuf danmaku segment request.
type SegmentParams struct {
	typeID       uint8
	oid          uint64
	segmentIndex uint32
	pid          ids.AID
	pullMode     *uint32
	rangeStart   *uint32
	rangeEnd     *uint32
}

// NewSegmentParams creates parameters with non-zero type, object ID, and
// segment index.
func NewSegmentParams(typeID uint8, oid uint64, segmentIndex uint32) (SegmentParams, error) {
	if typeID == 0 {
		return SegmentParams{}, &bpierr.ParameterError{Field: "type", Message: "danmaku type must be non-zero"}
	}
	if oid == 0 {
		return SegmentParams{}, &bpierr.ParameterError{Field: "oid", Message: "danmaku object ID must be non-zero"}
	}
	if segmentIndex == 0 {
		return SegmentParams{}, &bpierr.ParameterError{Field: "segment_index", Message: "segment index must be non-zero"}
	}
	return SegmentParams{typeID: typeID, oid: oid, segmentIndex: segmentIndex}, nil
}

// WithAID returns a copy with the owning archive ID.
func (p SegmentParams) WithAID(aid ids.AID) (SegmentParams, error) {
	if err := aid.Validate(); err != nil {
		return SegmentParams{}, &bpierr.ParameterError{Field: "pid", Message: "archive ID must be non-zero"}
	}
	p.pid = aid
	return p, nil
}

// WithPullMode returns a copy with an API pull-mode value.
func (p SegmentParams) WithPullMode(mode uint32) SegmentParams {
	p.pullMode = &mode
	return p
}

// WithRange returns a copy with an inclusive millisecond content range.
func (p SegmentParams) WithRange(start, end uint32) (SegmentParams, error) {
	if end < start {
		return SegmentParams{}, &bpierr.ParameterError{Field: "pe", Message: "range end must not precede range start"}
	}
	p.rangeStart = &start
	p.rangeEnd = &end
	return p, nil
}

// EncodeQuery returns a new URL query containing the request parameters.
func (p SegmentParams) EncodeQuery() (url.Values, error) {
	if p.typeID == 0 {
		return nil, &bpierr.ParameterError{Field: "type", Message: "danmaku type must be non-zero"}
	}
	if p.oid == 0 {
		return nil, &bpierr.ParameterError{Field: "oid", Message: "danmaku object ID must be non-zero"}
	}
	if p.segmentIndex == 0 {
		return nil, &bpierr.ParameterError{Field: "segment_index", Message: "segment index must be non-zero"}
	}
	values := url.Values{
		"type":          {strconv.FormatUint(uint64(p.typeID), 10)},
		"oid":           {strconv.FormatUint(p.oid, 10)},
		"segment_index": {strconv.FormatUint(uint64(p.segmentIndex), 10)},
	}
	if p.pid != 0 {
		if err := p.pid.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "pid", Message: "archive ID must be non-zero"}
		}
		values.Set("pid", p.pid.String())
	}
	if p.pullMode != nil {
		values.Set("pull_mode", strconv.FormatUint(uint64(*p.pullMode), 10))
	}
	if p.rangeStart != nil && p.rangeEnd != nil {
		if *p.rangeEnd < *p.rangeStart {
			return nil, &bpierr.ParameterError{Field: "pe", Message: "range end must not precede range start"}
		}
		values.Set("ps", strconv.FormatUint(uint64(*p.rangeStart), 10))
		values.Set("pe", strconv.FormatUint(uint64(*p.rangeEnd), 10))
	}
	return values, nil
}
