package video

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// PlayerInfoParams identifies one video part for video.player_info_v2.
type PlayerInfoParams struct {
	id        identifier
	cid       ids.CID
	seasonID  ids.SeasonID
	episodeID ids.EpisodeID
}

func PlayerInfoByAID(aid ids.AID, cid ids.CID) PlayerInfoParams {
	return PlayerInfoParams{id: byAID(aid), cid: cid}
}

func PlayerInfoByBVID(bvid ids.BVID, cid ids.CID) PlayerInfoParams {
	return PlayerInfoParams{id: byBVID(bvid), cid: cid}
}

func (p PlayerInfoParams) WithSeasonID(seasonID ids.SeasonID) PlayerInfoParams {
	p.seasonID = seasonID
	return p
}

func (p PlayerInfoParams) WithEpisodeID(episodeID ids.EpisodeID) PlayerInfoParams {
	p.episodeID = episodeID
	return p
}

func (p PlayerInfoParams) EncodeQuery() (url.Values, error) {
	values, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	if err := p.cid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID must be non-zero"}
	}
	values.Set("cid", p.cid.String())
	if p.seasonID != 0 {
		if err := p.seasonID.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "season_id", Message: "season ID must be non-zero"}
		}
		values.Set("season_id", p.seasonID.String())
	}
	if p.episodeID != 0 {
		if err := p.episodeID.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "ep_id", Message: "episode ID must be non-zero"}
		}
		values.Set("ep_id", p.episodeID.String())
	}
	return values, nil
}

// HomepageRecommendationsParams configures video.homepage_recommendations.
// Its zero value uses the promoted Web defaults.
type HomepageRecommendationsParams struct {
	pageSize   uint8
	freshIndex uint32
	fetchRow   uint32
}

func NewHomepageRecommendationsParams() HomepageRecommendationsParams {
	return HomepageRecommendationsParams{}
}

func (p HomepageRecommendationsParams) WithPageSize(size uint8) (HomepageRecommendationsParams, error) {
	if size == 0 || size > 30 {
		return HomepageRecommendationsParams{}, &bpierr.ParameterError{Field: "ps", Message: "value must be between 1 and 30"}
	}
	p.pageSize = size
	return p, nil
}

func (p HomepageRecommendationsParams) WithFreshIndex(index uint32) (HomepageRecommendationsParams, error) {
	if index == 0 {
		return HomepageRecommendationsParams{}, &bpierr.ParameterError{Field: "fresh_idx", Message: "value must be non-zero"}
	}
	p.freshIndex = index
	return p, nil
}

func (p HomepageRecommendationsParams) WithFetchRow(row uint32) (HomepageRecommendationsParams, error) {
	if row == 0 {
		return HomepageRecommendationsParams{}, &bpierr.ParameterError{Field: "fetch_row", Message: "value must be non-zero"}
	}
	p.fetchRow = row
	return p, nil
}

func (p HomepageRecommendationsParams) EncodeQuery() (url.Values, error) {
	pageSize := p.pageSize
	if pageSize == 0 {
		pageSize = 12
	}
	freshIndex := p.freshIndex
	if freshIndex == 0 {
		freshIndex = 1
	}
	fetchRow := p.fetchRow
	if fetchRow == 0 {
		fetchRow = 1
	}
	index := strconv.FormatUint(uint64(freshIndex), 10)
	return url.Values{
		"fresh_type":   {"4"},
		"ps":           {strconv.FormatUint(uint64(pageSize), 10)},
		"fresh_idx":    {index},
		"fresh_idx_1h": {index},
		"brush":        {index},
		"fetch_row":    {strconv.FormatUint(uint64(fetchRow), 10)},
	}, nil
}

// AISummaryParams identifies a video, content part, and uploader for
// video.ai_summary.
type AISummaryParams struct {
	id    identifier
	cid   ids.CID
	upMID ids.MID
}

func AISummaryByAID(aid ids.AID, cid ids.CID, upMID ids.MID) AISummaryParams {
	return AISummaryParams{id: byAID(aid), cid: cid, upMID: upMID}
}

func AISummaryByBVID(bvid ids.BVID, cid ids.CID, upMID ids.MID) AISummaryParams {
	return AISummaryParams{id: byBVID(bvid), cid: cid, upMID: upMID}
}

func (p AISummaryParams) EncodeQuery() (url.Values, error) {
	values, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	if err := p.cid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID must be non-zero"}
	}
	if err := p.upMID.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "up_mid", Message: "uploader member ID must be non-zero"}
	}
	values.Set("cid", p.cid.String())
	values.Set("up_mid", p.upMID.String())
	return values, nil
}

type PlayerInfo struct {
	AID               ids.AID            `json:"aid"`
	BVID              ids.BVID           `json:"bvid"`
	AllowBP           bool               `json:"allow_bp"`
	NoShare           bool               `json:"no_share"`
	CID               ids.CID            `json:"cid"`
	DMMask            *DMMask            `json:"dm_mask"`
	Subtitle          *SubtitleInfo      `json:"subtitle"`
	ViewPoints        []ViewPoint        `json:"view_points"`
	IPInfo            json.RawMessage    `json:"ip_info"`
	LoginMID          uint64             `json:"login_mid"`
	LoginMIDHash      *string            `json:"login_mid_hash"`
	IsOwner           bool               `json:"is_owner"`
	Name              string             `json:"name"`
	Permission        string             `json:"permission"`
	LevelInfo         json.RawMessage    `json:"level_info"`
	VIP               json.RawMessage    `json:"vip"`
	AnswerStatus      uint8              `json:"answer_status"`
	BlockTime         uint64             `json:"block_time"`
	Role              string             `json:"role"`
	LastPlayTime      int64              `json:"last_play_time"`
	LastPlayCID       int64              `json:"last_play_cid"`
	NowTime           uint64             `json:"now_time"`
	OnlineCount       *uint64            `json:"online_count"`
	NeedLoginSubtitle bool               `json:"need_login_subtitle"`
	PreviewToast      string             `json:"preview_toast"`
	Interaction       *Interaction       `json:"interaction"`
	Options           *PlayerOptions     `json:"options"`
	GuideAttention    json.RawMessage    `json:"guide_attention"`
	JumpCard          json.RawMessage    `json:"jump_card"`
	OperationCard     json.RawMessage    `json:"operation_card"`
	OnlineSwitch      json.RawMessage    `json:"online_switch"`
	Fawkes            json.RawMessage    `json:"fawkes"`
	ShowSwitch        json.RawMessage    `json:"show_switch"`
	BGMInfo           *BGMInfo           `json:"bgm_info"`
	ToastBlock        bool               `json:"toast_block"`
	IsUPowerExclusive bool               `json:"is_upower_exclusive"`
	IsUPowerPlay      bool               `json:"is_upower_play"`
	IsUGCPayPreview   bool               `json:"is_ugc_pay_preview"`
	ElectricHighLevel *ElectricHighLevel `json:"elec_high_level"`
	DisableShowUPInfo bool               `json:"disable_show_up_info"`
}

type DMMask struct {
	CID      uint64 `json:"cid"`
	Platform uint8  `json:"plat"`
	FPS      uint64 `json:"fps"`
	Time     uint64 `json:"time"`
	MaskURL  string `json:"mask_url"`
}

type SubtitleInfo struct {
	AllowSubmit bool           `json:"allow_submit"`
	Language    string         `json:"lan"`
	LanguageDoc string         `json:"lan_doc"`
	Subtitles   []SubtitleItem `json:"subtitles"`
}

type SubtitleItem struct {
	AIStatus    uint8  `json:"ai_status"`
	AIType      uint8  `json:"ai_type"`
	ID          uint64 `json:"id"`
	IDString    string `json:"id_str"`
	IsLocked    bool   `json:"is_lock"`
	Language    string `json:"lan"`
	LanguageDoc string `json:"lan_doc"`
	URL         string `json:"subtitle_url"`
	Type        uint8  `json:"type"`
}

type ViewPoint struct {
	Content  string `json:"content"`
	From     uint64 `json:"from"`
	To       uint64 `json:"to"`
	Type     uint8  `json:"type"`
	ImageURL string `json:"imgUrl"`
	LogoURL  string `json:"logoUrl"`
	TeamType string `json:"team_type"`
	TeamName string `json:"team_name"`
}

type Interaction struct {
	GraphVersion uint64  `json:"graph_version"`
	Message      *string `json:"msg"`
	ErrorToast   *string `json:"error_toast"`
	Mark         *uint8  `json:"mark"`
	NeedReload   *uint8  `json:"need_reload"`
}

type PlayerOptions struct {
	Is360      bool `json:"is_360"`
	WithoutVIP bool `json:"without_vip"`
}

type BGMInfo struct {
	MusicID    string `json:"music_id"`
	MusicTitle string `json:"music_title"`
	JumpURL    string `json:"jump_url"`
}

type ElectricHighLevel struct {
	PrivilegeType uint64          `json:"privilege_type"`
	Title         string          `json:"title"`
	Subtitle      string          `json:"sub_title"`
	ShowButton    bool            `json:"show_button"`
	ButtonText    string          `json:"button_text"`
	JumpURL       json.RawMessage `json:"jump_url"`
	Intro         string          `json:"intro"`
	Open          bool            `json:"open"`
	New           bool            `json:"new"`
	QuestionText  string          `json:"question_text"`
	QADetailLink  string          `json:"qa_detail_link"`
}

type RecommendationReason struct {
	Type    uint8   `json:"reason_type"`
	Content *string `json:"content"`
}

type HomeRecommendationStat struct {
	View    uint64 `json:"view"`
	Danmaku uint64 `json:"danmaku"`
	Like    uint64 `json:"like"`
}

type HomepageRecommendation struct {
	AVFeature    json.RawMessage         `json:"av_feature"`
	BusinessInfo json.RawMessage         `json:"business_info"`
	BVID         string                  `json:"bvid"`
	CID          uint64                  `json:"cid"`
	Duration     uint64                  `json:"duration"`
	Goto         string                  `json:"goto"`
	ID           uint64                  `json:"id"`
	IsFollowed   uint8                   `json:"is_followed"`
	IsStock      uint8                   `json:"is_stock"`
	Owner        Owner                   `json:"owner"`
	Picture      string                  `json:"pic"`
	Position     uint8                   `json:"pos"`
	PublishTime  uint64                  `json:"pubdate"`
	Reason       *RecommendationReason   `json:"rcmd_reason"`
	RoomInfo     json.RawMessage         `json:"room_info"`
	ShowInfo     uint8                   `json:"show_info"`
	Stat         *HomeRecommendationStat `json:"stat"`
	Title        string                  `json:"title"`
	TrackID      string                  `json:"track_id"`
	URI          string                  `json:"uri"`
}

type HomepageRecommendations struct {
	Items                     []HomepageRecommendation `json:"item"`
	MID                       uint64                   `json:"mid"`
	PreloadExposePercent      float32                  `json:"preload_expose_pct"`
	PreloadFloorExposePercent float32                  `json:"preload_floor_expose_pct"`
}

type AISummary struct {
	Code         int8             `json:"code"`
	ModelResult  *AISummaryResult `json:"model_result"`
	STID         *string          `json:"stid"`
	Status       *uint8           `json:"status"`
	LikeCount    uint64           `json:"like_num"`
	DislikeCount uint64           `json:"dislike_num"`
}

type AISummaryResult struct {
	Type    uint8              `json:"result_type"`
	Summary string             `json:"summary"`
	Outline []AISummaryOutline `json:"outline"`
}

type AISummaryOutline struct {
	Title     string                 `json:"title"`
	Parts     []AISummaryOutlinePart `json:"part_outline"`
	Timestamp uint64                 `json:"timestamp"`
}

type AISummaryOutlinePart struct {
	Timestamp uint64 `json:"timestamp"`
	Content   string `json:"content"`
}
