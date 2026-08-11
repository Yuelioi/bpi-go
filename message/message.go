// Package message contains parameters and stable response models for Bilibili
// notification and private-message counters.
package message

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// UnreadCountParams configures the aggregate notification-counter request.
type UnreadCountParams struct {
	build   string
	mobiApp string
}

func NewUnreadCountParams() UnreadCountParams {
	return UnreadCountParams{build: "0", mobiApp: "web"}
}

func (p UnreadCountParams) WithBuild(build string) (UnreadCountParams, error) {
	build = strings.TrimSpace(build)
	if build == "" {
		return UnreadCountParams{}, &bpierr.ParameterError{Field: "build", Message: "value cannot be blank"}
	}
	p.build = build
	return p, nil
}

func (p UnreadCountParams) WithMobiApp(mobiApp string) (UnreadCountParams, error) {
	mobiApp = strings.TrimSpace(mobiApp)
	if mobiApp == "" {
		return UnreadCountParams{}, &bpierr.ParameterError{Field: "mobi_app", Message: "value cannot be blank"}
	}
	p.mobiApp = mobiApp
	return p, nil
}

func (p UnreadCountParams) EncodeQuery() (url.Values, error) {
	return url.Values{"build": {p.build}, "mobi_app": {p.mobiApp}}, nil
}

// ReplyFeedParams configures the reply-notification feed.
type ReplyFeedParams struct {
	startID     *uint64
	startTime   *uint64
	webLocation string
}

func NewReplyFeedParams() ReplyFeedParams { return ReplyFeedParams{} }

func (p ReplyFeedParams) WithStartID(id uint64) (ReplyFeedParams, error) {
	if id == 0 {
		return ReplyFeedParams{}, &bpierr.ParameterError{Field: "id", Message: "value must be non-zero"}
	}
	p.startID = &id
	return p, nil
}

func (p ReplyFeedParams) WithStartTime(timestamp uint64) (ReplyFeedParams, error) {
	if timestamp == 0 {
		return ReplyFeedParams{}, &bpierr.ParameterError{Field: "reply_time", Message: "value must be non-zero"}
	}
	p.startTime = &timestamp
	return p, nil
}

func (p ReplyFeedParams) WithWebLocation(webLocation string) ReplyFeedParams {
	p.webLocation = webLocation
	return p
}

func (p ReplyFeedParams) EncodeQuery() (url.Values, error) {
	values := url.Values{
		"build":        {"0"},
		"mobi_app":     {"web"},
		"platform":     {"web"},
		"web_location": {p.webLocation},
	}
	if p.startID != nil {
		values.Set("id", strconv.FormatUint(*p.startID, 10))
	}
	if p.startTime != nil {
		values.Set("reply_time", strconv.FormatUint(*p.startTime, 10))
	}
	return values, nil
}

type SingleUnreadType uint32

const (
	SingleUnreadAll      SingleUnreadType = 0
	SingleUnreadFollow   SingleUnreadType = 1
	SingleUnreadUnfollow SingleUnreadType = 2
	SingleUnreadBlocked  SingleUnreadType = 3
)

// SingleUnreadParams configures private-message unread counters.
type SingleUnreadParams struct {
	unreadType   SingleUnreadType
	showUnfollow bool
	showDustbin  bool
}

func NewSingleUnreadParams() SingleUnreadParams { return SingleUnreadParams{} }

func (p SingleUnreadParams) WithUnreadType(unreadType SingleUnreadType) SingleUnreadParams {
	p.unreadType = unreadType
	if unreadType == SingleUnreadBlocked {
		p.showDustbin = true
	}
	return p
}

func (p SingleUnreadParams) WithCustomUnreadType(unreadType uint32) (SingleUnreadParams, error) {
	if unreadType == 0 {
		return SingleUnreadParams{}, &bpierr.ParameterError{Field: "unread_type", Message: "value must be non-zero"}
	}
	p.unreadType = SingleUnreadType(unreadType)
	return p, nil
}

func (p SingleUnreadParams) ShowUnfollowList(show bool) SingleUnreadParams {
	p.showUnfollow = show
	return p
}

func (p SingleUnreadParams) ShowDustbin(show bool) SingleUnreadParams {
	p.showDustbin = show
	return p
}

func (p SingleUnreadParams) EncodeQuery() (url.Values, error) {
	return url.Values{
		"build":              {"0"},
		"mobi_app":           {"web"},
		"show_dustbin":       {boolFlag(p.showDustbin)},
		"show_unfollow_list": {boolFlag(p.showUnfollow)},
		"unread_type":        {strconv.FormatUint(uint64(p.unreadType), 10)},
	}, nil
}

func boolFlag(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

type UnreadCount struct {
	Coin          uint32 `json:"coin"`
	Danmaku       uint32 `json:"danmu"`
	Favorite      uint32 `json:"favorite"`
	ReceivedLike  uint32 `json:"recv_like"`
	ReceivedReply uint32 `json:"recv_reply"`
	System        uint32 `json:"sys_msg"`
	UP            uint32 `json:"up"`
}

type ReplyFeed struct {
	Cursor     ReplyCursor `json:"cursor"`
	Items      []ReplyItem `json:"items"`
	LastViewAt uint64      `json:"last_view_at"`
}

type ReplyCursor struct {
	End  bool    `json:"is_end"`
	ID   *uint64 `json:"id"`
	Time *uint64 `json:"time"`
}

type ReplyItem struct {
	ID        uint64      `json:"id"`
	User      ReplyUser   `json:"user"`
	Detail    ReplyDetail `json:"item"`
	Count     uint32      `json:"counts"`
	Multi     uint32      `json:"is_multi"`
	ReplyTime uint64      `json:"reply_time"`
}

type ReplyUser struct {
	MID      uint64  `json:"mid"`
	Nickname string  `json:"nickname"`
	Avatar   string  `json:"avatar"`
	Follow   bool    `json:"follow"`
	Fans     *uint32 `json:"fans"`
	MIDLink  *string `json:"mid_link"`
}

type ReplyDetail struct {
	SubjectID          uint64         `json:"subject_id"`
	RootID             uint64         `json:"root_id"`
	SourceID           uint64         `json:"source_id"`
	TargetID           uint64         `json:"target_id"`
	Type               string         `json:"type"`
	BusinessID         uint32         `json:"business_id"`
	Business           string         `json:"business"`
	Title              string         `json:"title"`
	Description        string         `json:"desc"`
	URI                string         `json:"uri"`
	NativeURI          string         `json:"native_uri"`
	RootReplyContent   string         `json:"root_reply_content"`
	SourceContent      string         `json:"source_content"`
	TargetReplyContent string         `json:"target_reply_content"`
	AtDetails          []AtUserDetail `json:"at_details"`
	HideReplyButton    bool           `json:"hide_reply_button"`
	HideLikeButton     bool           `json:"hide_like_button"`
	LikeState          uint32         `json:"like_state"`
}

type AtUserDetail struct {
	MID      uint64 `json:"mid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Follow   bool   `json:"follow"`
}

type SingleUnread struct {
	UnfollowUnread         uint32 `json:"unfollow_unread"`
	FollowUnread           uint32 `json:"follow_unread"`
	UnfollowPushMessage    uint32 `json:"unfollow_push_msg"`
	DustbinPushMessage     uint32 `json:"dustbin_push_msg"`
	DustbinUnread          uint32 `json:"dustbin_unread"`
	BusinessUnfollowUnread uint32 `json:"biz_msg_unfollow_unread"`
	BusinessFollowUnread   uint32 `json:"biz_msg_follow_unread"`
	CustomUnread           uint32 `json:"custom_unread"`
}
