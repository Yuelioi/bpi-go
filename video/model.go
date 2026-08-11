package video

import (
	"encoding/json"

	"github.com/Yuelioi/bpi-go/ids"
)

// View is the stable payload returned by video.view.
type View struct {
	AID     ids.AID  `json:"aid"`
	BVID    ids.BVID `json:"bvid"`
	Videos  uint32   `json:"videos"`
	Title   string   `json:"title"`
	Picture string   `json:"pic"`
	Owner   Owner    `json:"owner"`
	Stat    Stat     `json:"stat"`
	CID     ids.CID  `json:"cid"`
	Pages   []Page   `json:"pages"`
}

// Detail is the stable payload returned by video.detail.
type Detail struct {
	View    View            `json:"View"`
	Tags    []Tag           `json:"Tags"`
	Related []Related       `json:"Related"`
	Card    json.RawMessage `json:"Card"`
	Reply   json.RawMessage `json:"Reply"`
}

// Owner is a video owner's stable public identity.
type Owner struct {
	MID  ids.MID `json:"mid"`
	Name string  `json:"name"`
	Face string  `json:"face"`
}

// Stat contains stable video counters.
type Stat struct {
	AID      ids.AID `json:"aid"`
	View     uint64  `json:"view"`
	Danmaku  uint64  `json:"danmaku"`
	Reply    uint64  `json:"reply"`
	Favorite *uint64 `json:"favorite"`
	Fav      *uint64 `json:"fav"`
	Coin     uint64  `json:"coin"`
	Share    uint64  `json:"share"`
	Like     uint64  `json:"like"`
}

// Page is one part of a multi-part video.
type Page struct {
	CID      ids.CID `json:"cid"`
	Number   uint32  `json:"page"`
	Part     string  `json:"part"`
	Duration uint64  `json:"duration"`
}

// Tag is one video tag returned by the detail endpoint.
type Tag struct {
	ID      uint64 `json:"tag_id"`
	Name    string `json:"tag_name"`
	JumpURL string `json:"jump_url"`
}

// Related contains the stable subset of one related-video entry.
type Related struct {
	AID   ids.AID  `json:"aid"`
	BVID  ids.BVID `json:"bvid"`
	Title string   `json:"title"`
	CID   *ids.CID `json:"cid"`
	Owner *Owner   `json:"owner"`
	Stat  *Stat    `json:"stat"`
}
