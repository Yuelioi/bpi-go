package search

import "encoding/json"

type Data[T any] struct {
	SEID       string        `json:"seid"`
	Page       int64         `json:"page"`
	PageSize   int64         `json:"pagesize"`
	NumResults int64         `json:"numResults"`
	NumPages   int64         `json:"numPages"`
	Result     *T            `json:"result"`
	PageInfo   *LivePageInfo `json:"pageinfo"`
}

type LivePageInfo struct {
	Users LivePageCounts `json:"live_user"`
	Rooms LivePageCounts `json:"live_room"`
}

type LivePageCounts struct {
	Total      int64 `json:"total"`
	NumResults int64 `json:"numResults"`
	Pages      int64 `json:"pages"`
	NumPages   int64 `json:"numPages"`
}

type Article struct {
	CategoryID   int64    `json:"category_id"`
	CategoryName string   `json:"category_name"`
	CommentURL   string   `json:"comment_url"`
	Description  string   `json:"desc"`
	ID           int64    `json:"id"`
	ImageURLs    []string `json:"image_urls"`
	IsComment    int64    `json:"is_comment"`
	IsFold       bool     `json:"is_fold"`
	IsRankOne    bool     `json:"is_rk1"`
	Like         int64    `json:"like"`
	MID          int64    `json:"mid"`
	PublishTime  int64    `json:"pub_time"`
	RankIndex    int64    `json:"rank_index"`
	RankOffset   int64    `json:"rank_offset"`
	Reply        int64    `json:"reply"`
	SpreadID     int64    `json:"spread_id"`
	Subtype      int64    `json:"sub_type"`
	TemplateID   int64    `json:"template_id"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	Version      string   `json:"version"`
	View         int64    `json:"view"`
}

type Media struct {
	Type           string     `json:"type"`
	MediaID        int64      `json:"media_id"`
	Title          string     `json:"title"`
	OriginalTitle  string     `json:"org_title"`
	MediaType      int64      `json:"media_type"`
	CV             string     `json:"cv"`
	Staff          string     `json:"staff"`
	SeasonID       int64      `json:"season_id"`
	IsAvid         bool       `json:"is_avid"`
	HitEpisodeIDs  string     `json:"hit_epids"`
	SeasonType     int64      `json:"season_type"`
	SeasonTypeName string     `json:"season_type_name"`
	SelectionStyle string     `json:"selection_style"`
	EpisodeSize    int64      `json:"ep_size"`
	URL            string     `json:"url"`
	ButtonText     string     `json:"button_text"`
	IsFollow       int64      `json:"is_follow"`
	IsSelection    int64      `json:"is_selection"`
	Episodes       []Episode  `json:"eps"`
	Badges         []Badge    `json:"badges"`
	Cover          string     `json:"cover"`
	Areas          string     `json:"areas"`
	Styles         string     `json:"styles"`
	GotoURL        string     `json:"goto_url"`
	Description    string     `json:"desc"`
	PublishTime    int64      `json:"pubtime"`
	MediaMode      int64      `json:"media_mode"`
	PublishText    string     `json:"fix_pubtime_str"`
	Score          MediaScore `json:"media_score"`
	DisplayInfo    []Badge    `json:"display_info"`
	PGCSeasonID    int64      `json:"pgc_season_id"`
	Corner         int64      `json:"corner"`
	IndexShow      string     `json:"index_show"`
}

type Bangumi = Media
type Movie = Media

type Episode struct {
	ID          int64   `json:"id"`
	Cover       string  `json:"cover"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	ReleaseDate string  `json:"release_date"`
	Badges      []Badge `json:"badges"`
	IndexTitle  string  `json:"index_title"`
	LongTitle   string  `json:"long_title"`
}

type Badge struct {
	Text            string `json:"text"`
	TextColor       string `json:"text_color"`
	TextColorNight  string `json:"text_color_night"`
	Background      string `json:"bg_color"`
	BackgroundNight string `json:"bg_color_night"`
	Border          string `json:"border_color"`
	BorderNight     string `json:"border_color_night"`
	BackgroundStyle int64  `json:"bg_style"`
}

type MediaScore struct {
	Score     float32 `json:"score"`
	UserCount int64   `json:"user_count"`
}

type Video struct {
	Type        string `json:"type"`
	ID          uint64 `json:"id"`
	Author      string `json:"author"`
	MID         uint64 `json:"mid"`
	TypeID      string `json:"typeid"`
	TypeName    string `json:"typename"`
	ArchiveURL  string `json:"arcurl"`
	AID         uint64 `json:"aid"`
	BVID        string `json:"bvid"`
	Title       string `json:"title"`
	Picture     string `json:"pic"`
	Play        uint64 `json:"play"`
	Danmaku     uint64 `json:"danmaku"`
	Favorites   uint64 `json:"favorites"`
	Like        uint64 `json:"like"`
	Tag         string `json:"tag"`
	Review      uint64 `json:"review"`
	PublishTime uint64 `json:"pubdate"`
	Duration    string `json:"duration"`
}

type LiveData struct {
	Rooms []LiveRoom `json:"live_room"`
	Users []LiveUser `json:"live_user"`
}

type LiveUser struct {
	Area       int64    `json:"area"`
	AreaV2ID   int64    `json:"area_v2_id"`
	Attentions int64    `json:"attentions"`
	Category   string   `json:"cate_name"`
	HitColumns []string `json:"hit_columns"`
	IsLive     bool     `json:"is_live"`
	LiveStatus int64    `json:"live_status"`
	LiveTime   string   `json:"live_time"`
	RankIndex  int64    `json:"rank_index"`
	RankOffset int64    `json:"rank_offset"`
	RoomID     int64    `json:"roomid"`
	Tags       string   `json:"tags"`
	Type       string   `json:"type"`
	Face       string   `json:"uface"`
	UID        int64    `json:"uid"`
	Name       string   `json:"uname"`
	ID         *int64   `json:"id"`
}

type LiveRoom struct {
	Area       int64        `json:"area"`
	Attentions int64        `json:"attentions"`
	Category   string       `json:"cate_name"`
	Cover      string       `json:"cover"`
	IsInline   int64        `json:"is_live_room_inline"`
	LiveStatus int64        `json:"live_status"`
	LiveTime   string       `json:"live_time"`
	Online     int64        `json:"online"`
	RankIndex  int64        `json:"rank_index"`
	RankOffset int64        `json:"rank_offset"`
	RoomID     int64        `json:"roomid"`
	ShortID    int64        `json:"short_id"`
	Tags       string       `json:"tags"`
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	Face       string       `json:"uface"`
	UID        int64        `json:"uid"`
	Name       string       `json:"uname"`
	UserCover  string       `json:"user_cover"`
	Watched    *WatchedShow `json:"watched_show"`
}

type WatchedShow struct {
	Switch       bool   `json:"switch"`
	Number       int64  `json:"num"`
	TextSmall    string `json:"text_small"`
	TextLarge    string `json:"text_large"`
	Icon         string `json:"icon"`
	IconLocation string `json:"icon_location"`
	IconWeb      string `json:"icon_web"`
}

type OfficialVerify struct {
	Type        int64  `json:"type"`
	Description string `json:"desc"`
}

type UserVideo struct {
	AID           int64  `json:"aid"`
	BVID          string `json:"bvid"`
	Title         string `json:"title"`
	PublishTime   int64  `json:"pubdate"`
	ArchiveURL    string `json:"arcurl"`
	Picture       string `json:"pic"`
	Play          string `json:"play"`
	Danmaku       int64  `json:"dm"`
	Coin          int64  `json:"coin"`
	Favorite      int64  `json:"fav"`
	Description   string `json:"desc"`
	Duration      string `json:"duration"`
	IsPay         int64  `json:"is_pay"`
	IsUnionVideo  int64  `json:"is_union_video"`
	IsChargeVideo int64  `json:"is_charge_video"`
	VT            int64  `json:"vt"`
	EnableVT      int64  `json:"enable_vt"`
	VTDisplay     string `json:"vt_display"`
}

type User struct {
	Type           string         `json:"type"`
	MID            int64          `json:"mid"`
	Name           string         `json:"uname"`
	Signature      string         `json:"usign"`
	Fans           int64          `json:"fans"`
	Videos         int64          `json:"videos"`
	Picture        string         `json:"upic"`
	FaceNFT        int64          `json:"face_nft"`
	FaceNFTType    int64          `json:"face_nft_type"`
	VerifyInfo     string         `json:"verify_info"`
	Level          int64          `json:"level"`
	Gender         int64          `json:"gender"`
	IsUploader     int64          `json:"is_upuser"`
	IsLive         int64          `json:"is_live"`
	RoomID         int64          `json:"room_id"`
	Results        []UserVideo    `json:"res"`
	OfficialVerify OfficialVerify `json:"official_verify"`
	IsSeniorMember int64          `json:"is_senior_member"`
}

type Default struct {
	SEID      string  `json:"seid"`
	ID        uint64  `json:"id"`
	Type      uint32  `json:"type"`
	ShowName  string  `json:"show_name"`
	Name      *string `json:"name"`
	GotoType  uint32  `json:"goto_type"`
	GotoValue string  `json:"goto_value"`
	URL       string  `json:"url"`
}

type Suggest struct {
	Tags []Suggestion `json:"tag"`
}

type Suggestion struct {
	Value *string `json:"value"`
	Name  *string `json:"name"`
	Type  *string `json:"type"`
}

type HotWords struct {
	Code  uint32    `json:"code"`
	Items []HotWord `json:"list"`
}

type HotWord struct {
	ID        uint64            `json:"hot_id"`
	Keyword   string            `json:"keyword"`
	ShowName  string            `json:"show_name"`
	HeatScore uint64            `json:"heat_score"`
	WordType  uint32            `json:"word_type"`
	LiveID    []json.RawMessage `json:"live_id"`
}
