package creativecenter

import "encoding/json"

type Season struct {
	ID           uint64  `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"desc"`
	Cover        string  `json:"cover"`
	Ended        uint32  `json:"isEnd"`
	MID          uint64  `json:"mid"`
	Activity     uint32  `json:"isAct"`
	Paid         uint32  `json:"is_pay"`
	State        int32   `json:"state"`
	PartState    uint32  `json:"partState"`
	SignState    uint32  `json:"signState"`
	RejectReason *string `json:"rejectReason"`
	CreatedAt    uint64  `json:"ctime"`
	ModifiedAt   uint64  `json:"mtime"`
	NoSection    uint32  `json:"no_section"`
	Forbidden    uint32  `json:"forbid"`
	ProtocolID   *string `json:"protocol_id"`
	EpisodeCount uint32  `json:"ep_num"`
	Price        uint32  `json:"season_price"`
	Opened       uint32  `json:"is_opened"`
	ChargingPay  uint32  `json:"has_charging_pay"`
	PUGVPay      *uint32 `json:"has_pugv_pay"`
	SeasonUpFrom *uint32 `json:"SeasonUpfrom"`
}

type Section struct {
	ID           uint64          `json:"id"`
	Type         uint32          `json:"type"`
	SeasonID     uint64          `json:"seasonId"`
	Title        string          `json:"title"`
	Order        uint32          `json:"order"`
	State        int32           `json:"state"`
	PartState    int32           `json:"partState"`
	CreatedAt    int64           `json:"ctime"`
	ModifiedAt   int64           `json:"mtime"`
	EpisodeCount int64           `json:"epCount"`
	Cover        string          `json:"cover"`
	Episodes     json.RawMessage `json:"Episodes"`
}

type Episode struct {
	ID           uint64  `json:"id"`
	Title        string  `json:"title"`
	AID          uint64  `json:"aid"`
	BVID         string  `json:"bvid"`
	CID          uint64  `json:"cid"`
	SeasonID     uint64  `json:"seasonId"`
	SectionID    uint64  `json:"sectionId"`
	Order        uint32  `json:"order"`
	VideoTitle   *string `json:"videoTitle"`
	ArchiveTitle *string `json:"archiveTitle"`
	ArchiveState int32   `json:"archiveState"`
	RejectReason *string `json:"rejectReason"`
	State        int32   `json:"state"`
	Cover        string  `json:"cover"`
	Free         uint32  `json:"is_free"`
	AIDOwner     bool    `json:"aid_owner"`
	ChargingPay  uint32  `json:"charging_pay"`
}

type SeasonStat struct {
	View         uint64 `json:"view"`
	Danmaku      uint64 `json:"danmaku"`
	Reply        uint64 `json:"reply"`
	Favorite     uint64 `json:"fav"`
	Coin         uint64 `json:"coin"`
	Share        uint64 `json:"share"`
	Like         uint64 `json:"like"`
	Subscription uint64 `json:"subscription"`
	NowRank      uint32 `json:"nowRank"`
	HistoryRank  uint32 `json:"hisRank"`
}

type SeasonItem struct {
	Season       Season          `json:"season"`
	CheckIn      json.RawMessage `json:"checkin"`
	Stat         *SeasonStat     `json:"seasonStat"`
	Sections     *Sections       `json:"sections"`
	PartEpisodes []Episode       `json:"part_episodes"`
}

type Sections struct {
	Sections []Section `json:"sections"`
	Total    int64     `json:"total"`
}

type SeasonList struct {
	Seasons  []SeasonItem    `json:"seasons"`
	Tip      json.RawMessage `json:"tip"`
	Total    uint32          `json:"total"`
	PlayType uint32          `json:"play_type"`
}

type SeasonInfo struct {
	Season       Season          `json:"season"`
	Course       json.RawMessage `json:"course"`
	CheckIn      json.RawMessage `json:"checkin"`
	Stat         json.RawMessage `json:"seasonStat"`
	Sections     Sections        `json:"sections"`
	PartEpisodes json.RawMessage `json:"part_episodes"`
}

type SectionEpisodes struct {
	Section  Section   `json:"section"`
	Episodes []Episode `json:"episodes"`
}

type ArchiveStat struct {
	AID         int64 `json:"aid"`
	View        int64 `json:"view"`
	Danmaku     int64 `json:"danmaku"`
	Reply       int64 `json:"reply"`
	Favorite    int64 `json:"favorite"`
	Coin        int64 `json:"coin"`
	Share       int64 `json:"share"`
	NowRank     int64 `json:"now_rank"`
	HistoryRank int64 `json:"his_rank"`
	Like        int64 `json:"like"`
	Dislike     int64 `json:"dislike"`
	VT          int64 `json:"vt"`
	VV          int64 `json:"vv"`
}

type Archive struct {
	AID         int64  `json:"aid"`
	BVID        string `json:"bvid"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Duration    int64  `json:"duration"`
	Description string `json:"desc"`
}

type ArchiveAudit struct {
	Archive  *Archive        `json:"Archive"`
	Videos   json.RawMessage `json:"Videos"`
	Stat     ArchiveStat     `json:"stat"`
	State    int64           `json:"state_panel"`
	TypeName *string         `json:"typename"`
}

type Page struct {
	Page     int64 `json:"pn"`
	PageSize int64 `json:"ps"`
	Count    int64 `json:"count"`
}

type ArchivesList struct {
	Archives []ArchiveAudit `json:"arc_audits"`
	Page     Page           `json:"page"`
	PlayType int64          `json:"play_type"`
}

type VideoPart struct {
	CID      int64  `json:"cid"`
	Index    int64  `json:"index"`
	Duration int64  `json:"duration"`
	Title    string `json:"title"`
}

type ArchiveVideos struct {
	Archive Archive     `json:"archive"`
	Videos  []VideoPart `json:"videos"`
}

type UpStat struct {
	IncrementCoin     int64 `json:"inc_coin"`
	IncrementElectric int64 `json:"inc_elec"`
	IncrementFavorite int64 `json:"inc_fav"`
	IncrementLike     int64 `json:"inc_like"`
	IncrementShare    int64 `json:"inc_share"`
	IncrementClick    int64 `json:"incr_click"`
	IncrementDanmaku  int64 `json:"incr_dm"`
	IncrementFans     int64 `json:"incr_fans"`
	IncrementReply    int64 `json:"incr_reply"`
	TotalClick        int64 `json:"total_click"`
	TotalCoin         int64 `json:"total_coin"`
	TotalDanmaku      int64 `json:"total_dm"`
	TotalElectric     int64 `json:"total_elec"`
	TotalFans         int64 `json:"total_fans"`
	TotalFavorite     int64 `json:"total_fav"`
	TotalLike         int64 `json:"total_like"`
	TotalReply        int64 `json:"total_reply"`
	TotalShare        int64 `json:"total_share"`
}

// ArchiveCompareStat is intentionally retained as raw JSON: the upstream
// diagnose schema contains dozens of experimental, frequently changing rates.
type ArchiveCompareItem struct {
	AID         int64           `json:"aid"`
	BVID        string          `json:"bvid"`
	Cover       string          `json:"cover"`
	Title       string          `json:"title"`
	PublishedAt int64           `json:"pubtime"`
	Duration    int64           `json:"duration"`
	Stat        json.RawMessage `json:"stat"`
	HourStat    json.RawMessage `json:"hour_stat"`
}
type ArchiveCompare struct {
	List []ArchiveCompareItem `json:"list"`
}

type ArticleStat struct {
	View              int64 `json:"view"`
	Reply             int64 `json:"reply"`
	Like              int64 `json:"like"`
	Coin              int64 `json:"coin"`
	Favorite          int64 `json:"fav"`
	Share             int64 `json:"share"`
	IncrementView     int64 `json:"incr_view"`
	IncrementReply    int64 `json:"incr_reply"`
	IncrementLike     int64 `json:"incr_like"`
	IncrementCoin     int64 `json:"incr_coin"`
	IncrementFavorite int64 `json:"incr_fav"`
	IncrementShare    int64 `json:"incr_share"`
}

type Trend struct {
	Date      int64 `json:"date_key"`
	Increment int64 `json:"total_inc"`
}

type PlaySource struct {
	PageSource struct {
		Dynamic, Other, Search, Space, Tenma int64
		RelatedVideo                         int64 `json:"related_video"`
	} `json:"page_source"`
	Proportion struct{ Android, H5, IOS, Out, PC int64 } `json:"play_proportion"`
}

type ViewerBaseDetail struct {
	Male     int64 `json:"male"`
	Female   int64 `json:"female"`
	AgeOne   int64 `json:"age_one"`
	AgeTwo   int64 `json:"age_two"`
	AgeThree int64 `json:"age_three"`
	AgeFour  int64 `json:"age_four"`
	PC       int64 `json:"plat_pc"`
	H5       int64 `json:"plat_h5"`
	Out      int64 `json:"plat_out"`
	IOS      int64 `json:"plat_ios"`
	Android  int64 `json:"plat_android"`
	OtherApp int64 `json:"plat_other_app"`
}

type ViewerData struct {
	Period     map[string]string `json:"period"`
	ViewerArea struct {
		Fan    map[string]int64 `json:"fan"`
		NotFan map[string]int64 `json:"not_fan"`
	} `json:"viewer_area"`
	ViewerBase struct {
		Fan    ViewerBaseDetail `json:"fan"`
		NotFan ViewerBaseDetail `json:"not_fan"`
	} `json:"viewer_base"`
}

type ElectromagneticInfo struct {
	MID       uint64 `json:"mid"`
	Level     uint32 `json:"level"`
	Score     uint32 `json:"score"`
	Credit    uint32 `json:"credit"`
	State     int32  `json:"state"`
	UpdatedAt uint64 `json:"update_date"`
}
