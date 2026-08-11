// Package login contains parameters and response models for Bilibili login
// and authenticated-session state.
package login

import (
	"encoding/json"
	"strings"

	"github.com/Yuelioi/bpi-go/ids"
)

// Nav is the stable login/navigation state returned by login.nav.
type Nav struct {
	IsLogin  bool     `json:"isLogin"`
	MID      *ids.MID `json:"mid,omitempty"`
	Username *string  `json:"uname,omitempty"`
	Face     *string  `json:"face,omitempty"`
	WBIImage WBIImage `json:"wbi_img"`
}

// WBIImage contains the current image-key URLs used for WBI signing.
type WBIImage struct {
	ImageURL string `json:"img_url"`
	SubURL   string `json:"sub_url"`
}

// UnmarshalJSON normalizes guest zero/blank identity fields to nil while
// retaining strict validation for non-zero member IDs.
func (n *Nav) UnmarshalJSON(data []byte) error {
	var raw struct {
		IsLogin  bool     `json:"isLogin"`
		MID      uint64   `json:"mid"`
		Username string   `json:"uname"`
		Face     string   `json:"face"`
		WBIImage WBIImage `json:"wbi_img"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*n = Nav{IsLogin: raw.IsLogin, WBIImage: raw.WBIImage}
	if raw.MID != 0 {
		mid, err := ids.NewMID(raw.MID)
		if err != nil {
			return err
		}
		n.MID = &mid
	}
	if strings.TrimSpace(raw.Username) != "" {
		username := raw.Username
		n.Username = &username
	}
	if strings.TrimSpace(raw.Face) != "" {
		face := raw.Face
		n.Face = &face
	}
	return nil
}
