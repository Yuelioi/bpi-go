package user

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/ids"
)

type CardProfile struct {
	Card         CardSummary `json:"card"`
	Following    bool        `json:"following"`
	ArchiveCount uint64      `json:"archive_count"`
	ArticleCount uint64      `json:"article_count"`
	Follower     uint64      `json:"follower"`
	LikeCount    uint64      `json:"like_num"`
}

type CardSummary struct {
	MID       ids.MID `json:"mid"`
	Name      string  `json:"name"`
	Sex       *string `json:"sex"`
	Face      string  `json:"face"`
	Sign      string  `json:"sign"`
	Fans      uint64  `json:"fans"`
	Attention uint64  `json:"attention"`
}

func (summary *CardSummary) UnmarshalJSON(data []byte) error {
	type alias CardSummary
	var value alias
	auxiliary := struct {
		MID flexibleMID `json:"mid"`
		*alias
	}{alias: &value}
	if err := json.Unmarshal(data, &auxiliary); err != nil {
		return err
	}
	value.MID = ids.MID(auxiliary.MID)
	*summary = CardSummary(value)
	return nil
}

type BatchCard struct {
	MID     ids.MID `json:"mid"`
	Name    string  `json:"name"`
	Face    string  `json:"face"`
	Sign    string  `json:"sign"`
	Rank    int32   `json:"rank"`
	Level   int32   `json:"level"`
	Silence int32   `json:"silence"`
}

type BatchInfo struct {
	MID           ids.MID          `json:"mid"`
	Name          string           `json:"name"`
	Sign          string           `json:"sign"`
	Rank          int32            `json:"rank"`
	Level         int32            `json:"level"`
	Silence       int32            `json:"silence"`
	Sex           *string          `json:"sex"`
	Face          string           `json:"face"`
	VIP           *BatchVIP        `json:"vip"`
	Official      *Official        `json:"official"`
	IsFakeAccount *uint32          `json:"is_fake_account"`
	ExpertInfo    *json.RawMessage `json:"expert_info"`
}

type BatchVIP struct {
	Type      int32 `json:"type"`
	Status    int32 `json:"status"`
	DueDate   int64 `json:"due_date"`
	PayType   int32 `json:"vip_pay_type"`
	ThemeType int32 `json:"theme_type"`
}

type SpaceProfile struct {
	MID        ids.MID        `json:"mid"`
	Name       string         `json:"name"`
	Sex        *string        `json:"sex"`
	Face       string         `json:"face"`
	Sign       string         `json:"sign"`
	Level      uint8          `json:"level"`
	Silence    uint8          `json:"silence"`
	Coins      float64        `json:"coins"`
	FansBadge  bool           `json:"fans_badge"`
	IsFollowed bool           `json:"is_followed"`
	TopPhoto   *string        `json:"top_photo"`
	Official   *Official      `json:"official"`
	VIP        *VIP           `json:"vip"`
	LiveRoom   *SpaceLiveRoom `json:"live_room"`
}

type Official struct {
	Role  int32  `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  int32  `json:"type"`
}

type VIP struct {
	Type   int32 `json:"type"`
	Status int32 `json:"status"`
}

type SpaceLiveRoom struct {
	RoomStatus uint8  `json:"roomStatus"`
	LiveStatus uint8  `json:"liveStatus"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	RoomID     uint64 `json:"roomid"`
}

// SpaceNotice wraps the string payload returned by /x/space/notice.
type SpaceNotice struct{ Content string }

func (notice *SpaceNotice) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &notice.Content)
}

func (notice SpaceNotice) MarshalJSON() ([]byte, error) { return json.Marshal(notice.Content) }

type BangumiFollowList struct {
	Items    []BangumiSeason `json:"list"`
	Page     uint32          `json:"pn"`
	PageSize uint32          `json:"ps"`
	Total    uint64          `json:"total"`
}

type BangumiSeason struct {
	SeasonID       int64          `json:"season_id"`
	MediaID        int64          `json:"media_id"`
	SeasonType     int64          `json:"season_type"`
	SeasonTypeName string         `json:"season_type_name"`
	Title          string         `json:"title"`
	Cover          string         `json:"cover"`
	TotalCount     int64          `json:"total_count"`
	IsFinish       int64          `json:"is_finish"`
	IsStarted      int64          `json:"is_started"`
	IsPlay         int64          `json:"is_play"`
	Badge          string         `json:"badge"`
	BadgeType      int64          `json:"badge_type"`
	LatestEpisode  BangumiEpisode `json:"new_ep"`
	Rating         *BangumiRating `json:"rating"`
	URL            string         `json:"url"`
	ShortURL       string         `json:"short_url"`
	Summary        string         `json:"summary"`
	Styles         []string       `json:"styles"`
	FollowStatus   int64          `json:"follow_status"`
	Progress       string         `json:"progress"`
	BothFollow     bool           `json:"both_follow"`
}

type BangumiEpisode struct {
	ID        int64   `json:"id"`
	IndexShow string  `json:"index_show"`
	Cover     string  `json:"cover"`
	Title     string  `json:"title"`
	LongTitle *string `json:"long_title"`
	PubTime   string  `json:"pub_time"`
	Duration  int64   `json:"duration"`
}

type BangumiRating struct {
	Score float64 `json:"score"`
	Count int64   `json:"count"`
}

type RelationStat struct {
	MID       ids.MID `json:"mid"`
	Following uint64  `json:"following"`
	Whisper   uint64  `json:"whisper"`
	Black     uint64  `json:"black"`
	Follower  uint64  `json:"follower"`
}

type Followings struct {
	List      []Following `json:"list"`
	REVersion uint32      `json:"re_version"`
	Total     uint64      `json:"total"`
}

type Following struct {
	MID            ids.MID           `json:"mid"`
	Attribute      uint8             `json:"attribute"`
	ModifiedAt     uint64            `json:"mtime"`
	Tags           []uint64          `json:"tag"`
	Special        uint8             `json:"special"`
	Name           string            `json:"uname"`
	Face           string            `json:"face"`
	Sign           string            `json:"sign"`
	FaceNFT        uint8             `json:"face_nft"`
	OfficialVerify *RelationOfficial `json:"official_verify"`
	VIP            *json.RawMessage  `json:"vip"`
}

type Followers struct {
	List      []Follower `json:"list"`
	Offset    string     `json:"offset"`
	REVersion uint32     `json:"re_version"`
	Total     uint64     `json:"total"`
}

type Follower struct {
	MID            ids.MID           `json:"mid"`
	Attribute      uint8             `json:"attribute"`
	ModifiedAt     *uint64           `json:"mtime"`
	Tags           []uint64          `json:"tag"`
	Special        uint8             `json:"special"`
	ContractInfo   *json.RawMessage  `json:"contract_info"`
	Name           string            `json:"uname"`
	Face           string            `json:"face"`
	Sign           string            `json:"sign"`
	FaceNFT        uint8             `json:"face_nft"`
	OfficialVerify *RelationOfficial `json:"official_verify"`
	VIP            *json.RawMessage  `json:"vip"`
}

type RelationOfficial struct {
	Type int8   `json:"type"`
	Desc string `json:"desc"`
}

type FollowTag struct {
	ID    int64   `json:"tagid"`
	Name  string  `json:"name"`
	Count int64   `json:"count"`
	Tip   *string `json:"tip"`
}

type MedalWall struct {
	List            []MedalWallItem `json:"list"`
	Count           uint32          `json:"count"`
	CloseSpaceMedal uint32          `json:"close_space_medal"`
	OnlyShowWearing uint32          `json:"only_show_wearing"`
	Name            string          `json:"name"`
	Icon            string          `json:"icon"`
	UID             ids.MID         `json:"uid"`
	Level           uint32          `json:"level"`
}

type MedalWallItem struct {
	MedalInfo  MedalInfo       `json:"medal_info"`
	TargetName string          `json:"target_name"`
	TargetIcon string          `json:"target_icon"`
	Link       string          `json:"link"`
	LiveStatus uint32          `json:"live_status"`
	Official   *uint32         `json:"offical"`
	OwnerInfo  *MedalOwnerInfo `json:"uinfo_medal"`
}

func (item MedalWallItem) TargetID() ids.MID { return item.MedalInfo.TargetID }

type MedalInfo struct {
	TargetID      ids.MID `json:"target_id"`
	Level         uint32  `json:"level"`
	Name          string  `json:"medal_name"`
	ColorStart    uint32  `json:"medal_color_start"`
	ColorEnd      uint32  `json:"medal_color_end"`
	ColorBorder   uint32  `json:"medal_color_border"`
	GuardLevel    uint32  `json:"guard_level"`
	WearingStatus uint32  `json:"wearing_status"`
	MedalID       uint64  `json:"medal_id"`
	Intimacy      uint64  `json:"intimacy"`
	NextIntimacy  uint64  `json:"next_intimacy"`
	TodayFeed     uint64  `json:"today_feed"`
	DayLimit      uint64  `json:"day_limit"`
	GuardIcon     *string `json:"guard_icon"`
	HonorIcon     *string `json:"honor_icon"`
}

type MedalOwnerInfo struct {
	Name             string  `json:"name"`
	Level            uint32  `json:"level"`
	ColorStart       uint32  `json:"color_start"`
	ColorEnd         uint32  `json:"color_end"`
	ColorBorder      uint32  `json:"color_border"`
	Color            uint32  `json:"color"`
	ID               uint64  `json:"id"`
	Type             uint32  `json:"typ"`
	IsLight          uint32  `json:"is_light"`
	RUID             ids.MID `json:"ruid"`
	GuardLevel       uint32  `json:"guard_level"`
	Score            uint64  `json:"score"`
	GuardIcon        *string `json:"guard_icon"`
	HonorIcon        *string `json:"honor_icon"`
	V2ColorStart     *string `json:"v2_medal_color_start"`
	V2ColorEnd       *string `json:"v2_medal_color_end"`
	V2ColorBorder    *string `json:"v2_medal_color_border"`
	V2ColorText      *string `json:"v2_medal_color_text"`
	V2ColorLevel     *string `json:"v2_medal_color_level"`
	UserReceiveCount *uint32 `json:"user_receive_count"`
}

type UpStat struct {
	Archive UpStatArchive `json:"archive"`
	Article UpStatArticle `json:"article"`
	Likes   uint64        `json:"likes"`
}

type UpStatArchive struct {
	View uint64 `json:"view"`
}
type UpStatArticle struct {
	View uint64 `json:"view"`
}

type NavStat struct {
	Video       uint64      `json:"video"`
	Bangumi     uint64      `json:"bangumi"`
	Cinema      uint64      `json:"cinema"`
	Channel     NavStatPair `json:"channel"`
	Favourite   NavStatPair `json:"favourite"`
	Tag         uint64      `json:"tag"`
	Article     uint64      `json:"article"`
	Playlist    uint64      `json:"playlist"`
	Album       uint64      `json:"album"`
	Audio       uint64      `json:"audio"`
	PUGV        uint64      `json:"pugv"`
	Opus        uint64      `json:"opus"`
	SeasonCount uint64      `json:"season_num"`
}

type NavStatPair struct {
	Master uint64 `json:"master"`
	Guest  uint64 `json:"guest"`
}

type AlbumCount struct {
	All   uint64 `json:"all_count"`
	Draw  uint64 `json:"draw_count"`
	Photo uint64 `json:"photo_count"`
	Daily uint64 `json:"daily_count"`
}

type NameToUID struct {
	Items []NameToUIDItem `json:"uid_list"`
}

type NameToUIDItem struct {
	Name string  `json:"name"`
	MID  ids.MID `json:"uid"`
}

func (item *NameToUIDItem) UnmarshalJSON(data []byte) error {
	type alias NameToUIDItem
	var value alias
	auxiliary := struct {
		MID flexibleMID `json:"uid"`
		*alias
	}{alias: &value}
	if err := json.Unmarshal(data, &auxiliary); err != nil {
		return err
	}
	value.MID = ids.MID(auxiliary.MID)
	*item = NameToUIDItem(value)
	return nil
}

type UploadedVideos struct {
	List           UploadedVideoList     `json:"list"`
	Page           UploadedVideosPage    `json:"page"`
	EpisodicButton *UploadedVideosButton `json:"episodic_button"`
	IsRisk         bool                  `json:"is_risk"`
}

type UploadedVideoList struct {
	TList  json.RawMessage `json:"tlist"`
	Videos []UploadedVideo `json:"vlist"`
}

type UploadedVideo struct {
	AID         ids.AID  `json:"aid"`
	BVID        ids.BVID `json:"bvid"`
	MID         ids.MID  `json:"mid"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Picture     string   `json:"pic"`
	Length      string   `json:"length"`
	Description string   `json:"description"`
	Created     uint64   `json:"created"`
	Play        uint64   `json:"play"`
	Comment     uint64   `json:"comment"`
	TypeID      uint64   `json:"typeid"`
	VideoReview uint64   `json:"video_review"`
	HideClick   bool     `json:"hide_click"`
}

type UploadedVideosPage struct {
	Count uint64 `json:"count"`
	Page  uint32 `json:"pn"`
	Size  uint32 `json:"ps"`
}

type UploadedVideosButton struct {
	Text string `json:"text"`
	URI  string `json:"uri"`
}

type flexibleMID ids.MID

func (mid *flexibleMID) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		value, err := ids.ParseMID(text)
		if err != nil {
			return err
		}
		*mid = flexibleMID(value)
		return nil
	}
	var number uint64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}
	value, err := ids.NewMID(number)
	if err != nil {
		return err
	}
	*mid = flexibleMID(value)
	return nil
}
