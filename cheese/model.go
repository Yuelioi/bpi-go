package cheese

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/video"
)

type Course struct {
	Brief             Brief           `json:"brief"`
	Cover             string          `json:"cover"`
	EpisodePage       EpisodePage     `json:"episode_page"`
	EpisodeSort       int32           `json:"episode_sort"`
	Episodes          []Episode       `json:"episodes"`
	ReleaseBottomInfo string          `json:"release_bottom_info"`
	ReleaseInfo       string          `json:"release_info"`
	ReleaseInfo2      string          `json:"release_info2"`
	ReleaseStatus     string          `json:"release_status"`
	SeasonID          uint64          `json:"season_id"`
	ShareURL          string          `json:"share_url"`
	ShortLink         string          `json:"short_link"`
	Stat              CourseStat      `json:"stat"`
	Status            int32           `json:"status"`
	Subtitle          string          `json:"subtitle"`
	Title             string          `json:"title"`
	Uploader          Uploader        `json:"up_info"`
	UserStatus        UserStatus      `json:"user_status"`
	Coupon            json.RawMessage `json:"coupon"`
	FAQ               json.RawMessage `json:"faq"`
	Payment           json.RawMessage `json:"payment"`
	PurchaseNote      json.RawMessage `json:"purchase_note"`
	PurchaseProtocol  json.RawMessage `json:"purchase_protocol"`
}

type Brief struct {
	Content string       `json:"content"`
	Images  []BriefImage `json:"img"`
	Title   string       `json:"title"`
	Type    int32        `json:"type"`
}

type BriefImage struct {
	AspectRatio float64 `json:"aspect_ratio"`
	URL         string  `json:"url"`
}

type EpisodePage struct {
	Next  bool   `json:"next"`
	Page  uint32 `json:"num"`
	Size  uint32 `json:"size"`
	Total uint32 `json:"total"`
}

type Episode struct {
	AID            uint64 `json:"aid"`
	CID            uint64 `json:"cid"`
	Cover          string `json:"cover"`
	Duration       uint64 `json:"duration"`
	From           string `json:"from"`
	ID             uint64 `json:"id"`
	Index          uint32 `json:"index"`
	Page           uint32 `json:"page"`
	Play           uint64 `json:"play"`
	Playable       bool   `json:"playable"`
	ReleaseDate    uint64 `json:"release_date"`
	Status         int32  `json:"status"`
	Subtitle       string `json:"subtitle"`
	Title          string `json:"title"`
	Watched        bool   `json:"watched"`
	WatchedHistory uint64 `json:"watchedHistory"`
}

type CourseStat struct {
	Play            uint64 `json:"play"`
	PlayDescription string `json:"play_desc"`
}

type Uploader struct {
	Avatar   string `json:"avatar"`
	Brief    string `json:"brief"`
	Follower uint64 `json:"follower"`
	IsFollow int32  `json:"is_follow"`
	Link     string `json:"link"`
	MID      uint64 `json:"mid"`
	UserName string `json:"uname"`
}

type UserStatus struct {
	Favored      int32     `json:"favored"`
	FavoredCount uint64    `json:"favored_count"`
	Paid         int32     `json:"payed"`
	Progress     *Progress `json:"progress"`
}

type Progress struct {
	LastEpisodeID    uint64 `json:"last_ep_id"`
	LastEpisodeIndex string `json:"last_ep_index"`
	LastTime         uint64 `json:"last_time"`
}

type EpisodeList struct {
	Items []Episode   `json:"items"`
	Page  EpisodePage `json:"page"`
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
	DURLs              []video.DURL          `json:"durls"`
	DASH               *video.DASH           `json:"dash"`
	SupportFormats     []video.SupportFormat `json:"support_formats"`
	Code               int32                 `json:"code"`
	FormatFlags        uint64                `json:"fnval"`
	FormatVersion      uint64                `json:"fnver"`
	VideoProject       bool                  `json:"video_project"`
	Type               string                `json:"type"`
	HasPaid            bool                  `json:"has_paid"`
	IsPreview          *uint32               `json:"is_preview"`
	NoRecode           int32                 `json:"no_rexcode"`
	Status             int32                 `json:"status"`
	Fragments          []FragmentVideo       `json:"fragment_videos"`
	Volume             *Volume               `json:"volume"`
}

type FragmentVideo struct {
	Info           FragmentInfo `json:"fragment_info"`
	PlayableStatus bool         `json:"playable_status"`
	VideoInfo      VideoInfo    `json:"video_info"`
}

type FragmentInfo struct {
	Type     string `json:"fragment_type"`
	Index    int64  `json:"index"`
	AID      int64  `json:"aid"`
	Position string `json:"fragment_position"`
	CID      int64  `json:"cid"`
}

type VideoInfo struct {
	NoRecode      int64           `json:"no_rexcode"`
	FormatFlags   int64           `json:"fnval"`
	VideoProject  bool            `json:"video_project"`
	ExpireTime    int64           `json:"expire_time"`
	FormatVersion int64           `json:"fnver"`
	Type          string          `json:"type"`
	URL           string          `json:"url"`
	Quality       int64           `json:"quality"`
	TimeLength    int64           `json:"timelength"`
	DASH          video.DASH      `json:"dash"`
	VideoCodecID  int64           `json:"video_codecid"`
	CID           int64           `json:"cid"`
	FileInfo      json.RawMessage `json:"file_info"`
}

type Volume struct {
	MeasuredI         float64         `json:"measured_i"`
	TargetI           float64         `json:"target_i"`
	TargetOffset      float64         `json:"target_offset"`
	MeasuredLRA       float64         `json:"measured_lra"`
	TargetTP          float64         `json:"target_tp"`
	MeasuredTP        float64         `json:"measured_tp"`
	MeasuredThreshold float64         `json:"measured_threshold"`
	MultiSceneArgs    json.RawMessage `json:"multi_scene_args"`
}
