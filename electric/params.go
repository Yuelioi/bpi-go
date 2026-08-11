// Package electric contains parameters and stable response models for
// Bilibili's public and account-scoped charging APIs.
package electric

import (
	"net/url"
	"strconv"
	"time"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type MonthUpListParams struct{ upMID ids.MID }

func NewMonthUpListParams(upMID ids.MID) MonthUpListParams {
	return MonthUpListParams{upMID: upMID}
}

func (p MonthUpListParams) EncodeQuery() (url.Values, error) {
	if err := p.upMID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"up_mid": {p.upMID.String()}}, nil
}

type VideoShowParams struct {
	mid  ids.MID
	aid  *ids.AID
	bvid *ids.BVID
}

func NewVideoShowParams(mid ids.MID) VideoShowParams { return VideoShowParams{mid: mid} }

func (p VideoShowParams) WithAID(aid ids.AID) VideoShowParams {
	p.aid = &aid
	return p
}

func (p VideoShowParams) WithBVID(bvid ids.BVID) VideoShowParams {
	p.bvid = &bvid
	return p
}

func (p VideoShowParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"mid": {p.mid.String()}}
	if p.aid != nil {
		if err := p.aid.Validate(); err != nil {
			return nil, err
		}
		values.Set("aid", p.aid.String())
	}
	if p.bvid != nil {
		if err := p.bvid.Validate(); err != nil {
			return nil, err
		}
		values.Set("bvid", p.bvid.String())
	}
	return values, nil
}

type RechargeListParams struct {
	page      uint64
	pageSize  uint64
	beginTime *time.Time
	endTime   *time.Time
}

func NewRechargeListParams(page, pageSize uint64) (RechargeListParams, error) {
	if err := validatePage(page, pageSize); err != nil {
		return RechargeListParams{}, err
	}
	return RechargeListParams{page: page, pageSize: pageSize}, nil
}

func (p RechargeListParams) WithDateRange(begin, end time.Time) (RechargeListParams, error) {
	if end.Before(begin) {
		return RechargeListParams{}, parameterError("endTime", "date must not precede beginTime")
	}
	p.beginTime, p.endTime = &begin, &end
	return p, nil
}

func (p RechargeListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{
		"customerId":  {"10026"},
		"currentPage": {strconv.FormatUint(p.page, 10)},
		"pageSize":    {strconv.FormatUint(p.pageSize, 10)},
	}
	if p.beginTime != nil {
		values.Set("beginTime", p.beginTime.Format(time.DateOnly))
	}
	if p.endTime != nil {
		values.Set("endTime", p.endTime.Format(time.DateOnly))
	}
	return values, nil
}

type PaginationParams struct {
	page     uint64
	pageSize uint64
}

func NewPaginationParams(page, pageSize uint64) (PaginationParams, error) {
	if err := validatePage(page, pageSize); err != nil {
		return PaginationParams{}, err
	}
	return PaginationParams{page: page, pageSize: pageSize}, nil
}

func (p PaginationParams) EncodeQuery() (url.Values, error) {
	return url.Values{"pn": {strconv.FormatUint(p.page, 10)}, "ps": {strconv.FormatUint(p.pageSize, 10)}}, nil
}

type ChargeRecordParams struct {
	page       uint64
	chargeType uint32
}

func NewChargeRecordParams(page uint64, chargeType uint32) (ChargeRecordParams, error) {
	if page == 0 {
		return ChargeRecordParams{}, parameterError("page", "value must be non-zero")
	}
	return ChargeRecordParams{page: page, chargeType: chargeType}, nil
}

func (p ChargeRecordParams) EncodeQuery() (url.Values, error) {
	return url.Values{
		"page": {strconv.FormatUint(p.page, 10)},
		"type": {strconv.FormatUint(uint64(p.chargeType), 10)},
	}, nil
}

type UpMIDParams struct{ upMID ids.MID }

func NewUpMIDParams(upMID ids.MID) UpMIDParams { return UpMIDParams{upMID: upMID} }

func (p UpMIDParams) EncodeQuery() (url.Values, error) {
	if err := p.upMID.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"up_mid": {p.upMID.String()}}, nil
}

type MemberRankParams struct {
	upMID         ids.MID
	page          uint64
	pageSize      uint64
	privilegeType *uint64
}

func NewMemberRankParams(upMID ids.MID, page, pageSize uint64) (MemberRankParams, error) {
	if err := validatePage(page, pageSize); err != nil {
		return MemberRankParams{}, err
	}
	return MemberRankParams{upMID: upMID, page: page, pageSize: pageSize}, nil
}

func (p MemberRankParams) WithPrivilegeType(privilegeType uint64) MemberRankParams {
	p.privilegeType = &privilegeType
	return p
}

func (p MemberRankParams) EncodeQuery() (url.Values, error) {
	if err := p.upMID.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{
		"up_mid": {p.upMID.String()},
		"pn":     {strconv.FormatUint(p.page, 10)},
		"ps":     {strconv.FormatUint(p.pageSize, 10)},
	}
	if p.privilegeType != nil {
		values.Set("privilege_type", strconv.FormatUint(*p.privilegeType, 10))
	}
	return values, nil
}

type RemarkListParams struct {
	page     uint64
	pageSize uint64
	begin    *time.Time
	end      *time.Time
}

func NewRemarkListParams(page, pageSize uint64) (RemarkListParams, error) {
	if err := validatePage(page, pageSize); err != nil {
		return RemarkListParams{}, err
	}
	return RemarkListParams{page: page, pageSize: pageSize}, nil
}

func (p RemarkListParams) WithDateRange(begin, end time.Time) (RemarkListParams, error) {
	if end.Before(begin) {
		return RemarkListParams{}, parameterError("end", "date must not precede begin")
	}
	p.begin, p.end = &begin, &end
	return p, nil
}

func (p RemarkListParams) EncodeQuery() (url.Values, error) {
	values := url.Values{"pn": {strconv.FormatUint(p.page, 10)}, "ps": {strconv.FormatUint(p.pageSize, 10)}}
	if p.begin != nil {
		values.Set("begin", p.begin.Format(time.DateOnly))
	}
	if p.end != nil {
		values.Set("end", p.end.Format(time.DateOnly))
	}
	return values, nil
}

type RemarkDetailParams struct{ id uint64 }

func NewRemarkDetailParams(id uint64) (RemarkDetailParams, error) {
	if id == 0 {
		return RemarkDetailParams{}, parameterError("id", "value must be non-zero")
	}
	return RemarkDetailParams{id: id}, nil
}

func (p RemarkDetailParams) EncodeQuery() (url.Values, error) {
	return url.Values{"id": {strconv.FormatUint(p.id, 10)}}, nil
}

func validatePage(page, pageSize uint64) error {
	if page == 0 {
		return parameterError("page", "value must be non-zero")
	}
	if pageSize == 0 {
		return parameterError("page_size", "value must be non-zero")
	}
	return nil
}

func parameterError(field, message string) error {
	return &bpierr.ParameterError{Field: field, Message: message}
}
