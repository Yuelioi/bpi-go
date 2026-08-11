package danmaku

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type HistoryDatesParams struct {
	typeID uint8
	oid    ids.CID
	month  string
}

func NewHistoryDatesParams(oid ids.CID, month string) (HistoryDatesParams, error) {
	if _, err := time.Parse("2006-01", month); err != nil {
		return HistoryDatesParams{}, parameterError("month", "month must use YYYY-MM format")
	}
	return HistoryDatesParams{typeID: 1, oid: oid, month: month}, nil
}

func (p HistoryDatesParams) WithType(typeID uint8) (HistoryDatesParams, error) {
	if typeID == 0 {
		return HistoryDatesParams{}, parameterError("type", "danmaku type must be non-zero")
	}
	p.typeID = typeID
	return p, nil
}

func (p HistoryDatesParams) EncodeQuery() (url.Values, error) {
	if err := p.oid.Validate(); err != nil {
		return nil, err
	}
	return url.Values{
		"type":  {strconv.FormatUint(uint64(p.typeID), 10)},
		"oid":   {p.oid.String()},
		"month": {p.month},
	}, nil
}

type SnapshotParams struct{ archive string }

func NewSnapshotByAID(aid ids.AID) SnapshotParams { return SnapshotParams{archive: aid.String()} }

func NewSnapshotByBVID(bvid ids.BVID) SnapshotParams { return SnapshotParams{archive: bvid.String()} }

func (p SnapshotParams) EncodeQuery() (url.Values, error) {
	if strings.TrimSpace(p.archive) == "" {
		return nil, parameterError("aid", "archive identifier cannot be blank")
	}
	return url.Values{"aid": {p.archive}}, nil
}

type ThumbupStatsParams struct {
	oid ids.CID
	ids []uint64
}

func NewThumbupStatsParams(oid ids.CID, danmakuIDs ...uint64) (ThumbupStatsParams, error) {
	if len(danmakuIDs) == 0 {
		return ThumbupStatsParams{}, parameterError("ids", "at least one danmaku id is required")
	}
	copyIDs := append([]uint64(nil), danmakuIDs...)
	for _, id := range copyIDs {
		if id == 0 {
			return ThumbupStatsParams{}, parameterError("ids", "danmaku ids must be non-zero")
		}
	}
	return ThumbupStatsParams{oid: oid, ids: copyIDs}, nil
}

func (p ThumbupStatsParams) EncodeQuery() (url.Values, error) {
	if err := p.oid.Validate(); err != nil {
		return nil, err
	}
	values := make([]string, len(p.ids))
	for i, id := range p.ids {
		values[i] = strconv.FormatUint(id, 10)
	}
	return url.Values{"oid": {p.oid.String()}, "ids": {strings.Join(values, ",")}}, nil
}

type AdvStateParams struct{ cid ids.CID }

func NewAdvStateParams(cid ids.CID) AdvStateParams { return AdvStateParams{cid: cid} }

func (p AdvStateParams) EncodeQuery() (url.Values, error) {
	if err := p.cid.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"cid": {p.cid.String()}, "mode": {"sp"}}, nil
}

type WebViewParams struct {
	typeID uint8
	oid    uint64
	pid    *ids.AID
}

func NewWebViewParams(typeID uint8, oid uint64) (WebViewParams, error) {
	if err := validateTypeOID(typeID, oid); err != nil {
		return WebViewParams{}, err
	}
	return WebViewParams{typeID: typeID, oid: oid}, nil
}

func (p WebViewParams) WithAID(aid ids.AID) WebViewParams { p.pid = &aid; return p }

func (p WebViewParams) EncodeQuery() (url.Values, error) {
	if err := validateTypeOID(p.typeID, p.oid); err != nil {
		return nil, err
	}
	values := url.Values{"type": {strconv.FormatUint(uint64(p.typeID), 10)}, "oid": {strconv.FormatUint(p.oid, 10)}}
	if p.pid != nil {
		if err := p.pid.Validate(); err != nil {
			return nil, err
		}
		values.Set("pid", p.pid.String())
	}
	return values, nil
}

type HistoryBytesParams struct {
	typeID uint8
	oid    uint64
	date   string
}

func NewHistoryBytesParams(typeID uint8, oid uint64, date string) (HistoryBytesParams, error) {
	if err := validateTypeOID(typeID, oid); err != nil {
		return HistoryBytesParams{}, err
	}
	if _, err := time.Parse(time.DateOnly, date); err != nil {
		return HistoryBytesParams{}, parameterError("date", "date must use YYYY-MM-DD format")
	}
	return HistoryBytesParams{typeID: typeID, oid: oid, date: date}, nil
}

func (p HistoryBytesParams) EncodeQuery() (url.Values, error) {
	if err := validateTypeOID(p.typeID, p.oid); err != nil {
		return nil, err
	}
	return url.Values{
		"type": {strconv.FormatUint(uint64(p.typeID), 10)},
		"oid":  {strconv.FormatUint(p.oid, 10)},
		"date": {p.date},
	}, nil
}

type XMLListParams struct{ cid ids.CID }

func NewXMLListParams(cid ids.CID) XMLListParams { return XMLListParams{cid: cid} }

func (p XMLListParams) EncodeQuery() (url.Values, error) {
	if err := p.cid.Validate(); err != nil {
		return nil, err
	}
	return url.Values{"oid": {p.cid.String()}}, nil
}

func (p XMLListParams) CommentURL() (string, error) {
	if err := p.cid.Validate(); err != nil {
		return "", err
	}
	return "https://comment.bilibili.com/" + p.cid.String() + ".xml", nil
}

func validateTypeOID(typeID uint8, oid uint64) error {
	if typeID == 0 {
		return parameterError("type", "danmaku type must be non-zero")
	}
	if oid == 0 {
		return parameterError("oid", "danmaku object ID must be non-zero")
	}
	return nil
}

func parameterError(field, message string) error {
	return &bpierr.ParameterError{Field: field, Message: message}
}
