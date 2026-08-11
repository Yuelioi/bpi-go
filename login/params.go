package login

import (
	"net/url"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type NoticeParams struct {
	mid   ids.MID
	buvid string
}

func NewNoticeParams(mid ids.MID) NoticeParams { return NoticeParams{mid: mid} }

func (p NoticeParams) WithBuvid(buvid string) (NoticeParams, error) {
	value, err := nonBlank("buvid", buvid)
	if err != nil {
		return NoticeParams{}, err
	}
	p.buvid = value
	return p, nil
}

func (p NoticeParams) EncodeQuery() (url.Values, error) {
	if err := p.mid.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{"mid": {p.mid.String()}}
	if p.buvid != "" {
		values.Set("buvid", p.buvid)
	}
	return values, nil
}

type LogParams struct {
	jsonp       string
	webLocation string
}

func NewLogParams() LogParams { return LogParams{jsonp: "jsonp", webLocation: "333.33"} }

func (p LogParams) WithJSONP(jsonp string) (LogParams, error) {
	value, err := nonBlank("jsonp", jsonp)
	if err != nil {
		return LogParams{}, err
	}
	p.jsonp = value
	return p, nil
}

func (p LogParams) WithWebLocation(location string) (LogParams, error) {
	value, err := nonBlank("web_location", location)
	if err != nil {
		return LogParams{}, err
	}
	p.webLocation = value
	return p, nil
}

func (p LogParams) EncodeQuery() (url.Values, error) {
	return url.Values{"jsonp": {p.jsonp}, "web_location": {p.webLocation}}, nil
}

type QRPollParams struct{ key string }

func NewQRPollParams(key string) (QRPollParams, error) {
	value, err := nonBlank("qrcode_key", key)
	if err != nil {
		return QRPollParams{}, err
	}
	return QRPollParams{key: value}, nil
}

func (p QRPollParams) EncodeQuery() (url.Values, error) {
	return url.Values{"qrcode_key": {p.key}}, nil
}

func nonBlank(field, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", &bpierr.ParameterError{Field: field, Message: "value cannot be blank"}
	}
	return value, nil
}
