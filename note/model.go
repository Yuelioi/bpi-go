package note

import "encoding/json"

type IsForbid struct {
	ForbidNoteEntrance bool `json:"forbid_note_entrance"`
}

type PrivateArc struct {
	OID         uint64 `json:"oid"`
	OIDType     uint8  `json:"oid_type"`
	Title       string `json:"title"`
	Picture     string `json:"pic"`
	Status      uint32 `json:"status"`
	Description string `json:"desc"`
}

type Tag struct {
	CID      uint64 `json:"cid"`
	Status   uint8  `json:"status"`
	Index    uint32 `json:"index"`
	Seconds  uint32 `json:"seconds"`
	Position uint32 `json:"pos"`
}

type PrivateInfo struct {
	Arc                PrivateArc `json:"arc"`
	AuditStatus        uint8      `json:"audit_status"`
	CIDCount           uint32     `json:"cid_count"`
	Content            string     `json:"content"`
	ForbidNoteEntrance bool       `json:"forbid_note_entrance"`
	PublishReason      *string    `json:"pub_reason"`
	PublishStatus      uint8      `json:"pub_status"`
	PublishVersion     uint32     `json:"pub_version"`
	Summary            string     `json:"summary"`
	Tags               []Tag      `json:"tags"`
	Title              string     `json:"title"`
}

type PublicArc struct {
	OID         uint64 `json:"oid"`
	OIDType     uint8  `json:"oid_type"`
	Title       string `json:"title"`
	Status      uint32 `json:"status"`
	Picture     string `json:"pic"`
	Description string `json:"desc"`
}

type PublicAuthor struct {
	MID     uint64          `json:"mid"`
	Name    string          `json:"name"`
	Face    string          `json:"face"`
	Level   uint8           `json:"level"`
	VIP     json.RawMessage `json:"vip_info"`
	Pendant json.RawMessage `json:"pendant"`
}

type PublicInfo struct {
	CVID               uint64       `json:"cvid"`
	NoteID             uint64       `json:"note_id"`
	Title              string       `json:"title"`
	Summary            string       `json:"summary"`
	Content            string       `json:"content"`
	CIDCount           uint32       `json:"cid_count"`
	PublishStatus      uint8        `json:"pub_status"`
	Tags               []Tag        `json:"tags"`
	Arc                PublicArc    `json:"arc"`
	Author             PublicAuthor `json:"author"`
	ForbidNoteEntrance bool         `json:"forbid_note_entrance"`
}

type ArchiveList struct {
	NoteIDs []string `json:"noteIds"`
}

type PrivateListArc struct {
	OID         uint64  `json:"oid"`
	Status      uint8   `json:"status"`
	OIDType     uint8   `json:"oid_type"`
	AID         uint64  `json:"aid"`
	BVID        *string `json:"bvid"`
	Picture     *string `json:"pic"`
	Description *string `json:"desc"`
}

type PrivateListItem struct {
	Title              string         `json:"title"`
	Summary            string         `json:"summary"`
	ModifiedTime       string         `json:"mtime"`
	Arc                PrivateListArc `json:"arc"`
	NoteID             uint64         `json:"note_id"`
	AuditStatus        uint8          `json:"audit_status"`
	WebURL             string         `json:"web_url"`
	NoteIDString       string         `json:"note_id_str"`
	Message            string         `json:"message"`
	ForbidNoteEntrance *bool          `json:"forbid_note_entrance"`
	Likes              uint64         `json:"likes"`
	HasLiked           bool           `json:"has_like"`
}

type Page struct {
	Total uint32 `json:"total"`
	Size  uint32 `json:"size"`
	Page  uint32 `json:"num"`
}

type PrivateList struct {
	Items []PrivateListItem `json:"list"`
	Page  *Page             `json:"page"`
}

type PublicListItem struct {
	CVID        uint64       `json:"cvid"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	PublishTime string       `json:"pubtime"`
	WebURL      string       `json:"web_url"`
	Message     string       `json:"message"`
	Author      PublicAuthor `json:"author"`
	Likes       uint64       `json:"likes"`
	HasLiked    bool         `json:"has_like"`
}

type PublicArchiveList struct {
	Items          []PublicListItem `json:"list"`
	Page           *Page            `json:"page"`
	ShowPublicNote bool             `json:"show_public_note"`
	Message        string           `json:"message"`
}

type PublicUserList struct {
	Items []PublicListItem `json:"list"`
	Page  *Page            `json:"page"`
}
