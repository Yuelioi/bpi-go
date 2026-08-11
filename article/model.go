package article

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/ids"
)

type Stats struct {
	Coin     int64 `json:"coin"`
	Dislike  int64 `json:"dislike"`
	Dynamic  int64 `json:"dynamic"`
	Favorite int64 `json:"favorite"`
	Like     int64 `json:"like"`
	Reply    int64 `json:"reply"`
	Share    int64 `json:"share"`
	View     int64 `json:"view"`
}

type Author struct {
	MID            ids.MID         `json:"mid"`
	Name           string          `json:"name"`
	Face           string          `json:"face"`
	Level          int32           `json:"level"`
	Fans           int64           `json:"fans"`
	OfficialVerify OfficialVerify  `json:"official_verify"`
	Nameplate      json.RawMessage `json:"nameplate"`
	Pendant        json.RawMessage `json:"pendant"`
	VIP            AuthorVIP       `json:"vip"`
}

type OfficialVerify struct {
	Type int32  `json:"type"`
	Desc string `json:"desc"`
}

type AuthorVIP struct {
	Type      int32            `json:"type"`
	Status    int32            `json:"status"`
	DueDate   int64            `json:"due_date"`
	PayType   int32            `json:"vip_pay_type"`
	ThemeType int32            `json:"theme_type"`
	Label     *json.RawMessage `json:"label"`
}

type Category struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	ParentID int32  `json:"parent_id"`
}

type Media struct {
	Area     string `json:"area"`
	Cover    string `json:"cover"`
	MediaID  int64  `json:"media_id"`
	Score    int32  `json:"score"`
	SeasonID int64  `json:"season_id"`
	Spoiler  int32  `json:"spoiler"`
	Title    string `json:"title"`
	TypeID   int32  `json:"type_id"`
	TypeName string `json:"type_name"`
}

type Info struct {
	Like            int32          `json:"like"`
	Attention       bool           `json:"attention"`
	Favorite        bool           `json:"favorite"`
	Coin            int32          `json:"coin"`
	Stats           Stats          `json:"stats"`
	Title           string         `json:"title"`
	BannerURL       string         `json:"banner_url"`
	MID             ids.MID        `json:"mid"`
	AuthorName      string         `json:"author_name"`
	IsAuthor        bool           `json:"is_author"`
	ImageURLs       []string       `json:"image_urls"`
	OriginImageURLs []string       `json:"origin_image_urls"`
	Shareable       bool           `json:"shareable"`
	ShowLaterWatch  bool           `json:"show_later_watch"`
	ShowSmallWindow bool           `json:"show_small_window"`
	InList          bool           `json:"in_list"`
	Previous        int64          `json:"pre"`
	Next            int64          `json:"next"`
	ShareChannels   []ShareChannel `json:"share_channels"`
	Type            int32          `json:"type"`
	VideoURL        string         `json:"video_url"`
	Location        string         `json:"location"`
	DisableShare    bool           `json:"disable_share"`
}

type ShareChannel struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Channel string `json:"share_channel"`
}

type Articles struct {
	List      Collection    `json:"list"`
	Articles  []ArticleItem `json:"articles"`
	Author    Author        `json:"author"`
	Last      ArticleItem   `json:"last"`
	Attention bool          `json:"attention"`
}

type Collection struct {
	ID           uint64  `json:"id"`
	MID          ids.MID `json:"mid"`
	Name         string  `json:"name"`
	ImageURL     string  `json:"image_url"`
	UpdateTime   int64   `json:"update_time"`
	CreatedAt    int64   `json:"ctime"`
	PublishedAt  int64   `json:"publish_time"`
	Summary      string  `json:"summary"`
	Words        int64   `json:"words"`
	Read         int64   `json:"read"`
	ArticleCount int32   `json:"articles_count"`
	State        int32   `json:"state"`
	Reason       string  `json:"reason"`
	ApplyTime    string  `json:"apply_time"`
	CheckTime    string  `json:"check_time"`
}

type ArticleItem struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	State       int32      `json:"state"`
	PublishTime int64      `json:"publish_time"`
	Words       int64      `json:"words"`
	ImageURLs   []string   `json:"image_urls"`
	Category    Category   `json:"category"`
	Categories  []Category `json:"categories"`
	Summary     string     `json:"summary"`
	Stats       *Stats     `json:"stats"`
	LikeState   *int32     `json:"like_state"`
}

type View struct {
	ActID             int64            `json:"act_id"`
	ApplyTime         string           `json:"apply_time"`
	Attributes        *int32           `json:"attributes"`
	AuthenMark        *json.RawMessage `json:"authenMark"`
	Author            Author           `json:"author"`
	BannerURL         string           `json:"banner_url"`
	Categories        []Category       `json:"categories"`
	Category          Category         `json:"category"`
	CheckState        int32            `json:"check_state"`
	CheckTime         string           `json:"check_time"`
	Content           string           `json:"content"`
	ContentPictures   *json.RawMessage `json:"content_pic_list"`
	CoverAID          int64            `json:"cover_avid"`
	CreatedAt         int64            `json:"ctime"`
	Dispute           *json.RawMessage `json:"dispute"`
	DynamicID         string           `json:"dyn_id_str"`
	Dynamic           *string          `json:"dynamic"`
	ID                ids.CVID         `json:"id"`
	ImageURLs         []string         `json:"image_urls"`
	IsLike            bool             `json:"is_like"`
	Keywords          string           `json:"keywords"`
	List              *Collection      `json:"list"`
	Media             Media            `json:"media"`
	ModifiedAt        int64            `json:"mtime"`
	Opus              *json.RawMessage `json:"opus"`
	OriginImageURLs   []string         `json:"origin_image_urls"`
	OriginTemplateID  int32            `json:"origin_template_id"`
	Original          int32            `json:"original"`
	PrivatePublish    int32            `json:"private_pub"`
	PublishTime       int64            `json:"publish_time"`
	Reprint           int32            `json:"reprint"`
	State             int32            `json:"state"`
	Stats             Stats            `json:"stats"`
	Summary           string           `json:"summary"`
	Tags              []Tag            `json:"tags"`
	TemplateID        int32            `json:"template_id"`
	Title             string           `json:"title"`
	TopVideoInfo      *json.RawMessage `json:"top_video_info"`
	TotalArticleCount int64            `json:"total_art_num"`
	Type              int32            `json:"type"`
	VersionID         int64            `json:"version_id"`
	Words             int64            `json:"words"`
}

type Tag struct {
	ID   int32  `json:"tid"`
	Name string `json:"name"`
}
