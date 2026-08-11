package article

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/ids"
)

type CardData map[string]CardItem

type CardKind string

const (
	CardVideo   CardKind = "video"
	CardArticle CardKind = "article"
	CardLive    CardKind = "live"
	CardUnknown CardKind = "unknown"
)

// CardItem is one typed article embed. Raw is populated only when Bilibili
// returns an unknown or schema-incompatible card variant.
type CardItem struct {
	Video   *VideoCard
	Article *ArticleCard
	Live    *LiveCard
	Raw     json.RawMessage
}

func (item CardItem) Kind() CardKind {
	switch {
	case item.Video != nil:
		return CardVideo
	case item.Article != nil:
		return CardArticle
	case item.Live != nil:
		return CardLive
	default:
		return CardUnknown
	}
}

func (item *CardItem) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, exists := fields["aid"]; exists {
		var value VideoCard
		if err := json.Unmarshal(data, &value); err == nil {
			*item = CardItem{Video: &value}
			return nil
		}
	}
	if _, exists := fields["room_id"]; exists {
		var value LiveCard
		if err := json.Unmarshal(data, &value); err == nil {
			*item = CardItem{Live: &value}
			return nil
		}
	}
	if _, exists := fields["id"]; exists {
		var value ArticleCard
		if err := json.Unmarshal(data, &value); err == nil {
			*item = CardItem{Article: &value}
			return nil
		}
	}
	*item = CardItem{Raw: append(json.RawMessage(nil), data...)}
	return nil
}

type VideoCard struct {
	AID         ids.AID        `json:"aid"`
	BVID        ids.BVID       `json:"bvid"`
	CID         ids.CID        `json:"cid"`
	Copyright   int32          `json:"copyright"`
	Picture     string         `json:"pic"`
	CreatedAt   int64          `json:"ctime"`
	Description string         `json:"desc"`
	Dimension   VideoDimension `json:"dimension"`
	Duration    int64          `json:"duration"`
	Dynamic     string         `json:"dynamic"`
	Owner       VideoOwner     `json:"owner"`
	PublishedAt int64          `json:"pubdate"`
	Rights      VideoRights    `json:"rights"`
	ShortLink   string         `json:"short_link_v2"`
	Stat        VideoStat      `json:"stat"`
	State       int32          `json:"state"`
	TID         int32          `json:"tid"`
	Title       string         `json:"title"`
	TypeName    string         `json:"tname"`
	Videos      int32          `json:"videos"`
	VTSwitch    bool           `json:"vt_switch"`
}

type VideoDimension struct {
	Height int32 `json:"height"`
	Rotate int32 `json:"rotate"`
	Width  int32 `json:"width"`
}

type VideoOwner struct {
	Face string  `json:"face"`
	MID  ids.MID `json:"mid"`
	Name string  `json:"name"`
}

type VideoRights struct {
	ArcPay        int32 `json:"arc_pay"`
	Autoplay      int32 `json:"autoplay"`
	BP            int32 `json:"bp"`
	Download      int32 `json:"download"`
	Electric      int32 `json:"elec"`
	HD5           int32 `json:"hd5"`
	IsCooperation int32 `json:"is_cooperation"`
	Movie         int32 `json:"movie"`
	NoBackground  int32 `json:"no_background"`
	NoReprint     int32 `json:"no_reprint"`
	Pay           int32 `json:"pay"`
	PayFreeWatch  int32 `json:"pay_free_watch"`
	UGCPay        int32 `json:"ugc_pay"`
	UGCPayPreview int32 `json:"ugc_pay_preview"`
}

type VideoStat struct {
	AID      ids.AID `json:"aid"`
	Coin     int64   `json:"coin"`
	Danmaku  int64   `json:"danmaku"`
	Dislike  int64   `json:"dislike"`
	Favorite int64   `json:"favorite"`
	HisRank  int32   `json:"his_rank"`
	Like     int64   `json:"like"`
	NowRank  int32   `json:"now_rank"`
	Reply    int64   `json:"reply"`
	Share    int64   `json:"share"`
	View     int64   `json:"view"`
	VT       int32   `json:"vt"`
	VV       int32   `json:"vv"`
}

type ArticleCard struct {
	ActID            int64            `json:"act_id"`
	ApplyTime        string           `json:"apply_time"`
	Attributes       *int32           `json:"attributes"`
	AuthenMark       *json.RawMessage `json:"authenMark"`
	Author           Author           `json:"author"`
	BannerURL        string           `json:"banner_url"`
	Categories       []Category       `json:"categories"`
	Category         Category         `json:"category"`
	CheckState       int32            `json:"check_state"`
	CheckTime        string           `json:"check_time"`
	ContentPictures  *json.RawMessage `json:"content_pic_list"`
	CoverAID         int64            `json:"cover_avid"`
	CreatedAt        int64            `json:"ctime"`
	Dispute          *json.RawMessage `json:"dispute"`
	Dynamic          *string          `json:"dynamic"`
	ID               ids.CVID         `json:"id"`
	ImageURLs        []string         `json:"image_urls"`
	IsLike           bool             `json:"is_like"`
	List             *Collection      `json:"list"`
	Media            Media            `json:"media"`
	ModifiedAt       int64            `json:"mtime"`
	OriginImageURLs  []string         `json:"origin_image_urls"`
	OriginTemplateID int32            `json:"origin_template_id"`
	Original         int32            `json:"original"`
	PrivatePublish   int32            `json:"private_pub"`
	PublishTime      int64            `json:"publish_time"`
	Reprint          int32            `json:"reprint"`
	State            int32            `json:"state"`
	Stats            Stats            `json:"stats"`
	Summary          string           `json:"summary"`
	TemplateID       int32            `json:"template_id"`
	Title            string           `json:"title"`
	TopVideoInfo     *json.RawMessage `json:"top_video_info"`
	Type             int32            `json:"type"`
	Words            int64            `json:"words"`
}

type LiveCard struct {
	AreaName       string     `json:"area_v2_name"`
	Cover          string     `json:"cover"`
	Face           string     `json:"face"`
	LiveStatus     int32      `json:"live_status"`
	Online         int64      `json:"online"`
	Pendant        string     `json:"pendent_ru"`
	PendantColor   string     `json:"pendent_ru_color"`
	PendantPicture string     `json:"pendent_ru_pic"`
	Role           int32      `json:"role"`
	RoomID         ids.RoomID `json:"room_id"`
	Title          string     `json:"title"`
	UID            ids.MID    `json:"uid"`
	Name           string     `json:"uname"`
}
