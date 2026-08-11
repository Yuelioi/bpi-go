package dynamic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Feed struct {
	HasMore        bool       `json:"has_more"`
	Items          []FeedItem `json:"items"`
	Offset         string     `json:"offset"`
	UpdateBaseline string     `json:"update_baseline"`
	UpdateNum      uint64     `json:"-"`
}

func (feed *Feed) UnmarshalJSON(data []byte) error {
	var value struct {
		HasMore        bool            `json:"has_more"`
		Items          []FeedItem      `json:"items"`
		Offset         string          `json:"offset"`
		UpdateBaseline string          `json:"update_baseline"`
		UpdateNum      json.RawMessage `json:"update_num"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	updateNum, err := decodeUint64(value.UpdateNum)
	if err != nil {
		return fmt.Errorf("dynamic update_num: %w", err)
	}
	*feed = Feed{HasMore: value.HasMore, Items: value.Items, Offset: value.Offset, UpdateBaseline: value.UpdateBaseline, UpdateNum: updateNum}
	return nil
}

type FeedItem struct {
	Basic   Basic           `json:"basic"`
	ID      string          `json:"id_str"`
	Modules json.RawMessage `json:"modules"`
	Type    string          `json:"type"`
	Visible bool            `json:"visible"`
}

type Basic struct {
	CommentID   string          `json:"comment_id_str"`
	CommentType int64           `json:"comment_type"`
	LikeIcon    json.RawMessage `json:"like_icon"`
	RID         string          `json:"rid_str"`
	OnlyFans    *bool           `json:"is_only_fans"`
	JumpURL     *string         `json:"jump_url"`
	Editable    *bool           `json:"editable"`
}

type Update struct {
	UpdateNum uint64 `json:"update_num"`
}

type BannerFeed struct {
	Banners []Banner `json:"banners"`
}

type Banner struct {
	ID        uint64 `json:"banner_id"`
	EndTime   uint64 `json:"end_time"`
	ImageURL  string `json:"img_url"`
	Link      string `json:"link"`
	Platform  uint64 `json:"platform"`
	Position  string `json:"position"`
	StartTime uint64 `json:"start_time"`
	Title     string `json:"title"`
	Weight    uint64 `json:"weight"`
}

type Detail struct {
	Item DetailItem `json:"item"`
}

type DetailItem struct {
	ID       string          `json:"id_str"`
	Basic    Basic           `json:"basic"`
	Modules  json.RawMessage `json:"modules"`
	Original *DetailItem     `json:"orig"`
	Type     string          `json:"type"`
	Visible  bool            `json:"visible"`
}

type Reaction struct {
	Action      string `json:"action"`
	Attend      uint8  `json:"attend"`
	Description string `json:"desc"`
	Face        string `json:"face"`
	MID         string `json:"mid"`
	Name        string `json:"name"`
}

type Reactions struct {
	HasMore bool       `json:"has_more"`
	Items   []Reaction `json:"items"`
	Offset  string     `json:"offset"`
	Total   uint64     `json:"total"`
}

type Forward struct {
	Description ForwardDescription `json:"desc"`
	ID          string             `json:"id_str"`
	PublishTime string             `json:"pub_time"`
	User        ForwardUser        `json:"user"`
}

type ForwardDescription struct {
	Nodes []RichTextNode `json:"rich_text_nodes"`
	Text  string         `json:"text"`
}

type RichTextNode struct {
	OriginalText string `json:"orig_text"`
	Text         string `json:"text"`
	Type         string `json:"type"`
}

type ForwardUser struct {
	Face    string `json:"face"`
	FaceNFT bool   `json:"face_nft"`
	MID     int64  `json:"mid"`
	Name    string `json:"name"`
}

type Forwards struct {
	HasMore bool      `json:"has_more"`
	Items   []Forward `json:"items"`
	Offset  string    `json:"offset"`
	Total   uint64    `json:"total"`
}

type ForwardInfo struct {
	Item Forward `json:"item"`
}

type Picture struct {
	Height uint64  `json:"height"`
	Size   float64 `json:"size"`
	Source string  `json:"src"`
	Width  uint64  `json:"width"`
}

type LotteryWinner struct {
	UID          uint64   `json:"uid"`
	Name         string   `json:"name"`
	Face         string   `json:"face"`
	HongbaoMoney *float64 `json:"hongbao_money"`
}

type LotteryResult struct {
	First  []LotteryWinner `json:"first_prize_result"`
	Second []LotteryWinner `json:"second_prize_result"`
	Third  []LotteryWinner `json:"third_prize_result"`
}

type LotteryPrizeType struct {
	Type  uint8                 `json:"type"`
	Value LotteryPrizeTypeValue `json:"value"`
}

type LotteryPrizeTypeValue struct {
	Count uint64 `json:"count"`
	Type  uint8  `json:"stype"`
}

type LotteryNotice struct {
	LotteryID          uint64            `json:"lottery_id"`
	SenderUID          uint64            `json:"sender_uid"`
	BusinessType       uint8             `json:"business_type"`
	BusinessID         uint64            `json:"business_id"`
	Status             uint8             `json:"status"`
	LotteryTime        uint64            `json:"lottery_time"`
	Participants       uint64            `json:"participants"`
	FirstPrize         uint32            `json:"first_prize"`
	FirstPrizeComment  string            `json:"first_prize_cmt"`
	FirstPrizePicture  string            `json:"first_prize_pic"`
	SecondPrize        uint32            `json:"second_prize"`
	SecondPrizeComment *string           `json:"second_prize_cmt"`
	SecondPrizePicture string            `json:"second_prize_pic"`
	ThirdPrize         uint32            `json:"third_prize"`
	ThirdPrizeComment  *string           `json:"third_prize_cmt"`
	ThirdPrizePicture  string            `json:"third_prize_pic"`
	Result             *LotteryResult    `json:"lottery_result"`
	Followed           bool              `json:"followed"`
	HasChargeRight     bool              `json:"has_charge_right"`
	AtCount            uint32            `json:"lottery_at_num"`
	DetailURL          string            `json:"lottery_detail_url"`
	FeedLimit          uint32            `json:"lottery_feed_limit"`
	NeedPost           uint8             `json:"need_post"`
	Participated       bool              `json:"participated"`
	FirstPrizeType     *LotteryPrizeType `json:"prize_type_first"`
	Reposted           bool              `json:"reposted"`
	Timestamp          uint64            `json:"ts"`
	UPowerRedirectURL  string            `json:"upower_redirect_url"`
	VIPBatchSign       string            `json:"vip_batch_sign"`
	VIPRedirectURL     string            `json:"vip_redirect_url"`
}

type NavAuthor struct {
	Face string `json:"face"`
	MID  uint64 `json:"-"`
	Name string `json:"name"`
}

func (author *NavAuthor) UnmarshalJSON(data []byte) error {
	var value struct {
		Face string          `json:"face"`
		MID  json.RawMessage `json:"mid"`
		Name string          `json:"name"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	mid, err := decodeUint64(value.MID)
	if err != nil {
		return fmt.Errorf("dynamic nav author mid: %w", err)
	}
	*author = NavAuthor{Face: value.Face, MID: mid, Name: value.Name}
	return nil
}

type NavItem struct {
	Author      NavAuthor `json:"author"`
	Cover       string    `json:"cover"`
	ID          string    `json:"id_str"`
	PublishTime string    `json:"pub_time"`
	RID         uint64    `json:"-"`
	Title       string    `json:"title"`
	Type        uint8     `json:"type"`
	Visible     bool      `json:"visible"`
}

func (item *NavItem) UnmarshalJSON(data []byte) error {
	var value struct {
		Author      NavAuthor       `json:"author"`
		Cover       string          `json:"cover"`
		ID          string          `json:"id_str"`
		PublishTime string          `json:"pub_time"`
		RID         json.RawMessage `json:"rid"`
		Title       string          `json:"title"`
		Type        uint8           `json:"type"`
		Visible     bool            `json:"visible"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	rid, err := decodeUint64(value.RID)
	if err != nil {
		return fmt.Errorf("dynamic nav rid: %w", err)
	}
	*item = NavItem{Author: value.Author, Cover: value.Cover, ID: value.ID, PublishTime: value.PublishTime, RID: rid, Title: value.Title, Type: value.Type, Visible: value.Visible}
	return nil
}

type NavFeed struct {
	HasMore        bool      `json:"has_more"`
	Items          []NavItem `json:"items"`
	Offset         string    `json:"offset"`
	UpdateBaseline string    `json:"update_baseline"`
	UpdateNum      uint64    `json:"-"`
}

func (feed *NavFeed) UnmarshalJSON(data []byte) error {
	var value struct {
		HasMore        bool            `json:"has_more"`
		Items          []NavItem       `json:"items"`
		Offset         string          `json:"offset"`
		UpdateBaseline string          `json:"update_baseline"`
		UpdateNum      json.RawMessage `json:"update_num"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	updateNum, err := decodeUint64(value.UpdateNum)
	if err != nil {
		return fmt.Errorf("dynamic nav update_num: %w", err)
	}
	*feed = NavFeed{HasMore: value.HasMore, Items: value.Items, Offset: value.Offset, UpdateBaseline: value.UpdateBaseline, UpdateNum: updateNum}
	return nil
}

type LiveUser struct {
	Face  string `json:"face"`
	Link  string `json:"link"`
	Title string `json:"title"`
	UID   uint64 `json:"uid"`
	Name  string `json:"uname"`
}

type LiveUsers struct {
	Count uint64     `json:"count"`
	Group string     `json:"group"`
	Items []LiveUser `json:"items"`
}

type UpUser struct {
	Profile UpUserProfile `json:"user_profile"`
}

type UpUserProfile struct {
	Info UpUserInfo `json:"info"`
}

type UpUserInfo struct {
	UID  uint64 `json:"uid"`
	Name string `json:"uname"`
	Face string `json:"face"`
}

type UpUsers struct {
	ButtonStatement string   `json:"button_statement"`
	Items           []UpUser `json:"items"`
}

type RecentUp struct {
	LiveUsers json.RawMessage `json:"live_users"`
	MyInfo    *MyInfo         `json:"my_info"`
	Items     []RecentUpUser  `json:"up_list"`
}

type MyInfo struct {
	Dynamics  int64  `json:"-"`
	Face      string `json:"face"`
	Follower  string `json:"follower"`
	Following int64  `json:"-"`
	MID       int64  `json:"-"`
	Name      string `json:"name"`
	SpaceBG   string `json:"space_bg"`
}

func (info *MyInfo) UnmarshalJSON(data []byte) error {
	var value struct {
		Dynamics  json.RawMessage `json:"dyns"`
		Face      string          `json:"face"`
		Follower  string          `json:"follower"`
		Following json.RawMessage `json:"following"`
		MID       json.RawMessage `json:"mid"`
		Name      string          `json:"name"`
		SpaceBG   string          `json:"space_bg"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	dynamics, err := decodeInt64(value.Dynamics)
	if err != nil {
		return fmt.Errorf("dynamic my_info dyns: %w", err)
	}
	following, err := decodeInt64(value.Following)
	if err != nil {
		return fmt.Errorf("dynamic my_info following: %w", err)
	}
	mid, err := decodeInt64(value.MID)
	if err != nil {
		return fmt.Errorf("dynamic my_info mid: %w", err)
	}
	*info = MyInfo{Dynamics: dynamics, Face: value.Face, Follower: value.Follower, Following: following, MID: mid, Name: value.Name, SpaceBG: value.SpaceBG}
	return nil
}

type RecentUpUser struct {
	Face            string `json:"face"`
	HasUpdate       bool   `json:"has_update"`
	IsReserveRecall bool   `json:"is_reserve_recall"`
	MID             int64  `json:"mid"`
	Name            string `json:"uname"`
}

func decodeUint64(raw json.RawMessage) (uint64, error) {
	value, err := decodeNumericText(raw)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(value, 10, 64)
}

func decodeInt64(raw json.RawMessage) (int64, error) {
	value, err := decodeNumericText(raw)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

func decodeNumericText(raw json.RawMessage) (string, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return "0", nil
	}
	if raw[0] == '"' {
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return "", err
		}
		return value, nil
	}
	return string(raw), nil
}
