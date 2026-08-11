// Package misc contains validated parameters and stable response models for
// Bilibili session bootstrap and utility endpoints.
package misc

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

const (
	DefaultPlatform     = "unix"
	DefaultShareChannel = "COPY"
	DefaultShareID      = "main.ugc-video-detail.0.0.pv"
	DefaultShareMode    = uint32(4)
	DefaultBuvid        = "qwq"
	DefaultBuild        = uint64(6_114_514)
)

type ShortLinkParams struct {
	aid          ids.AID
	platform     string
	shareChannel string
	shareID      string
	shareMode    uint32
	buvid        string
	build        uint64
}

func NewShortLinkParams(aid ids.AID) ShortLinkParams {
	return ShortLinkParams{
		aid: aid, platform: DefaultPlatform, shareChannel: DefaultShareChannel,
		shareID: DefaultShareID, shareMode: DefaultShareMode, buvid: DefaultBuvid, build: DefaultBuild,
	}
}

func (p ShortLinkParams) WithPlatform(value string) (ShortLinkParams, error) {
	return p.withString("platform", value, func(value string) { p.platform = value })
}

func (p ShortLinkParams) WithShareChannel(value string) (ShortLinkParams, error) {
	return p.withString("share_channel", value, func(value string) { p.shareChannel = value })
}

func (p ShortLinkParams) WithShareID(value string) (ShortLinkParams, error) {
	return p.withString("share_id", value, func(value string) { p.shareID = value })
}

func (p ShortLinkParams) WithBuvid(value string) (ShortLinkParams, error) {
	return p.withString("buvid", value, func(value string) { p.buvid = value })
}

func (p ShortLinkParams) WithShareMode(value uint32) ShortLinkParams {
	p.shareMode = value
	return p
}

func (p ShortLinkParams) WithBuild(value uint64) ShortLinkParams {
	p.build = value
	return p
}

func (p ShortLinkParams) EncodeForm() (url.Values, error) {
	if err := p.aid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "oid", Message: "video ID is invalid"}
	}
	for field, value := range map[string]string{"platform": p.platform, "share_channel": p.shareChannel, "share_id": p.shareID, "buvid": p.buvid} {
		if strings.TrimSpace(value) == "" {
			return nil, &bpierr.ParameterError{Field: field, Message: "value cannot be blank"}
		}
	}
	return url.Values{
		"platform":      {p.platform},
		"share_channel": {p.shareChannel},
		"share_id":      {p.shareID},
		"share_mode":    {strconv.FormatUint(uint64(p.shareMode), 10)},
		"oid":           {p.aid.String()},
		"buvid":         {p.buvid},
		"build":         {strconv.FormatUint(p.build, 10)},
	}, nil
}

func (p ShortLinkParams) withString(field, value string, set func(string)) (ShortLinkParams, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return ShortLinkParams{}, &bpierr.ParameterError{Field: field, Message: "value cannot be blank"}
	}
	set(value)
	return p, nil
}

type Buvid3 struct {
	Buvid string `json:"buvid"`
}

type Buvid struct {
	Buvid3 string `json:"b_3"`
	Buvid4 string `json:"b_4"`
}

type ShortLink struct {
	Content string `json:"content"`
	Count   int32  `json:"count"`
	Link    string `json:"-"`
	Title   string `json:"-"`
}

func (link *ShortLink) Extract() {
	const prefix = "https://b23.tv/"
	position := strings.Index(link.Content, prefix)
	if position < 0 {
		link.Title = strings.TrimSpace(link.Content)
		link.Link = ""
		return
	}
	link.Title = strings.TrimSpace(link.Content[:position])
	link.Link = strings.TrimSpace(link.Content[position:])
}

type Ticket struct {
	Ticket     string           `json:"ticket"`
	CreatedAt  int64            `json:"created_at"`
	TTL        int32            `json:"ttl"`
	Context    map[string]any   `json:"context"`
	Navigation TicketNavigation `json:"nav"`
}

type TicketNavigation struct {
	Image string `json:"img"`
	Sub   string `json:"sub"`
}
