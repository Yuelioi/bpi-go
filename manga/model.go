package manga

import "encoding/json"

type PointInfo struct {
	Point       int32  `json:"point"`
	OriginPoint int32  `json:"origin_point"`
	IsActivity  bool   `json:"is_activity"`
	Title       string `json:"title"`
}

type ClockInInfo struct {
	DayCount       int32       `json:"day_count"`
	Status         int32       `json:"status"`
	Points         []int32     `json:"points"`
	CreditIcon     string      `json:"credit_icon"`
	SignBeforeIcon string      `json:"sign_before_icon"`
	SignTodayIcon  string      `json:"sign_today_icon"`
	BreatheIcon    string      `json:"breathe_icon"`
	NewCreditIcon  string      `json:"new_credit_x_icon"`
	CouponPicture  string      `json:"coupon_pic"`
	PointInfos     []PointInfo `json:"point_infos"`
}

type Coupon struct {
	ID           int32  `json:"ID"`
	RemainAmount int32  `json:"remain_amount"`
	TotalAmount  uint32 `json:"total_amount"`
}

type CouponInfo struct {
	RemainCoupon     int64 `json:"remain_coupon"`
	RemainSilver     int64 `json:"remain_silver"`
	RemainShopCoupon int64 `json:"remain_shop_coupon"`
}

type Coupons struct {
	TotalRemainAmount int32      `json:"total_remain_amount"`
	Items             []Coupon   `json:"user_coupons"`
	Info              CouponInfo `json:"coupon_info"`
}

type ProductLimit struct {
	Type  int32  `json:"type"`
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type Product struct {
	ID                 int64          `json:"id"`
	Type               int32          `json:"type"`
	Title              string         `json:"title"`
	Image              string         `json:"image"`
	Amount             int32          `json:"amount"`
	Cost               int32          `json:"cost"`
	RealCost           int32          `json:"real_cost"`
	RemainAmount       int32          `json:"remain_amount"`
	ComicID            int64          `json:"comic_id"`
	Limits             []ProductLimit `json:"limits"`
	Discount           int32          `json:"discount"`
	ProductType        int32          `json:"product_type"`
	PendantURL         string         `json:"pendant_url"`
	PendantExpire      int32          `json:"pendant_expire"`
	ExchangeLimit      int32          `json:"exchange_limit"`
	AddressDeadline    string         `json:"address_deadline"`
	ActivityType       int32          `json:"act_type"`
	HasExchanged       bool           `json:"has_exchanged"`
	MainCouponDeadline string         `json:"main_coupon_deadline"`
	Deadline           string         `json:"deadline"`
	Point              string         `json:"point"`
}

type UserPoint struct {
	Point string `json:"point"`
}

// The season's task, welfare, text, and rank schemas are intentionally kept
// raw: the promoted fixture establishes no stable fields beyond their shape.
type SeasonInfo struct {
	CurrentTime  string            `json:"current_time"`
	StartTime    string            `json:"start_time"`
	EndTime      string            `json:"end_time"`
	RemainAmount int32             `json:"remain_amount"`
	SeasonID     string            `json:"season_id"`
	Tasks        []json.RawMessage `json:"tasks"`
	Welfare      []json.RawMessage `json:"welfare"`
	Cover        string            `json:"cover"`
	TodayTasks   []json.RawMessage `json:"today_tasks"`
	Text         json.RawMessage   `json:"text"`
	SeasonTitle  string            `json:"season_title"`
	Rank         json.RawMessage   `json:"rank"`
}
