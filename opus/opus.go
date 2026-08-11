// Package opus contains validated parameters and response models for public
// Bilibili opus feeds.
package opus

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

const spaceFeedWebLocation = "333.1387"

// SpaceFeedKind selects the kind of entries included in a space feed.
type SpaceFeedKind string

const (
	SpaceFeedAll     SpaceFeedKind = "all"
	SpaceFeedArticle SpaceFeedKind = "article"
	SpaceFeedDynamic SpaceFeedKind = "dynamic"
)

func (kind SpaceFeedKind) effective() (SpaceFeedKind, error) {
	if kind == "" {
		return SpaceFeedAll, nil
	}
	switch kind {
	case SpaceFeedAll, SpaceFeedArticle, SpaceFeedDynamic:
		return kind, nil
	default:
		return "", &bpierr.ParameterError{Field: "type", Message: "kind must be all, article, or dynamic"}
	}
}

// SpaceFeedParams configures opus.space_feed. Page zero and kind all are the
// endpoint's default first-page request.
type SpaceFeedParams struct {
	mid    ids.MID
	page   uint32
	offset string
	kind   SpaceFeedKind
}

// NewSpaceFeedParams creates first-page feed parameters for a member.
func NewSpaceFeedParams(mid ids.MID) (SpaceFeedParams, error) {
	if err := mid.Validate(); err != nil {
		return SpaceFeedParams{}, &bpierr.ParameterError{Field: "host_mid", Message: "member ID must be non-zero"}
	}
	return SpaceFeedParams{mid: mid}, nil
}

// WithPage returns a copy configured for the supplied zero-based page.
func (p SpaceFeedParams) WithPage(page uint32) SpaceFeedParams {
	p.page = page
	return p
}

// WithOffset returns a copy configured with a non-blank continuation token.
func (p SpaceFeedParams) WithOffset(offset string) (SpaceFeedParams, error) {
	if strings.TrimSpace(offset) == "" {
		return SpaceFeedParams{}, &bpierr.ParameterError{Field: "offset", Message: "offset cannot be blank"}
	}
	p.offset = offset
	return p, nil
}

// WithKind returns a copy configured for one supported feed kind.
func (p SpaceFeedParams) WithKind(kind SpaceFeedKind) (SpaceFeedParams, error) {
	if _, err := kind.effective(); err != nil {
		return SpaceFeedParams{}, err
	}
	p.kind = kind
	return p, nil
}

// EncodeQuery returns a new URL query containing the effective feed options.
func (p SpaceFeedParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "host_mid", Message: "member ID must be non-zero"}
	}
	kind, err := p.kind.effective()
	if err != nil {
		return nil, err
	}
	values := url.Values{
		"host_mid":     {p.mid.String()},
		"page":         {strconv.FormatUint(uint64(p.page), 10)},
		"type":         {string(kind)},
		"web_location": {spaceFeedWebLocation},
	}
	if p.offset != "" {
		values.Set("offset", p.offset)
	}
	return values, nil
}

type SpaceCover struct {
	Height uint32 `json:"height"`
	URL    string `json:"url"`
	Width  uint32 `json:"width"`
}

type SpaceStat struct {
	Like string  `json:"like"`
	View *string `json:"view"`
}

type SpaceItem struct {
	Content string      `json:"content"`
	Cover   *SpaceCover `json:"cover"`
	JumpURL string      `json:"jump_url"`
	OpusID  string      `json:"opus_id"`
	Stat    SpaceStat   `json:"stat"`
}

// SpaceFeed is one page of public opus entries from a member's space.
type SpaceFeed struct {
	HasMore   bool        `json:"has_more"`
	Items     []SpaceItem `json:"items"`
	Offset    string      `json:"offset"`
	UpdateNum uint32      `json:"update_num"`
}
