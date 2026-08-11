// Package dynamic contains validated parameters and stable response models for
// Bilibili's dynamic-feed endpoints.
package dynamic

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

const (
	DefaultAllFeatures    = "itemOpusStyle,listOnlyfans,opusBigCover,onlyfansVote,decorationCard,onlyfansAssetsV2,forwardListHidden,ugcDelete"
	DefaultWebLocation    = "333.1365"
	DefaultDetailFeatures = "htmlNewStyle,itemOpusStyle,decorationCard"
)

type AllParams struct {
	features       string
	webLocation    string
	hostMID        *ids.MID
	offset         string
	updateBaseline string
}

func NewAllParams() AllParams {
	return AllParams{features: DefaultAllFeatures, webLocation: DefaultWebLocation}
}

func (p AllParams) WithFeatures(features string) (AllParams, error) {
	value, err := nonBlank("features", features)
	if err != nil {
		return AllParams{}, err
	}
	p.features = value
	return p, nil
}

func (p AllParams) WithWebLocation(location string) (AllParams, error) {
	value, err := nonBlank("web_location", location)
	if err != nil {
		return AllParams{}, err
	}
	p.webLocation = value
	return p, nil
}

func (p AllParams) WithHostMID(mid ids.MID) AllParams {
	p.hostMID = &mid
	return p
}

func (p AllParams) WithOffset(offset string) (AllParams, error) {
	value, err := nonBlank("offset", offset)
	if err != nil {
		return AllParams{}, err
	}
	p.offset = value
	return p, nil
}

func (p AllParams) WithUpdateBaseline(baseline string) (AllParams, error) {
	value, err := nonBlank("update_baseline", baseline)
	if err != nil {
		return AllParams{}, err
	}
	p.updateBaseline = value
	return p, nil
}

func (p AllParams) EncodeQuery() (url.Values, error) {
	if p.features == "" || p.webLocation == "" {
		return nil, &bpierr.ParameterError{Field: "params", Message: "use NewAllParams to initialize required defaults"}
	}
	values := url.Values{"features": {p.features}, "web_location": {p.webLocation}}
	if p.hostMID != nil {
		if err := p.hostMID.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "host_mid", Message: "member ID is invalid"}
		}
		values.Set("host_mid", p.hostMID.String())
	}
	if p.offset != "" {
		values.Set("offset", p.offset)
	}
	if p.updateBaseline != "" {
		values.Set("update_baseline", p.updateBaseline)
	}
	return values, nil
}

type CheckNewParams struct {
	baseline string
	typeName string
}

func NewCheckNewParams(baseline string) (CheckNewParams, error) {
	value, err := nonBlank("update_baseline", baseline)
	if err != nil {
		return CheckNewParams{}, err
	}
	return CheckNewParams{baseline: value}, nil
}

func (p CheckNewParams) WithType(typeName string) (CheckNewParams, error) {
	value, err := nonBlank("type", typeName)
	if err != nil {
		return CheckNewParams{}, err
	}
	p.typeName = value
	return p, nil
}

func (p CheckNewParams) EncodeQuery() (url.Values, error) {
	if p.baseline == "" {
		return nil, &bpierr.ParameterError{Field: "update_baseline", Message: "value cannot be blank"}
	}
	values := url.Values{"update_baseline": {p.baseline}}
	if p.typeName != "" {
		values.Set("type", p.typeName)
	}
	return values, nil
}

type NavFeedParams struct {
	updateBaseline string
	offset         string
}

func NewNavFeedParams() NavFeedParams { return NavFeedParams{} }

func (p NavFeedParams) WithUpdateBaseline(baseline string) (NavFeedParams, error) {
	value, err := nonBlank("update_baseline", baseline)
	if err != nil {
		return NavFeedParams{}, err
	}
	p.updateBaseline = value
	return p, nil
}

func (p NavFeedParams) WithOffset(offset string) (NavFeedParams, error) {
	value, err := nonBlank("offset", offset)
	if err != nil {
		return NavFeedParams{}, err
	}
	p.offset = value
	return p, nil
}

func (p NavFeedParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.updateBaseline != "" {
		values.Set("update_baseline", p.updateBaseline)
	}
	if p.offset != "" {
		values.Set("offset", p.offset)
	}
	return values, nil
}

type LiveUsersParams struct{ size uint32 }

func NewLiveUsersParams() LiveUsersParams { return LiveUsersParams{} }

func (p LiveUsersParams) WithSize(size uint32) (LiveUsersParams, error) {
	if size == 0 {
		return LiveUsersParams{}, &bpierr.ParameterError{Field: "size", Message: "value must be non-zero"}
	}
	p.size = size
	return p, nil
}

func (p LiveUsersParams) EncodeQuery() (url.Values, error) {
	values := url.Values{}
	if p.size != 0 {
		values.Set("size", strconv.FormatUint(uint64(p.size), 10))
	}
	return values, nil
}

type UpUsersParams struct{ teenagersMode bool }

func NewUpUsersParams() UpUsersParams { return UpUsersParams{} }

func (p UpUsersParams) WithTeenagersMode(enabled bool) UpUsersParams {
	p.teenagersMode = enabled
	return p
}

func (p UpUsersParams) EncodeQuery() (url.Values, error) {
	value := "0"
	if p.teenagersMode {
		value = "1"
	}
	return url.Values{"teenagers_mode": {value}}, nil
}

type DetailParams struct {
	id       ids.DynamicID
	features string
}

func NewDetailParams(id ids.DynamicID) DetailParams {
	return DetailParams{id: id, features: DefaultDetailFeatures}
}

func (p DetailParams) WithFeatures(features string) (DetailParams, error) {
	value, err := nonBlank("features", features)
	if err != nil {
		return DetailParams{}, err
	}
	p.features = value
	return p, nil
}

func (p DetailParams) EncodeQuery() (url.Values, error) {
	if err := p.id.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "id", Message: "dynamic ID is invalid"}
	}
	if p.features == "" {
		return nil, &bpierr.ParameterError{Field: "features", Message: "value cannot be blank"}
	}
	return url.Values{"id": {p.id.String()}, "features": {p.features}}, nil
}

type OffsetParams struct {
	id     ids.DynamicID
	offset string
}

func NewOffsetParams(id ids.DynamicID) OffsetParams { return OffsetParams{id: id} }

func (p OffsetParams) WithOffset(offset string) (OffsetParams, error) {
	value, err := nonBlank("offset", offset)
	if err != nil {
		return OffsetParams{}, err
	}
	p.offset = value
	return p, nil
}

func (p OffsetParams) EncodeQuery() (url.Values, error) {
	if err := p.id.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "id", Message: "dynamic ID is invalid"}
	}
	values := url.Values{"id": {p.id.String()}}
	if p.offset != "" {
		values.Set("offset", p.offset)
	}
	return values, nil
}

type ItemParams struct{ id ids.DynamicID }

func NewItemParams(id ids.DynamicID) ItemParams { return ItemParams{id: id} }

func (p ItemParams) EncodeQuery() (url.Values, error) {
	if err := p.id.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "id", Message: "dynamic ID is invalid"}
	}
	return url.Values{"id": {p.id.String()}}, nil
}

type LotteryNoticeParams struct{ businessID ids.DynamicID }

func NewLotteryNoticeParams(id ids.DynamicID) LotteryNoticeParams {
	return LotteryNoticeParams{businessID: id}
}

func (p LotteryNoticeParams) EncodeQuery(csrf string) (url.Values, error) {
	if err := p.businessID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "business_id", Message: "dynamic ID is invalid"}
	}
	return url.Values{
		"business_id":   {p.businessID.String()},
		"business_type": {"1"},
		"csrf":          {csrf},
	}, nil
}

func FeedBannerQuery() url.Values {
	return url.Values{"platform": {"1"}, "position": {"web动态"}, "web_location": {DefaultWebLocation}}
}

func nonBlank(field, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", &bpierr.ParameterError{Field: field, Message: "value cannot be blank"}
	}
	return value, nil
}
