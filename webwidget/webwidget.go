// Package webwidget contains validated parameters and response models for
// Bilibili's public Web-widget endpoints.
package webwidget

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

const defaultHeaderResourceID = uint32(142)

// HeaderPageParams configures web_widget.header_page. Its zero value selects
// Bilibili's default Web header resource 142.
type HeaderPageParams struct {
	resourceID uint32
}

// NewHeaderPageParams returns parameters for the default header resource.
func NewHeaderPageParams() HeaderPageParams { return HeaderPageParams{} }

// WithResourceID returns a copy configured for a non-zero resource ID.
func (p HeaderPageParams) WithResourceID(resourceID uint32) (HeaderPageParams, error) {
	if resourceID == 0 {
		return HeaderPageParams{}, &bpierr.ParameterError{Field: "resource_id", Message: "value must be non-zero"}
	}
	p.resourceID = resourceID
	return p, nil
}

// EncodeQuery returns a new URL query with the effective resource ID.
func (p HeaderPageParams) EncodeQuery() (url.Values, error) {
	resourceID := p.resourceID
	if resourceID == 0 {
		resourceID = defaultHeaderResourceID
	}
	return url.Values{"resource_id": {strconv.FormatUint(uint64(resourceID), 10)}}, nil
}

// RegionBannerParams configures web_widget.region_banner.
type RegionBannerParams struct {
	regionID uint32
}

// NewRegionBannerParams creates parameters for a non-zero video-region ID.
func NewRegionBannerParams(regionID uint32) (RegionBannerParams, error) {
	if regionID == 0 {
		return RegionBannerParams{}, &bpierr.ParameterError{Field: "region_id", Message: "value must be non-zero"}
	}
	return RegionBannerParams{regionID: regionID}, nil
}

// EncodeQuery returns a new URL query for the configured video region.
func (p RegionBannerParams) EncodeQuery() (url.Values, error) {
	if p.regionID == 0 {
		return nil, &bpierr.ParameterError{Field: "region_id", Message: "value must be non-zero"}
	}
	return url.Values{"region_id": {strconv.FormatUint(uint64(p.regionID), 10)}}, nil
}

// RegionBanner is one public regional carousel item.
type RegionBanner struct {
	Image    string `json:"image"`
	Title    string `json:"title"`
	Subtitle string `json:"sub_title"`
	URL      string `json:"url"`
	RegionID int64  `json:"rid"`
}

// RegionBannerData is the regional carousel payload.
type RegionBannerData struct {
	Items []RegionBanner `json:"region_banner_list"`
}

// HeaderData is the public home-page header payload. SplitLayerJSON preserves
// the wire representation while SplitLayer exposes its parsed structure.
type HeaderData struct {
	Name             string      `json:"name"`
	Picture          string      `json:"pic"`
	LitPicture       string      `json:"litpic"`
	URL              string      `json:"url"`
	IsSplitLayer     uint32      `json:"is_split_layer"`
	SplitLayerJSON   string      `json:"split_layer"`
	ParsedSplitLayer *SplitLayer `json:"-"`
}

// UnmarshalJSON decodes the JSON object embedded in the split_layer string so
// callers do not need a second parsing step.
func (h *HeaderData) UnmarshalJSON(data []byte) error {
	type wireHeader HeaderData
	var wire wireHeader
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if strings.TrimSpace(wire.SplitLayerJSON) != "" {
		var layer SplitLayer
		if err := json.Unmarshal([]byte(wire.SplitLayerJSON), &layer); err != nil {
			return fmt.Errorf("decode split_layer: %w", err)
		}
		wire.ParsedSplitLayer = &layer
	}
	*h = HeaderData(wire)
	return nil
}

// SplitLayer is the parsed, versioned header-layer definition.
type SplitLayer struct {
	Version string  `json:"version"`
	Layers  []Layer `json:"layers"`
}

// Layer is one independently transformed header asset layer.
type Layer struct {
	Resources []Resource `json:"resources"`
	Scale     Scale      `json:"scale"`
	Rotate    Rotate     `json:"rotate"`
	Translate Translate  `json:"translate"`
	Blur      Blur       `json:"blur"`
	Opacity   Opacity    `json:"opacity"`
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
}

// Resource is one asset used by a header layer.
type Resource struct {
	Source string `json:"src"`
	ID     int64  `json:"id"`
}

type Scale struct {
	Initial *float64 `json:"initial"`
}

type Rotate struct {
	Offset *int64 `json:"offset"`
}

type Translate struct {
	Offset  []int64 `json:"offset"`
	Initial []int64 `json:"initial"`
}

type Blur struct {
	Initial *int64 `json:"initial"`
}

type Opacity struct {
	Wrap    string   `json:"wrap"`
	Initial *float64 `json:"initial"`
}

// RegionCount maps Bilibili region identifiers to today's upload counts.
type RegionCount map[string]uint64

// OnlineData is the public per-region online-upload payload.
type OnlineData struct {
	RegionCount RegionCount `json:"region_count"`
}
