package electric

import "encoding/json"

type ChargeVIPInfo struct {
	DueMilliseconds int64 `json:"vipDueMsec"`
	Status          int32 `json:"vipStatus"`
	Type            int32 `json:"vipType"`
}

type ChargeUser struct {
	Name    string        `json:"uname"`
	Avatar  string        `json:"avatar"`
	MID     int64         `json:"mid"`
	PayMID  int64         `json:"pay_mid"`
	Rank    int32         `json:"rank"`
	VIP     ChargeVIPInfo `json:"vip_info"`
	Message string        `json:"message"`
}

type MonthUpList struct {
	Count      int32        `json:"count"`
	List       []ChargeUser `json:"list"`
	TotalCount int32        `json:"total_count"`
}

type VideoShowHighLevel struct {
	PrivilegeType int32  `json:"privilege_type"`
	Title         string `json:"title"`
	Subtitle      string `json:"sub_title"`
	ShowButton    bool   `json:"show_button"`
}

type VideoShowInfo struct {
	Show       bool               `json:"show"`
	State      int32              `json:"state"`
	Title      string             `json:"title"`
	JumpURL    string             `json:"jump_url"`
	Icon       string             `json:"icon"`
	HighLevel  VideoShowHighLevel `json:"high_level"`
	QuestionID int64              `json:"with_qa_id"`
}

type VideoShow struct {
	ShowInfo   VideoShowInfo `json:"show_info"`
	AVCount    int32         `json:"av_count"`
	Count      int32         `json:"count"`
	TotalCount int32         `json:"total_count"`
	List       []ChargeUser  `json:"list"`
}

type RechargePage struct {
	CurrentPage uint64 `json:"currentPage"`
	PageSize    uint64 `json:"pageSize"`
	TotalCount  uint64 `json:"totalCount"`
	TotalPage   uint64 `json:"totalPage"`
}

type RechargeRecord struct {
	MID               uint64  `json:"mid"`
	Name              string  `json:"name"`
	Avatar            string  `json:"avatar"`
	OriginalThirdCoin float64 `json:"originalThirdCoin"`
	Brokerage         float64 `json:"brokerage"`
	Remark            string  `json:"remark"`
	CreatedAt         string  `json:"ctime"`
}

type RechargeList struct {
	Page   RechargePage     `json:"page"`
	Result []RechargeRecord `json:"result"`
}

type RankPager struct {
	Current uint64 `json:"current"`
	Size    uint64 `json:"size"`
	Total   uint64 `json:"total"`
}

type RecentRankRecord struct {
	AID       uint64  `json:"aid"`
	BVID      string  `json:"bvid"`
	Electric  float64 `json:"elec_num"`
	Title     string  `json:"title"`
	Name      string  `json:"uname"`
	Avatar    string  `json:"avatar"`
	CreatedAt string  `json:"ctime"`
}

type RecentRank struct {
	List  []RecentRankRecord `json:"list"`
	Pager RankPager          `json:"pager"`
}

type Renew struct {
	UID             uint64 `json:"uid"`
	RUID            uint64 `json:"ruid"`
	GoodsID         uint64 `json:"goods_id"`
	Status          uint8  `json:"status"`
	NextExecuteTime uint64 `json:"next_execute_time"`
	SignedTime      uint64 `json:"signed_time"`
	SignedPrice     uint64 `json:"signed_price"`
	PayChannel      uint8  `json:"pay_channel"`
	Period          uint64 `json:"period"`
	MobileApp       string `json:"mobile_app"`
}

type ChargeItem struct {
	PrivilegeType uint64  `json:"privilege_type"`
	Icon          string  `json:"icon"`
	Name          string  `json:"name"`
	ExpireTime    uint64  `json:"expire_time"`
	Renew         *Renew  `json:"renew"`
	StartTime     uint64  `json:"start_time"`
	RenewList     []Renew `json:"renew_list"`
}

type ChargeUP struct {
	UPMID          uint64       `json:"up_uid"`
	Name           string       `json:"user_name"`
	Face           string       `json:"user_face"`
	Items          []ChargeItem `json:"item"`
	Start          uint64       `json:"start"`
	HighLevelState uint8        `json:"high_level_state"`
	ReplyState     uint8        `json:"elec_reply_state"`
}

type ChargeRecord struct {
	List       []ChargeUP `json:"list"`
	Page       uint64     `json:"page"`
	PageSize   uint64     `json:"page_size"`
	TotalPage  uint64     `json:"total_page"`
	TotalCount uint64     `json:"total_num"`
	HasMore    uint8      `json:"is_more"`
}

type RankUser struct {
	Rank     uint64 `json:"rank"`
	MID      uint64 `json:"mid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type UpowerRank struct {
	Total       uint64     `json:"total"`
	Description string     `json:"total_desc"`
	List        []RankUser `json:"list"`
}

type ItemIntro struct {
	IntroVideoAID string `json:"intro_video_aid"`
	Welcome       string `json:"welcomes"`
}

type UserCard struct {
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type ItemDetail struct {
	Rank              UpowerRank        `json:"upower_rank"`
	Item              ItemIntro         `json:"item"`
	UserCard          UserCard          `json:"user_card"`
	Level             uint8             `json:"upower_level"`
	ReplyState        uint8             `json:"elec_reply_state"`
	VoucherState      json.RawMessage   `json:"voucher_state"`
	RightCounts       map[string]uint64 `json:"upower_right_count"`
	OnlyContainsMedal bool              `json:"only_contain_medal"`
	PrivilegeType     uint64            `json:"privilege_type"`
}

type UPCard struct {
	MID           uint64 `json:"mid"`
	Nickname      string `json:"nickname"`
	OfficialTitle string `json:"official_title"`
	Avatar        string `json:"avatar"`
}

type ChallengeInfo struct {
	ID               string            `json:"challenge_id"`
	Description      string            `json:"description"`
	Type             int64             `json:"challenge_type"`
	RemainingDays    int64             `json:"remaining_days"`
	EndTime          string            `json:"end_time"`
	Progress         int64             `json:"progress"`
	Targets          []json.RawMessage `json:"targets"`
	State            int64             `json:"state"`
	EndTimeUnix      int64             `json:"end_time_unix"`
	PublishedDynamic int64             `json:"pub_dyn"`
	DynamicContent   string            `json:"dyn_content"`
}

type FollowInfo struct {
	Days              uint64        `json:"days"`
	UPCard            UPCard        `json:"up_card"`
	UserCard          UserCard      `json:"user_card"`
	RemainingDays     int64         `json:"remain_days"`
	RemainingLessDay  uint8         `json:"remain_less_1day"`
	Rank              UpowerRank    `json:"upower_rank"`
	Icon              string        `json:"upower_icon"`
	RightCount        int64         `json:"upower_right_count"`
	OnlyContainsMedal bool          `json:"only_contain_medal"`
	PrivilegeType     uint64        `json:"privilege_type"`
	Challenge         ChallengeInfo `json:"challenge_info"`
}

type UPInfo struct {
	MID         uint64 `json:"mid"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Type        int32  `json:"type"`
	Title       string `json:"title"`
	UpowerState uint8  `json:"upower_state"`
}

type MemberInfo struct {
	MID           uint64 `json:"mid"`
	Nickname      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	Rank          int64  `json:"rank"`
	Day           uint64 `json:"day"`
	ExpiresAt     uint64 `json:"expire_at"`
	RemainingDays uint64 `json:"remain_days"`
}

type LevelInfo struct {
	PrivilegeType uint64 `json:"privilege_type"`
	Name          string `json:"name"`
	Price         uint64 `json:"price"`
	MemberTotal   uint64 `json:"member_total"`
}

type MemberRank struct {
	UPInfo        UPInfo       `json:"up_info"`
	RankInfo      []MemberInfo `json:"rank_info"`
	UserInfo      MemberInfo   `json:"user_info"`
	MemberTotal   uint64       `json:"member_total"`
	PrivilegeType uint64       `json:"privilege_type"`
	Charged       bool         `json:"is_charge"`
	Tabs          []uint64     `json:"tabs"`
	Levels        []LevelInfo  `json:"level_info"`
}

type RemarkRecord struct {
	AID          uint64 `json:"aid"`
	BVID         string `json:"bvid"`
	ID           uint64 `json:"id"`
	MID          uint64 `json:"mid"`
	ReplyMID     uint64 `json:"reply_mid"`
	Electric     uint64 `json:"elec_num"`
	State        uint8  `json:"state"`
	Message      string `json:"msg"`
	ArchiveName  string `json:"aname"`
	UserName     string `json:"uname"`
	Avatar       string `json:"avator"`
	ReplyName    string `json:"reply_name"`
	ReplyAvatar  string `json:"reply_avator"`
	ReplyMessage string `json:"reply_msg"`
	CreatedAt    uint64 `json:"ctime"`
	ReplyTime    uint64 `json:"reply_time"`
}

type RemarkList struct {
	List  []RemarkRecord `json:"list"`
	Pager RankPager      `json:"pager"`
}

type RemarkDetail = RemarkRecord
