package video

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// PlayURLParams configures the WBI-signed video.play_url endpoint.
type PlayURLParams struct {
	id          identifier
	cid         ids.CID
	quality     *uint64
	formatFlags *uint64
	formatVer   *uint64
	fourK       *bool
	platform    string
	highQuality *bool
	tryLook     *bool
}

// PlayURLByAID identifies a playback request by AV ID and content ID.
func PlayURLByAID(aid ids.AID, cid ids.CID) PlayURLParams {
	return PlayURLParams{id: byAID(aid), cid: cid}
}

// PlayURLByBVID identifies a playback request by BV ID and content ID.
func PlayURLByBVID(bvid ids.BVID, cid ids.CID) PlayURLParams {
	return PlayURLParams{id: byBVID(bvid), cid: cid}
}

// WithQuality returns a copy with a Bilibili quality code.
func (p PlayURLParams) WithQuality(quality uint64) PlayURLParams {
	p.quality = &quality
	return p
}

// WithFormatFlags returns a copy with the Bilibili stream-format bitmask.
func (p PlayURLParams) WithFormatFlags(flags uint64) PlayURLParams {
	p.formatFlags = &flags
	return p
}

// WithFormatVersion returns a copy with the stream-format version.
func (p PlayURLParams) WithFormatVersion(version uint64) PlayURLParams {
	p.formatVer = &version
	return p
}

// With4K returns a copy that controls 4K stream inclusion.
func (p PlayURLParams) With4K(enabled bool) PlayURLParams {
	p.fourK = &enabled
	return p
}

// WithPlatform returns a copy with a non-blank API platform marker.
func (p PlayURLParams) WithPlatform(platform string) (PlayURLParams, error) {
	if strings.TrimSpace(platform) == "" {
		return PlayURLParams{}, &bpierr.ParameterError{Field: "platform", Message: "platform cannot be blank"}
	}
	p.platform = platform
	return p, nil
}

// WithHighQuality returns a copy that controls the high-quality flag.
func (p PlayURLParams) WithHighQuality(enabled bool) PlayURLParams {
	p.highQuality = &enabled
	return p
}

// WithTryLook returns a copy that controls trial playback.
func (p PlayURLParams) WithTryLook(enabled bool) PlayURLParams {
	p.tryLook = &enabled
	return p
}

// EncodeQuery returns the unsigned query. The root VideoClient applies WBI
// signing immediately before sending it.
func (p PlayURLParams) EncodeQuery() (url.Values, error) {
	if err := p.cid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID is invalid"}
	}
	idValues, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	values := url.Values{"cid": {p.cid.String()}}
	if aid := idValues.Get("aid"); aid != "" {
		values.Set("avid", aid)
	} else {
		values.Set("bvid", idValues.Get("bvid"))
	}
	if p.quality != nil {
		values.Set("qn", strconv.FormatUint(*p.quality, 10))
	}
	if p.formatFlags != nil {
		values.Set("fnval", strconv.FormatUint(*p.formatFlags, 10))
	}
	if p.formatVer != nil {
		values.Set("fnver", strconv.FormatUint(*p.formatVer, 10))
	}
	if p.fourK != nil {
		values.Set("fourk", boolFlag(*p.fourK))
	}
	platform := p.platform
	if platform == "" {
		platform = "pc"
	}
	values.Set("platform", platform)
	if p.highQuality != nil {
		values.Set("high_quality", boolFlag(*p.highQuality))
	}
	if p.tryLook != nil {
		values.Set("try_look", boolFlag(*p.tryLook))
	}
	return values, nil
}

func boolFlag(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

// PlayURL is the stable playback payload returned by video.play_url.
type PlayURL struct {
	From               string          `json:"from"`
	Result             string          `json:"result"`
	Message            string          `json:"message"`
	Quality            uint64          `json:"quality"`
	Format             string          `json:"format"`
	TimeLength         uint64          `json:"timelength"`
	AcceptFormat       string          `json:"accept_format"`
	AcceptDescriptions []string        `json:"accept_description"`
	AcceptQualities    []uint64        `json:"accept_quality"`
	VideoCodecID       uint8           `json:"video_codecid"`
	SeekParameter      string          `json:"seek_param"`
	SeekType           string          `json:"seek_type"`
	DURL               []DURL          `json:"durl"`
	DASH               *DASH           `json:"dash"`
	SupportFormats     []SupportFormat `json:"support_formats"`
	HighFormat         json.RawMessage `json:"high_format"`
	LastPlayTime       int64           `json:"last_play_time"`
	LastPlayCID        int64           `json:"last_play_cid"`
}

// DASH contains adaptive video and audio streams.
type DASH struct {
	Video    []DASHStream `json:"video"`
	Audio    []DASHStream `json:"audio"`
	Dolby    *DASHDolby   `json:"dolby"`
	FLAC     *DASHFLAC    `json:"flac"`
	Duration uint64       `json:"duration"`
}

// DASHDolby contains optional Dolby audio streams.
type DASHDolby struct {
	Type  uint8        `json:"type"`
	Audio []DASHStream `json:"audio"`
}

// UnmarshalJSON treats non-numeric Dolby type markers as the protocol's zero
// value. PUGV responses use values such as "NONE", while the regular player
// endpoint returns an integer for the otherwise equivalent field.
func (dolby *DASHDolby) UnmarshalJSON(data []byte) error {
	var value struct {
		Type  json.RawMessage `json:"type"`
		Audio []DASHStream    `json:"audio"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var kind uint8
	if len(value.Type) != 0 {
		_ = json.Unmarshal(value.Type, &kind)
	}
	*dolby = DASHDolby{Type: kind, Audio: value.Audio}
	return nil
}

// DASHFLAC contains FLAC audio streams.
type DASHFLAC struct {
	Audio []DASHStream `json:"audio"`
}

// DASHStream is one adaptive media stream.
type DASHStream struct {
	ID           uint64          `json:"id"`
	BaseURL      string          `json:"baseUrl"`
	BackupURLs   []string        `json:"backupUrl"`
	Bandwidth    uint64          `json:"bandwidth"`
	MIMEType     string          `json:"mimeType"`
	Codecs       string          `json:"codecs"`
	Width        *uint32         `json:"width"`
	Height       *uint32         `json:"height"`
	FrameRate    *string         `json:"frameRate"`
	SAR          *string         `json:"sar"`
	StartWithSAP *uint8          `json:"start_with_sap"`
	SegmentBase  json.RawMessage `json:"segment_base"`
	MD5          *string         `json:"md5"`
	Size         *uint64         `json:"size"`
	DBType       *uint8          `json:"db_type"`
	Type         *string         `json:"type"`
	StreamName   *string         `json:"stream_name"`
	Orientation  *uint8          `json:"orientation"`
}

// UnmarshalJSON accepts the camelCase Web schema and the snake_case PUGV
// schema used by otherwise equivalent DASH stream endpoints.
func (stream *DASHStream) UnmarshalJSON(data []byte) error {
	type alias DASHStream
	var value alias
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var snake struct {
		BaseURL    string   `json:"base_url"`
		BackupURLs []string `json:"backup_url"`
		MIMEType   string   `json:"mime_type"`
		FrameRate  *string  `json:"frame_rate"`
	}
	if err := json.Unmarshal(data, &snake); err != nil {
		return err
	}
	if value.BaseURL == "" {
		value.BaseURL = snake.BaseURL
	}
	if value.BackupURLs == nil {
		value.BackupURLs = snake.BackupURLs
	}
	if value.MIMEType == "" {
		value.MIMEType = snake.MIMEType
	}
	if value.FrameRate == nil {
		value.FrameRate = snake.FrameRate
	}
	*stream = DASHStream(value)
	return nil
}

// DURL is one progressive-download media segment.
type DURL struct {
	Order      uint32   `json:"order"`
	Length     uint64   `json:"length"`
	Size       uint64   `json:"size"`
	Ahead      string   `json:"ahead"`
	VHead      string   `json:"vhead"`
	URL        string   `json:"url"`
	BackupURLs []string `json:"backup_url"`
}

// SupportFormat describes one advertised playback quality/format.
type SupportFormat struct {
	Quality        uint64   `json:"quality"`
	Format         string   `json:"format"`
	NewDescription string   `json:"new_description"`
	Display        string   `json:"display_desc"`
	Superscript    string   `json:"superscript"`
	Codecs         []string `json:"codecs"`
}
