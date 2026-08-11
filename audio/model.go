package audio

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/ids"
)

type Info struct {
	ID               ids.AudioID      `json:"id"`
	UID              ids.MID          `json:"uid"`
	UserName         string           `json:"uname"`
	Author           string           `json:"author"`
	Title            string           `json:"title"`
	Cover            string           `json:"cover"`
	Intro            string           `json:"intro"`
	LyricURL         string           `json:"lyric"`
	CRType           int32            `json:"crtype"`
	Duration         int64            `json:"duration"`
	PassTime         int64            `json:"passtime"`
	CurrentTime      int64            `json:"curtime"`
	AID              int64            `json:"aid"`
	BVID             string           `json:"bvid"`
	CID              int64            `json:"cid"`
	MSID             int64            `json:"msid"`
	Attribute        int64            `json:"attr"`
	Limit            int64            `json:"limit"`
	ActivityID       int64            `json:"activityId"`
	LimitDescription string           `json:"limitdesc"`
	CreatedAt        *json.RawMessage `json:"ctime"`
	Statistic        Statistic        `json:"statistic"`
	VIP              VIPInfo          `json:"vipInfo"`
	CollectIDs       []int64          `json:"collectIds"`
	CoinCount        int64            `json:"coin_num"`
}

type Statistic struct {
	SID     ids.AudioID `json:"sid"`
	Play    int64       `json:"play"`
	Collect int64       `json:"collect"`
	Comment int64       `json:"comment"`
	Share   int64       `json:"share"`
}

type VIPInfo struct {
	Type    int32 `json:"type"`
	Status  int32 `json:"status"`
	DueDate int64 `json:"due_date"`
	PayType int32 `json:"vip_pay_type"`
}

type Tag struct {
	Type    string `json:"type"`
	Subtype int32  `json:"subtype"`
	Key     int32  `json:"key"`
	Info    string `json:"info"`
}

type MemberGroup struct {
	List []Member `json:"list"`
	Type int32    `json:"type"`
}

type Member struct {
	MID      int64  `json:"mid"`
	Name     string `json:"name"`
	MemberID int64  `json:"member_id"`
}

type StatusNumber struct {
	SID     ids.AudioID `json:"sid"`
	Play    int64       `json:"play"`
	Collect int64       `json:"collect"`
	Comment int64       `json:"comment"`
	Share   int64       `json:"share"`
}

type Page[T any] struct {
	CurrentPage int32 `json:"curPage"`
	PageCount   int32 `json:"pageCount"`
	TotalSize   int32 `json:"totalSize"`
	PageSize    int32 `json:"pageSize"`
	Data        []T   `json:"data"`
}

type Collection struct {
	ID          int64               `json:"id"`
	UID         int64               `json:"uid"`
	UserName    string              `json:"uname"`
	Title       string              `json:"title"`
	Type        int32               `json:"type"`
	Published   int32               `json:"published"`
	Cover       string              `json:"cover"`
	CreatedAt   int64               `json:"ctime"`
	SongCount   int32               `json:"song"`
	Description string              `json:"desc"`
	SIDs        []int64             `json:"sids"`
	MenuID      int64               `json:"menuId"`
	Statistic   CollectionStatistic `json:"statistic"`
}

type CollectionStatistic struct {
	SID     int64  `json:"sid"`
	Play    int64  `json:"play"`
	Collect int64  `json:"collect"`
	Comment *int64 `json:"comment"`
	Share   int64  `json:"share"`
}

type HotMenu struct {
	MenuID       int64         `json:"menuId"`
	UID          int64         `json:"uid"`
	UserName     string        `json:"uname"`
	Title        string        `json:"title"`
	Cover        string        `json:"cover"`
	Intro        string        `json:"intro"`
	Type         int32         `json:"type"`
	Off          int32         `json:"off"`
	CreatedAt    int64         `json:"ctime"`
	CurrentTime  int64         `json:"curtime"`
	Statistic    MenuStatistic `json:"statistic"`
	SongCount    int32         `json:"snum"`
	Attribute    int32         `json:"attr"`
	IsDefault    int32         `json:"isDefault"`
	CollectionID int64         `json:"collectionId"`
}

type MenuStatistic struct {
	SID     int64 `json:"sid"`
	Play    int64 `json:"play"`
	Collect int64 `json:"collect"`
	Comment int64 `json:"comment"`
	Share   int64 `json:"share"`
}

type RankMenu struct {
	MenuID       int64         `json:"menuId"`
	UID          int64         `json:"uid"`
	UserName     string        `json:"uname"`
	Title        string        `json:"title"`
	Cover        string        `json:"cover"`
	Intro        string        `json:"intro"`
	Type         int32         `json:"type"`
	Off          int32         `json:"off"`
	CreatedAt    int64         `json:"ctime"`
	CurrentTime  int64         `json:"curtime"`
	Statistic    MenuStatistic `json:"statistic"`
	SongCount    int32         `json:"snum"`
	Attribute    int32         `json:"attr"`
	IsDefault    int32         `json:"isDefault"`
	CollectionID int64         `json:"collectionId"`
	Audios       []RankItem    `json:"audios"`
}

type RankItem struct {
	ID       ids.AudioID `json:"id"`
	Title    string      `json:"title"`
	Duration int64       `json:"duration"`
}

type StreamURLWeb struct {
	SID       ids.AudioID      `json:"sid"`
	Type      uint32           `json:"type"`
	Info      string           `json:"info"`
	Timeout   uint64           `json:"timeout"`
	Size      uint64           `json:"size"`
	CDNs      []string         `json:"cdns"`
	Qualities *json.RawMessage `json:"qualities"`
	Title     string           `json:"title"`
	Cover     string           `json:"cover"`
}

type StreamURL struct {
	SID       ids.AudioID   `json:"sid"`
	Type      uint32        `json:"type"`
	Info      string        `json:"info"`
	Timeout   uint64        `json:"timeout"`
	Size      uint64        `json:"size"`
	CDNs      []string      `json:"cdns"`
	Qualities []QualityInfo `json:"qualities"`
	Title     string        `json:"title"`
	Cover     string        `json:"cover"`
}

type QualityInfo struct {
	Type               uint32 `json:"type"`
	Description        string `json:"desc"`
	Size               uint64 `json:"size"`
	BPS                string `json:"bps"`
	Tag                string `json:"tag"`
	Require            uint32 `json:"require"`
	RequireDescription string `json:"requiredesc"`
}

type RankPeriods struct {
	List map[string][]RankPeriod `json:"list"`
}

type RankPeriod struct {
	ID          uint64 `json:"ID"`
	Period      uint64 `json:"priod"`
	PublishTime uint64 `json:"publish_time"`
}

type RankDetail struct {
	ListenFID   uint64 `json:"listen_fid"`
	AllFID      uint64 `json:"all_fid"`
	FavoriteMID uint64 `json:"fav_mid"`
	CoverURL    string `json:"cover_url"`
	IsSubscribe bool   `json:"is_subscribe"`
	ListenCount uint64 `json:"listen_count"`
}

type RankMusicList struct {
	List []RankMusic `json:"list"`
}

type RankMusic struct {
	MusicID          string   `json:"music_id"`
	MusicTitle       string   `json:"music_title"`
	Singer           string   `json:"singer"`
	Album            string   `json:"album"`
	MVAID            uint64   `json:"mv_aid"`
	MVBVID           string   `json:"mv_bvid"`
	MVCover          string   `json:"mv_cover"`
	Heat             uint64   `json:"heat"`
	Rank             uint64   `json:"rank"`
	CanListen        bool     `json:"can_listen"`
	Recommendation   string   `json:"recommendation"`
	CreationAID      uint64   `json:"creation_aid"`
	CreationBVID     string   `json:"creation_bvid"`
	CreationCover    string   `json:"creation_cover"`
	CreationTitle    string   `json:"creation_title"`
	CreationUP       uint64   `json:"creation_up"`
	CreationNickname string   `json:"creation_nickname"`
	CreationDuration uint64   `json:"creation_duration"`
	CreationPlay     uint64   `json:"creation_play"`
	CreationReason   string   `json:"creation_reason"`
	Achievements     []string `json:"achievements"`
	MaterialID       uint64   `json:"material_id"`
	MaterialUseCount uint64   `json:"material_use_num"`
	MaterialDuration uint64   `json:"material_duration"`
	MaterialShow     uint64   `json:"material_show"`
	SongType         uint64   `json:"song_type"`
}
