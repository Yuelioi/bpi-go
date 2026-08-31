package bangumi

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/video"
)

type Info struct {
	Media  Media   `json:"media"`
	Review *Review `json:"review"`
}

type Media struct {
	Areas             []Area     `json:"areas"`
	Cover             string     `json:"cover"`
	HorizontalPicture string     `json:"horizontal_picture"`
	MediaID           uint64     `json:"media_id"`
	NewEpisode        NewEpisode `json:"new_ep"`
	Rating            Rating     `json:"rating"`
	SeasonID          uint64     `json:"season_id"`
	ShareURL          string     `json:"share_url"`
	Title             string     `json:"title"`
	Type              uint32     `json:"type"`
	TypeName          string     `json:"type_name"`
}

type Area struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type NewEpisode struct {
	ID        uint64 `json:"id"`
	Index     string `json:"index"`
	IndexShow string `json:"index_show"`
}

type Rating struct {
	Count uint64  `json:"count"`
	Score float64 `json:"score"`
}

type Review struct {
	IsCoin uint32 `json:"is_coin"`
	IsOpen uint32 `json:"is_open"`
}

type Detail struct {
	Actors                string           `json:"actors"`
	Alias                 string           `json:"alias"`
	Areas                 []Area           `json:"areas"`
	BackgroundCover       string           `json:"bkg_cover"`
	Cover                 string           `json:"cover"`
	DeliveryFragmentVideo bool             `json:"delivery_fragment_video"`
	EnableVT              bool             `json:"enable_vt"`
	Episodes              []Episode        `json:"episodes"`
	Evaluate              string           `json:"evaluate"`
	JapaneseTitle         string           `json:"jp_title"`
	Link                  string           `json:"link"`
	MediaID               uint64           `json:"media_id"`
	Mode                  uint32           `json:"mode"`
	NewEpisode            DetailNewEpisode `json:"new_ep"`
	Publish               Publish          `json:"publish"`
	Rating                *Rating          `json:"rating"`
	Record                string           `json:"record"`
	Rights                Rights           `json:"rights"`
	SeasonID              uint64           `json:"season_id"`
	SeasonTitle           string           `json:"season_title"`
	Seasons               []Season         `json:"seasons"`
	ShareCopy             string           `json:"share_copy"`
	ShareSubtitle         string           `json:"share_sub_title"`
	ShareURL              string           `json:"share_url"`
	ShowSeasonType        uint32           `json:"show_season_type"`
	SquareCover           string           `json:"square_cover"`
	Staff                 string           `json:"staff"`
	Stat                  Stat             `json:"stat"`
	Status                uint32           `json:"status"`
	Styles                []string         `json:"styles"`
	Subtitle              string           `json:"subtitle"`
	Title                 string           `json:"title"`
	Total                 int32            `json:"total"`
	Type                  uint32           `json:"type"`
	UserStatus            *UserStatus      `json:"user_status"`
	Activity              *json.RawMessage `json:"activity"`
	Payment               *json.RawMessage `json:"payment"`
	PlayStrategy          *json.RawMessage `json:"play_strategy"`
}

type DetailNewEpisode struct {
	ID    uint64 `json:"id"`
	Desc  string `json:"desc"`
	IsNew uint32 `json:"is_new"`
	Title string `json:"title"`
}

type Episode struct {
	AID                uint64           `json:"aid"`
	Badge              string           `json:"badge"`
	BadgeInfo          *BadgeInfo       `json:"badge_info"`
	BadgeType          uint32           `json:"badge_type"`
	BVID               string           `json:"bvid"`
	CID                uint64           `json:"cid"`
	Cover              string           `json:"cover"`
	Dimension          *Dimension       `json:"dimension"`
	Duration           uint64           `json:"duration"`
	EnableVT           bool             `json:"enable_vt"`
	EpisodeID          uint64           `json:"ep_id"`
	From               string           `json:"from"`
	ID                 uint64           `json:"id"`
	IsViewHidden       bool             `json:"is_view_hide"`
	Link               string           `json:"link"`
	LongTitle          string           `json:"long_title"`
	PublishTime        uint64           `json:"pub_time"`
	PV                 uint64           `json:"pv"`
	ReleaseDate        string           `json:"release_date"`
	Rights             *EpisodeRights   `json:"rights"`
	SectionType        uint32           `json:"section_type"`
	ShareCopy          string           `json:"share_copy"`
	ShareURL           string           `json:"share_url"`
	ShortLink          string           `json:"short_link"`
	ShowTitle          string           `json:"show_title"`
	ShowDRMLoginDialog bool             `json:"showDrmLoginDialog"`
	Skip               *json.RawMessage `json:"skip"`
	Status             uint32           `json:"status"`
	Subtitle           string           `json:"subtitle"`
	Title              string           `json:"title"`
	VID                string           `json:"vid"`
}

type BadgeInfo struct {
	Background      string `json:"bg_color"`
	NightBackground string `json:"bg_color_night"`
	Text            string `json:"text"`
}

type Dimension struct {
	Height uint32 `json:"height"`
	Rotate uint32 `json:"rotate"`
	Width  uint32 `json:"width"`
}

type EpisodeRights struct {
	AllowDanmaku  uint32  `json:"allow_dm"`
	AllowDownload uint32  `json:"allow_download"`
	AreaLimit     uint32  `json:"area_limit"`
	AllowDemand   *uint32 `json:"allow_demand"`
}

type Publish struct {
	IsFinish       uint32 `json:"is_finish"`
	IsStarted      uint32 `json:"is_started"`
	PublishTime    string `json:"pub_time"`
	PublishDisplay string `json:"pub_time_show"`
	UnknownDate    uint32 `json:"unknow_pub_date"`
	Weekday        uint32 `json:"weekday"`
}

type Rights struct {
	AllowDownload   uint32 `json:"allow_download"`
	AllowReview     uint32 `json:"allow_review"`
	AreaLimit       uint32 `json:"area_limit"`
	CanWatch        uint32 `json:"can_watch"`
	Copyright       string `json:"copyright"`
	IsPreview       uint32 `json:"is_preview"`
	OnlyVIPDownload uint32 `json:"only_vip_download"`
	Resource        string `json:"resource"`
}

type Season struct {
	Badge       string     `json:"badge"`
	BadgeInfo   *BadgeInfo `json:"badge_info"`
	BadgeType   uint32     `json:"badge_type"`
	Cover       string     `json:"cover"`
	EnableVT    bool       `json:"enable_vt"`
	MediaID     uint64     `json:"media_id"`
	SeasonID    uint64     `json:"season_id"`
	SeasonTitle string     `json:"season_title"`
	SeasonType  uint32     `json:"season_type"`
}

type Stat struct {
	Coins      uint64  `json:"coins"`
	Danmakus   uint64  `json:"danmakus"`
	Favorite   uint64  `json:"favorite"`
	Favorites  uint64  `json:"favorites"`
	FollowText string  `json:"follow_text"`
	Hot        *uint64 `json:"hot"`
	Likes      uint64  `json:"likes"`
	Reply      uint64  `json:"reply"`
	Share      uint64  `json:"share"`
	Views      uint64  `json:"views"`
	VT         uint64  `json:"vt"`
}

// UnmarshalJSON accepts both the full season-stat names and the compact
// section-episode aliases observed in historical responses.
func (stat *Stat) UnmarshalJSON(data []byte) error {
	type statAlias Stat
	var value statAlias
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["coins"]; !present {
		if raw, ok := fields["coin"]; ok {
			if err := json.Unmarshal(raw, &value.Coins); err != nil {
				return err
			}
		}
	}
	if _, present := fields["views"]; !present {
		if raw, ok := fields["play"]; ok {
			if err := json.Unmarshal(raw, &value.Views); err != nil {
				return err
			}
		}
	}
	*stat = Stat(value)
	return nil
}

type UserStatus struct {
	AreaLimit   uint32 `json:"area_limit"`
	Follow      uint32 `json:"follow"`
	FollowState uint32 `json:"follow_status"`
	Login       uint32 `json:"login"`
	Pay         uint32 `json:"pay"`
	Sponsor     uint32 `json:"sponsor"`
}

type Sections struct {
	Main     MainSection   `json:"main_section"`
	Sections []MainSection `json:"section"`
}

type MainSection struct {
	Episodes []SectionEpisode `json:"episodes"`
	ID       uint64           `json:"id"`
	Type     uint32           `json:"type"`
	Title    string           `json:"title"`
}

type SectionEpisode struct {
	AID        uint64    `json:"aid"`
	Badge      string    `json:"badge"`
	BadgeInfo  BadgeInfo `json:"badge_info"`
	BadgeType  uint32    `json:"badge_type"`
	CID        uint64    `json:"cid"`
	Cover      string    `json:"cover"`
	From       string    `json:"from"`
	ID         uint64    `json:"id"`
	IsPremiere uint32    `json:"is_premiere"`
	LongTitle  string    `json:"long_title"`
	ShareURL   string    `json:"share_url"`
	Status     uint32    `json:"status"`
	Title      string    `json:"title"`
	VID        string    `json:"vid"`
}

type PlayURL struct {
	From               string                `json:"from"`
	Result             string                `json:"result"`
	Message            string                `json:"message"`
	Quality            uint64                `json:"quality"`
	Format             string                `json:"format"`
	TimeLength         uint64                `json:"timelength"`
	AcceptFormat       string                `json:"accept_format"`
	AcceptDescriptions []string              `json:"accept_description"`
	AcceptQualities    []uint64              `json:"accept_quality"`
	VideoCodecID       uint8                 `json:"video_codecid"`
	SeekParameter      string                `json:"seek_param"`
	SeekType           string                `json:"seek_type"`
	DURLs              []video.DURL          `json:"durl"`
	DASH               *video.DASH           `json:"dash"`
	SupportFormats     []video.SupportFormat `json:"support_formats"`
	Code               uint32                `json:"code"`
	FormatVersion      uint32                `json:"fnver"`
	VideoProject       bool                  `json:"video_project"`
	Type               string                `json:"type"`
	BP                 uint32                `json:"bp"`
	VIPType            *uint32               `json:"vip_type"`
	VIPStatus          *uint32               `json:"vip_status"`
	IsDRM              bool                  `json:"is_drm"`
	NoRecode           uint32                `json:"no_rexcode"`
	RecordInfo         *RecordInfo           `json:"record_info"`
}

type RecordInfo struct {
	Icon   string `json:"record_icon"`
	Record string `json:"record"`
}
