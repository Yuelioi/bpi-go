package danmaku

type AdvState struct {
	Coins   uint8 `json:"coins"`
	Confirm uint8 `json:"confirm"`
	Accept  bool  `json:"accept"`
	HasBuy  bool  `json:"hasBuy"`
}

type ThumbupStat struct {
	Likes    int64  `json:"likes"`
	UserLike int32  `json:"user_like"`
	ID       string `json:"id_str"`
}

type ThumbupStats map[string]ThumbupStat
