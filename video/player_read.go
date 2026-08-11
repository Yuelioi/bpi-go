package video

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// OnlineTotalParams identifies one video part for video.online_total.
type OnlineTotalParams struct {
	id  identifier
	cid ids.CID
}

func OnlineTotalByAID(aid ids.AID, cid ids.CID) OnlineTotalParams {
	return OnlineTotalParams{id: byAID(aid), cid: cid}
}

func OnlineTotalByBVID(bvid ids.BVID, cid ids.CID) OnlineTotalParams {
	return OnlineTotalParams{id: byBVID(bvid), cid: cid}
}

func (p OnlineTotalParams) EncodeQuery() (url.Values, error) {
	values, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	if err := p.cid.Validate(); err != nil {
		return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID must be non-zero"}
	}
	values.Set("cid", p.cid.String())
	return values, nil
}

// RelatedParams identifies a video for video.related_videos.
type RelatedParams struct{ id identifier }

func RelatedByAID(aid ids.AID) RelatedParams             { return RelatedParams{id: byAID(aid)} }
func RelatedByBVID(bvid ids.BVID) RelatedParams          { return RelatedParams{id: byBVID(bvid)} }
func (p RelatedParams) EncodeQuery() (url.Values, error) { return p.id.encode() }

// TagsParams identifies a video, and optionally a part, for video.tags.
type TagsParams struct {
	id  identifier
	cid ids.CID
}

func TagsByAID(aid ids.AID) TagsParams    { return TagsParams{id: byAID(aid)} }
func TagsByBVID(bvid ids.BVID) TagsParams { return TagsParams{id: byBVID(bvid)} }

func (p TagsParams) WithCID(cid ids.CID) TagsParams {
	p.cid = cid
	return p
}

func (p TagsParams) EncodeQuery() (url.Values, error) {
	values, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	if p.cid != 0 {
		if err := p.cid.Validate(); err != nil {
			return nil, &bpierr.ParameterError{Field: "cid", Message: "content ID must be non-zero"}
		}
		values.Set("cid", p.cid.String())
	}
	return values, nil
}

// InteractiveInfoParams identifies one interactive-video graph and optional
// edge for video.interactive_video_info.
type InteractiveInfoParams struct {
	id           identifier
	graphVersion uint64
	edgeID       uint64
}

func InteractiveInfoByAID(aid ids.AID, graphVersion uint64) (InteractiveInfoParams, error) {
	return newInteractiveInfoParams(byAID(aid), graphVersion)
}

func InteractiveInfoByBVID(bvid ids.BVID, graphVersion uint64) (InteractiveInfoParams, error) {
	return newInteractiveInfoParams(byBVID(bvid), graphVersion)
}

func newInteractiveInfoParams(id identifier, graphVersion uint64) (InteractiveInfoParams, error) {
	if graphVersion == 0 {
		return InteractiveInfoParams{}, &bpierr.ParameterError{Field: "graph_version", Message: "value must be non-zero"}
	}
	return InteractiveInfoParams{id: id, graphVersion: graphVersion}, nil
}

func (p InteractiveInfoParams) WithEdgeID(edgeID uint64) (InteractiveInfoParams, error) {
	if edgeID == 0 {
		return InteractiveInfoParams{}, &bpierr.ParameterError{Field: "edge_id", Message: "value must be non-zero"}
	}
	p.edgeID = edgeID
	return p, nil
}

func (p InteractiveInfoParams) EncodeQuery() (url.Values, error) {
	values, err := p.id.encode()
	if err != nil {
		return nil, err
	}
	if p.graphVersion == 0 {
		return nil, &bpierr.ParameterError{Field: "graph_version", Message: "value must be non-zero"}
	}
	values.Set("graph_version", strconv.FormatUint(p.graphVersion, 10))
	if p.edgeID != 0 {
		values.Set("edge_id", strconv.FormatUint(p.edgeID, 10))
	}
	return values, nil
}

type OnlineTotal struct {
	Total      string            `json:"total"`
	Count      string            `json:"count"`
	ShowSwitch OnlineTotalSwitch `json:"show_switch"`
}

type OnlineTotalSwitch struct {
	Total bool `json:"total"`
	Count bool `json:"count"`
}

type VideoTag struct {
	ID      *uint64 `json:"tag_id"`
	Name    string  `json:"tag_name"`
	MusicID *string `json:"music_id"`
	Type    string  `json:"tag_type"`
	JumpURL *string `json:"jump_url"`
}

type Rights struct {
	BP            uint8 `json:"bp"`
	Electric      uint8 `json:"elec"`
	Download      uint8 `json:"download"`
	Movie         uint8 `json:"movie"`
	Pay           uint8 `json:"pay"`
	HD5           uint8 `json:"hd5"`
	NoReprint     uint8 `json:"no_reprint"`
	Autoplay      uint8 `json:"autoplay"`
	UGCPay        uint8 `json:"ugc_pay"`
	IsCooperation uint8 `json:"is_cooperation"`
	UGCPayPreview uint8 `json:"ugc_pay_preview"`
	NoBackground  uint8 `json:"no_background"`
}

type Dimension struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
	Rotate uint8  `json:"rotate"`
}

type RelatedVideo struct {
	AID         ids.AID   `json:"aid"`
	Videos      uint32    `json:"videos"`
	TID         uint32    `json:"tid"`
	TypeName    string    `json:"tname"`
	Copyright   uint8     `json:"copyright"`
	Picture     string    `json:"pic"`
	Title       string    `json:"title"`
	PublishTime uint64    `json:"pubdate"`
	CreateTime  uint64    `json:"ctime"`
	Description string    `json:"desc"`
	State       int8      `json:"state"`
	Duration    uint64    `json:"duration"`
	Rights      Rights    `json:"rights"`
	Owner       Owner     `json:"owner"`
	Stat        Stat      `json:"stat"`
	Dynamic     string    `json:"dynamic"`
	CID         ids.CID   `json:"cid"`
	Dimension   Dimension `json:"dimension"`
	BVID        ids.BVID  `json:"bvid"`
	ShortLink   string    `json:"short_link_v2"`
}

type InteractiveInfo struct {
	Title          string                 `json:"title"`
	EdgeID         uint64                 `json:"edge_id"`
	StoryList      []InteractiveStory     `json:"story_list"`
	Edges          *InteractiveEdges      `json:"edges"`
	Preload        *InteractivePreload    `json:"preload"`
	HiddenVars     []InteractiveHiddenVar `json:"hidden_vars"`
	IsLeaf         uint8                  `json:"is_leaf"`
	NoTutorial     uint8                  `json:"no_tutorial"`
	NoBacktracking uint8                  `json:"no_backtracking"`
	NoEvaluation   uint8                  `json:"no_evaluation"`
}

type InteractiveStory struct {
	NodeID    uint64 `json:"node_id"`
	EdgeID    uint64 `json:"edge_id"`
	Title     string `json:"title"`
	CID       uint64 `json:"cid"`
	StartPos  uint64 `json:"start_pos"`
	Cover     string `json:"cover"`
	IsCurrent uint8  `json:"is_current"`
	Cursor    uint64 `json:"cursor"`
}

type InteractiveEdges struct {
	Dimension *InteractiveDimension `json:"dimension"`
	Questions []InteractiveQuestion `json:"questions"`
	Skin      json.RawMessage       `json:"skin"`
}

type InteractiveDimension struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
	Rotate uint8  `json:"rotate"`
	SAR    string `json:"sar"`
}

type InteractiveQuestion struct {
	ID         uint64              `json:"id"`
	Type       uint8               `json:"type"`
	StartTimeR uint32              `json:"start_time_r"`
	Duration   int64               `json:"duration"`
	PauseVideo uint8               `json:"pause_video"`
	Title      string              `json:"title"`
	Choices    []InteractiveChoice `json:"choices"`
}

type InteractiveChoice struct {
	ID             uint64 `json:"id"`
	PlatformAction string `json:"platform_action"`
	NativeAction   string `json:"native_action"`
	Condition      string `json:"condition"`
	CID            uint64 `json:"cid"`
	Option         string `json:"option"`
	IsDefault      *uint8 `json:"is_default"`
	IsHidden       *uint8 `json:"is_hidden"`
}

type InteractivePreload struct {
	Videos []InteractivePreloadVideo `json:"video"`
}

type InteractivePreloadVideo struct {
	AID uint64 `json:"aid"`
	CID uint64 `json:"cid"`
}

type InteractiveHiddenVar struct {
	Value  int64  `json:"value"`
	ID     string `json:"id"`
	IDV2   string `json:"id_v2"`
	Type   uint8  `json:"type"`
	IsShow uint8  `json:"is_show"`
	Name   string `json:"name"`
}
