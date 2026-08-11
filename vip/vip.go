// Package vip contains parameters and stable response models for Bilibili VIP
// center endpoints.
package vip

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type CenterParams struct{ build uint32 }

func NewCenterParams() CenterParams { return CenterParams{} }

func (p CenterParams) WithBuild(build uint32) CenterParams {
	p.build = build
	return p
}

func (p CenterParams) EncodeQuery() (url.Values, error) {
	return url.Values{"build": {strconv.FormatUint(uint64(p.build), 10)}}, nil
}

type Center struct {
	User     User   `json:"user"`
	Wallet   Wallet `json:"wallet"`
	InReview bool   `json:"in_review"`
}

type User struct {
	Account              json.RawMessage `json:"account"`
	VIP                  json.RawMessage `json:"vip"`
	TV                   *TVVIP          `json:"tv"`
	BackgroundImageSmall string          `json:"background_image_small"`
	BackgroundImageBig   string          `json:"background_image_big"`
	PanelTitle           string          `json:"panel_title"`
	VIPOverdueExplain    string          `json:"vip_overdue_explain"`
	TVOverdueExplain     string          `json:"tv_overdue_explain"`
	AccountExceptionText string          `json:"account_exception_text"`
	AutoRenew            bool            `json:"is_auto_renew"`
	TVAutoRenew          bool            `json:"is_tv_auto_renew"`
	SurplusSeconds       uint64          `json:"surplus_seconds"`
	VIPKeepTime          uint64          `json:"vip_keep_time"`
}

type TVVIP struct {
	Type       uint32 `json:"type"`
	PayType    uint32 `json:"vip_pay_type"`
	Status     uint32 `json:"status"`
	Expiration uint64 `json:"due_date"`
}

type Wallet struct {
	Coupon            uint64 `json:"coupon"`
	Point             uint64 `json:"point"`
	PrivilegeReceived bool   `json:"privilege_received"`
}
