package live

import "encoding/json"

type SubArea struct {
	ID         string `json:"id"`
	ParentID   string `json:"parent_id"`
	OldAreaID  string `json:"old_area_id"`
	Name       string `json:"name"`
	ActivityID string `json:"act_id"`
	PKStatus   string `json:"pk_status"`
	HotStatus  int32  `json:"hot_status"`
	LockStatus string `json:"lock_status"`
	Picture    string `json:"pic"`
	ParentName string `json:"parent_name"`
	AreaType   int32  `json:"area_type"`
}

type ParentArea struct {
	ID   int32     `json:"id"`
	Name string    `json:"name"`
	List []SubArea `json:"list"`
}

type RoomPendantFrame struct {
	Name              string `json:"name"`
	Value             string `json:"value"`
	Position          int32  `json:"position"`
	Description       string `json:"desc"`
	Area              int32  `json:"area"`
	OldArea           int32  `json:"area_old"`
	BackgroundColor   string `json:"bg_color"`
	BackgroundPicture string `json:"bg_pic"`
	UseOldArea        bool   `json:"use_old_area"`
}

type RoomPendantBadge struct {
	Name        string `json:"name"`
	Position    int32  `json:"position"`
	Value       string `json:"value"`
	Description string `json:"desc"`
}

type RoomPendants struct {
	Frame       RoomPendantFrame  `json:"frame"`
	MobileFrame *RoomPendantFrame `json:"mobile_frame"`
	Badge       *RoomPendantBadge `json:"badge"`
	MobileBadge *RoomPendantBadge `json:"mobile_badge"`
}

// RoomStudioInfo is intentionally open because the upstream field is not
// documented and has no stable promoted shape.
type RoomStudioInfo map[string]json.RawMessage

type RoomInfo struct {
	UID                  int64           `json:"uid"`
	RoomID               int64           `json:"room_id"`
	ShortID              int64           `json:"short_id"`
	Attention            int64           `json:"attention"`
	Online               int64           `json:"online"`
	IsPortrait           bool            `json:"is_portrait"`
	Description          string          `json:"description"`
	LiveStatus           int32           `json:"live_status"`
	AreaID               int32           `json:"area_id"`
	ParentAreaID         int32           `json:"parent_area_id"`
	ParentAreaName       string          `json:"parent_area_name"`
	OldAreaID            int32           `json:"old_area_id"`
	Background           string          `json:"background"`
	Title                string          `json:"title"`
	UserCover            string          `json:"user_cover"`
	Keyframe             string          `json:"keyframe"`
	LiveTime             string          `json:"live_time"`
	Tags                 string          `json:"tags"`
	RoomSilentType       string          `json:"room_silent_type"`
	RoomSilentLevel      int32           `json:"room_silent_level"`
	RoomSilentSecond     int64           `json:"room_silent_second"`
	AreaName             string          `json:"area_name"`
	HotWords             []string        `json:"hot_words"`
	HotWordsStatus       int32           `json:"hot_words_status"`
	NewPendants          RoomPendants    `json:"new_pendants"`
	PKStatus             int32           `json:"pk_status"`
	PKID                 int64           `json:"pk_id"`
	AllowChangeAreaTime  int64           `json:"allow_change_area_time"`
	AllowUploadCoverTime int64           `json:"allow_upload_cover_time"`
	StudioInfo           *RoomStudioInfo `json:"studio_info"`
}

type QualityDescription struct {
	QN          int32  `json:"qn"`
	Description string `json:"desc"`
}

type StreamURL struct {
	URL        string `json:"url"`
	Order      int32  `json:"order"`
	StreamType int32  `json:"stream_type"`
	P2PType    int32  `json:"p2p_type"`
}

type Stream struct {
	CurrentQuality     int32                `json:"current_quality"`
	AcceptQuality      []string             `json:"accept_quality"`
	CurrentQN          int32                `json:"current_qn"`
	QualityDescription []QualityDescription `json:"quality_description"`
	URLs               []StreamURL          `json:"durl"`
}

type WatchedShow struct {
	Switch       bool   `json:"switch"`
	Number       int32  `json:"num"`
	TextSmall    string `json:"text_small"`
	TextLarge    string `json:"text_large"`
	Icon         string `json:"icon"`
	IconLocation int32  `json:"icon_location"`
	IconWeb      string `json:"icon_web"`
}

type RecommendRoom struct {
	HeadBox              json.RawMessage `json:"head_box"`
	AreaV2ID             int32           `json:"area_v2_id"`
	AreaV2ParentID       int32           `json:"area_v2_parent_id"`
	AreaV2Name           string          `json:"area_v2_name"`
	AreaV2ParentName     string          `json:"area_v2_parent_name"`
	BroadcastType        int32           `json:"broadcast_type"`
	Cover                string          `json:"cover"`
	Link                 string          `json:"link"`
	Online               int32           `json:"online"`
	PendantInfo          json.RawMessage `json:"pendant_Info"`
	RoomID               int64           `json:"roomid"`
	Title                string          `json:"title"`
	UserName             string          `json:"uname"`
	Face                 string          `json:"face"`
	Verify               json.RawMessage `json:"verify"`
	UID                  int64           `json:"uid"`
	Keyframe             string          `json:"keyframe"`
	IsAutoPlay           int32           `json:"is_auto_play"`
	HeadBoxType          int32           `json:"head_box_type"`
	Flag                 int32           `json:"flag"`
	SessionID            string          `json:"session_id"`
	ShowCallback         string          `json:"show_callback"`
	ClickCallback        string          `json:"click_callback"`
	SpecialID            int32           `json:"special_id"`
	WatchedShow          WatchedShow     `json:"watched_show"`
	IsNFT                int32           `json:"is_nft"`
	NFTMark              string          `json:"nft_dmark"`
	IsAd                 bool            `json:"is_ad"`
	AdTransparentContent json.RawMessage `json:"ad_transparent_content"`
	ShowAdIcon           bool            `json:"show_ad_icon"`
	Status               bool            `json:"status"`
	Followers            int32           `json:"followers"`
}

type Recommend struct {
	Rooms     []RecommendRoom `json:"recommend_room_list"`
	TopRoomID int64           `json:"top_room_id"`
}

type Version struct {
	CurrentVersion   string `json:"curr_version"`
	Build            uint64 `json:"build"`
	Instruction      string `json:"instruction"`
	FileSize         string `json:"file_size"`
	FileMD5          string `json:"file_md5"`
	Content          string `json:"content"`
	DownloadURL      string `json:"download_url"`
	HDiffPatchSwitch uint8  `json:"hdiffpatch_switch"`
}

type GiftType struct {
	GiftID   int64  `json:"gift_id"`
	GiftName string `json:"gift_name"`
	Price    int64  `json:"price"`
}

type GiftItem struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Price             int64  `json:"price"`
	Type              int32  `json:"type"`
	CoinType          string `json:"coin_type"`
	Effect            int32  `json:"effect"`
	StayTime          int32  `json:"stay_time"`
	AnimationFrameNum int32  `json:"animation_frame_num"`
	Description       string `json:"desc"`
	BasicImage        string `json:"img_basic"`
	GIF               string `json:"gif"`
}

type GiftConfig struct {
	List []GiftItem `json:"list"`
}

type GiftBaseConfig struct {
	BaseConfig GiftConfig `json:"base_config"`
}

type RoomGiftList struct {
	GiftConfig   *GiftBaseConfig `json:"gift_config"`
	GiftData     json.RawMessage `json:"gift_data"`
	GlobalConfig json.RawMessage `json:"global_config"`
}

type BlindGift struct {
	GiftID    int64  `json:"gift_id"`
	Price     int64  `json:"price"`
	GiftName  string `json:"gift_name"`
	GiftImage string `json:"gift_img"`
	Chance    string `json:"chance"`
}

type BlindGiftInfo struct {
	NoteText      string      `json:"note_text"`
	BlindPrice    int64       `json:"blind_price"`
	BlindGiftName string      `json:"blind_gift_name"`
	Gifts         []BlindGift `json:"gifts"`
}

type DanmuHost struct {
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	WSSPort uint32 `json:"wss_port"`
	WSPort  uint32 `json:"ws_port"`
}

type DanmuInfo struct {
	Token string      `json:"token"`
	Hosts []DanmuHost `json:"host_list"`
}

type Emoticon struct {
	BulgeDisplay    int32  `json:"bulge_display"`
	Description     string `json:"descript"`
	Emoji           string `json:"emoji"`
	ID              int64  `json:"emoticon_id"`
	Unique          string `json:"emoticon_unique"`
	ValueType       int32  `json:"emoticon_value_type"`
	Height          int32  `json:"height"`
	Identity        int32  `json:"identity"`
	InPlayerArea    int32  `json:"in_player_area"`
	IsDynamic       int32  `json:"is_dynamic"`
	Permission      int32  `json:"perm"`
	UnlockNeedGift  int32  `json:"unlock_need_gift"`
	UnlockNeedLevel int32  `json:"unlock_need_level"`
	UnlockShowColor string `json:"unlock_show_color"`
	UnlockShowImage string `json:"unlock_show_image"`
	UnlockShowText  string `json:"unlock_show_text"`
	URL             string `json:"url"`
	Width           int32  `json:"width"`
}

type TopShowItem struct {
	Image string `json:"image"`
	Text  string `json:"text"`
}

type TopShow struct {
	TopLeft  TopShowItem `json:"top_left"`
	TopRight TopShowItem `json:"top_right"`
}

type EmoticonPackage struct {
	CurrentCover          string            `json:"current_cover"`
	Emoticons             []Emoticon        `json:"emoticons"`
	Description           string            `json:"pkg_descript"`
	ID                    int64             `json:"pkg_id"`
	Name                  string            `json:"pkg_name"`
	Permission            int32             `json:"pkg_perm"`
	Type                  int32             `json:"pkg_type"`
	RecentlyUsedEmoticons []json.RawMessage `json:"recently_used_emoticons"`
	TopShow               *TopShow          `json:"top_show"`
	TopShowRecent         *TopShow          `json:"top_show_recent"`
	UnlockIdentity        int32             `json:"unlock_identity"`
	UnlockNeedGift        int32             `json:"unlock_need_gift"`
}

type EmoticonData struct {
	Packages    []EmoticonPackage `json:"data"`
	FansBrand   int32             `json:"fans_brand"`
	PurchaseURL *string           `json:"purchase_url"`
}

type RedPocketAward struct {
	GiftID      int64  `json:"gift_id"`
	Number      int32  `json:"num"`
	GiftName    string `json:"gift_name"`
	GiftPicture string `json:"gift_pic"`
}

type PopularityRedPocket struct {
	LotteryID       int64            `json:"lot_id"`
	SenderUID       int64            `json:"sender_uid"`
	SenderName      string           `json:"sender_name"`
	SenderFace      string           `json:"sender_face"`
	JoinRequirement int32            `json:"join_requirement"`
	Danmu           string           `json:"danmu"`
	Awards          []RedPocketAward `json:"awards"`
	StartTime       int64            `json:"start_time"`
	EndTime         int64            `json:"end_time"`
	LastTime        int64            `json:"last_time"`
	RemoveTime      int64            `json:"remove_time"`
	ReplaceTime     int64            `json:"replace_time"`
	CurrentTime     int64            `json:"current_time"`
	LotteryStatus   int32            `json:"lot_status"`
	H5URL           string           `json:"h5_url"`
	UserStatus      int32            `json:"user_status"`
	LotteryConfigID int64            `json:"lot_config_id"`
	TotalPrice      int64            `json:"total_price"`
}

type ActivityBoxInfo map[string]json.RawMessage

// LotteryInfo keeps documented fields typed and preserves only the
// explicitly unstable lottery sections in Extra.
type LotteryInfo struct {
	PopularityRedPocket []PopularityRedPocket      `json:"popularity_red_pocket"`
	ActivityBoxInfo     *ActivityBoxInfo           `json:"activity_box_info"`
	Extra               map[string]json.RawMessage `json:"-"`
}

func (l *LotteryInfo) UnmarshalJSON(data []byte) error {
	type wire LotteryInfo
	var known wire
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	delete(raw, "popularity_red_pocket")
	delete(raw, "activity_box_info")
	*l = LotteryInfo(known)
	l.Extra = raw
	return nil
}

type MedalPageInfo struct {
	TotalPage   int32 `json:"total_page"`
	CurrentPage int32 `json:"cur_page"`
}

type FansMedal struct {
	CanDelete        bool   `json:"can_deleted"`
	DayLimit         int32  `json:"day_limit"`
	GuardLevel       int32  `json:"guard_level"`
	GuardMedalTitle  string `json:"guard_medal_title"`
	Intimacy         int32  `json:"intimacy"`
	IsLighted        int32  `json:"is_lighted"`
	Level            int32  `json:"level"`
	MedalName        string `json:"medal_name"`
	MedalColorBorder int32  `json:"medal_color_border"`
	MedalColorStart  int32  `json:"medal_color_start"`
	MedalColorEnd    int32  `json:"medal_color_end"`
	MedalID          int64  `json:"medal_id"`
	NextIntimacy     int32  `json:"next_intimacy"`
	TodayFeed        int32  `json:"today_feed"`
	RoomID           int64  `json:"roomid"`
	Status           int32  `json:"status"`
	TargetID         int64  `json:"target_id"`
	TargetName       string `json:"target_name"`
	UserName         string `json:"uname"`
}

type MyMedals struct {
	Count    int32         `json:"count"`
	Items    []FansMedal   `json:"items"`
	PageInfo MedalPageInfo `json:"page_info"`
}

type FollowUpLive struct {
	RoomID         int64  `json:"roomid"`
	UID            int64  `json:"uid"`
	UserName       string `json:"uname"`
	Title          string `json:"title"`
	Face           string `json:"face"`
	LiveStatus     int32  `json:"live_status"`
	RecordLiveTime int64  `json:"record_live_time"`
	AreaNameV2     string `json:"area_name_v2"`
	RoomNews       string `json:"room_news"`
	TextSmall      string `json:"text_small"`
	RoomCover      string `json:"room_cover"`
	ParentAreaID   int32  `json:"parent_area_id"`
	AreaID         int32  `json:"area_id"`
}

type FollowUpList struct {
	Title           string         `json:"title"`
	PageSize        int32          `json:"pageSize"`
	TotalPage       int32          `json:"totalPage"`
	List            []FollowUpLive `json:"list"`
	Count           int32          `json:"count"`
	NeverLivedCount int32          `json:"never_lived_count"`
	LiveCount       int32          `json:"live_count"`
	NeverLivedFaces []string       `json:"never_lived_faces"`
}

type FollowedRoom struct {
	Title            string `json:"title"`
	RoomID           int64  `json:"room_id"`
	UID              int64  `json:"uid"`
	Online           int32  `json:"online"`
	LiveTime         int64  `json:"live_time"`
	LiveStatus       int32  `json:"live_status"`
	ShortID          int32  `json:"short_id"`
	Area             int32  `json:"area"`
	AreaName         string `json:"area_name"`
	AreaV2ID         int32  `json:"area_v2_id"`
	AreaV2Name       string `json:"area_v2_name"`
	AreaV2ParentName string `json:"area_v2_parent_name"`
	AreaV2ParentID   int32  `json:"area_v2_parent_id"`
	UserName         string `json:"uname"`
	Face             string `json:"face"`
	TagName          string `json:"tag_name"`
	Tags             string `json:"tags"`
	UserCover        string `json:"cover_from_user"`
	Keyframe         string `json:"keyframe"`
	LockTill         string `json:"lock_till"`
	HiddenTill       string `json:"hidden_till"`
	BroadcastType    int32  `json:"broadcast_type"`
	IsEncrypted      bool   `json:"is_encrypt"`
	Link             string `json:"link"`
	Nickname         string `json:"nickname"`
	RoomName         string `json:"roomname"`
	LegacyRoomID     int64  `json:"roomid"`
}

type FollowUpWebList struct {
	Rooms        []FollowedRoom `json:"rooms"`
	List         []FollowedRoom `json:"list"`
	Count        int32          `json:"count"`
	NotLivingNum int32          `json:"not_living_num"`
}

type ReplayLiveInfo struct {
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	LiveTime int64  `json:"live_time"`
	LiveType int32  `json:"live_type"`
}

type ReplayVideoInfo struct {
	ReplayStatus  int32   `json:"replay_status"`
	EstimatedTime string  `json:"estimated_time"`
	Duration      int32   `json:"duration"`
	DownloadURL   *string `json:"download_url"`
	AlertCode     *int32  `json:"alert_code"`
	AlertMessage  *string `json:"alert_message"`
}

type ReplayAlarmInfo struct {
	Code         int32  `json:"code"`
	Message      string `json:"message"`
	CurrentTime  int64  `json:"cur_time"`
	IsBanPublish bool   `json:"is_ban_publish"`
}

type Replay struct {
	ReplayID  int64           `json:"replay_id"`
	LiveInfo  ReplayLiveInfo  `json:"live_info"`
	VideoInfo ReplayVideoInfo `json:"video_info"`
	AlarmInfo ReplayAlarmInfo `json:"alarm_info"`
	RoomID    int64           `json:"room_id"`
	LiveKey   string          `json:"live_key"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
}

type ReplayPagination struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Total    *int32 `json:"total"`
}

type ReplayList struct {
	ReplayInfo []Replay         `json:"replay_info"`
	Pagination ReplayPagination `json:"pagination"`
}

type GuardUserOrigin struct {
	Name string `json:"name"`
	Face string `json:"face"`
}

type GuardUserOfficial struct {
	Role        int32  `json:"role"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Type        int32  `json:"type"`
}

type GuardUserBase struct {
	Name          string            `json:"name"`
	Face          string            `json:"face"`
	NameColor     int32             `json:"name_color"`
	IsMystery     bool              `json:"is_mystery"`
	RiskControl   json.RawMessage   `json:"risk_ctrl_info"`
	OriginInfo    GuardUserOrigin   `json:"origin_info"`
	OfficialInfo  GuardUserOfficial `json:"official_info"`
	NameColorText string            `json:"name_color_str"`
}

type GuardUserMedal struct {
	Name             string `json:"name"`
	Level            int32  `json:"level"`
	ColorStart       int32  `json:"color_start"`
	ColorEnd         int32  `json:"color_end"`
	ColorBorder      int32  `json:"color_border"`
	Color            int32  `json:"color"`
	ID               int32  `json:"id"`
	Type             int32  `json:"typ"`
	IsLight          int32  `json:"is_light"`
	AnchorUID        int64  `json:"ruid"`
	GuardLevel       int32  `json:"guard_level"`
	Score            int32  `json:"score"`
	GuardIcon        string `json:"guard_icon"`
	HonorIcon        string `json:"honor_icon"`
	V2ColorStart     string `json:"v2_medal_color_start"`
	V2ColorEnd       string `json:"v2_medal_color_end"`
	V2ColorBorder    string `json:"v2_medal_color_border"`
	V2ColorText      string `json:"v2_medal_color_text"`
	V2ColorLevel     string `json:"v2_medal_color_level"`
	UserReceiveCount int32  `json:"user_receive_count"`
}

type GuardUserGuard struct {
	Level       int32  `json:"level"`
	ExpiredText string `json:"expired_str"`
}

type GuardUser struct {
	UID         int64           `json:"uid"`
	Base        GuardUserBase   `json:"base"`
	Medal       GuardUserMedal  `json:"medal"`
	Wealth      json.RawMessage `json:"wealth"`
	Title       json.RawMessage `json:"title"`
	Guard       GuardUserGuard  `json:"guard"`
	HeadFrame   json.RawMessage `json:"uhead_frame"`
	GuardLeader json.RawMessage `json:"guard_leader"`
}

type GuardMember struct {
	AnchorUID     int64     `json:"ruid"`
	Rank          int32     `json:"rank"`
	AccompanyDays int32     `json:"accompany"`
	User          GuardUser `json:"uinfo"`
	Score         int32     `json:"score"`
}

type GuardListInfo struct {
	Number                 int32    `json:"num"`
	PageCount              int32    `json:"page"`
	CurrentPage            int32    `json:"now"`
	AchievementLevel       int32    `json:"achievement_level"`
	AnchorAchievementLevel int32    `json:"anchor_guard_achieve_level"`
	AchievementIcon        string   `json:"achievement_icon_src"`
	BuyGuardIcon           string   `json:"buy_guard_icon_src"`
	RuleDocument           string   `json:"rule_doc_src"`
	Background             string   `json:"ex_background_src"`
	ColorStart             string   `json:"color_start"`
	ColorEnd               string   `json:"color_end"`
	TabColors              []string `json:"tab_color"`
	TitleColors            []string `json:"title_color"`
}

type GuardList struct {
	Info GuardListInfo `json:"info"`
	Top3 []GuardMember `json:"top3"`
	List []GuardMember `json:"list"`
}

type SilentUser struct {
	TargetUID    int64  `json:"tuid"`
	TargetName   string `json:"tname"`
	OperatorUID  int64  `json:"uid"`
	OperatorName string `json:"name"`
	CreatedAt    string `json:"ctime"`
	ID           int64  `json:"id"`
	IsAnchor     int8   `json:"is_anchor"`
	Face         string `json:"face"`
	Message      string `json:"msg"`
	AdminLevel   int8   `json:"admin_level"`
	IsMystery    bool   `json:"is_mystery"`
	BlockEndTime string `json:"block_end_time"`
	Type         int8   `json:"type"`
}

type SilentUsers struct {
	Data      []SilentUser `json:"data"`
	Total     int32        `json:"total"`
	TotalPage int32        `json:"total_page"`
	Page      int32        `json:"pn"`
	PageSize  int32        `json:"ps"`
}

type BannedUser struct {
	UID          int64  `json:"uid"`
	ModifiedAt   string `json:"mtime"`
	Face         string `json:"face"`
	Name         string `json:"name"`
	IsAnchor     bool   `json:"is_anchor"`
	OperatorName string `json:"operator_name"`
	AdminLevel   int8   `json:"admin_level"`
	IsMystery    bool   `json:"is_mystery"`
}

type BannedUsers struct {
	Data      []BannedUser `json:"data"`
	Total     int32        `json:"total"`
	TotalPage int32        `json:"total_page"`
	Page      int32        `json:"pn"`
	PageSize  int32        `json:"ps"`
}

type ShieldKeyword struct {
	Keyword  string `json:"keyword"`
	UID      int64  `json:"uid"`
	Name     string `json:"name"`
	IsAnchor int8   `json:"is_anchor"`
}

type ShieldKeywords struct {
	Keywords []ShieldKeyword `json:"keyword_list"`
	MaxLimit int32           `json:"max_limit"`
}

type HeartBeat struct {
	NextInterval int32 `json:"next_interval"`
}
