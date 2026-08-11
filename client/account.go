package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

const redactedValue = "<redacted>"

// Account contains the four common Cookie values used for a complete Bilibili
// account projection. Raw Cookie headers may contain fewer or additional
// pairs. Account's formatted, logged, and JSON representations are redacted.
type Account struct {
	DedeUserID string
	SESSDATA   string
	BiliJCT    string
	Buvid3     string
}

// Validate reports whether every field in the complete account projection is
// present.
func (a Account) Validate() error {
	if a.DedeUserID == "" || a.SESSDATA == "" || a.BiliJCT == "" || a.Buvid3 == "" {
		return &ParameterError{
			Field:   "account",
			Message: "account requires DedeUserID, SESSDATA, bili_jct, and buvid3",
		}
	}
	return nil
}

// CSRF returns the account CSRF token.
func (a Account) CSRF() (string, error) {
	if a.BiliJCT == "" {
		return "", ErrAuthenticationRequired
	}
	return a.BiliJCT, nil
}

// String deliberately excludes all credential values.
func (a Account) String() string {
	return fmt.Sprintf("Account{DedeUserID:%s SESSDATA:%s BiliJCT:%s Buvid3:%s}",
		redact(a.DedeUserID), redact(a.SESSDATA), redact(a.BiliJCT), redact(a.Buvid3))
}

// LogValue deliberately excludes all credential values.
func (a Account) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("dede_user_id", redact(a.DedeUserID)),
		slog.String("sessdata", redact(a.SESSDATA)),
		slog.String("bili_jct", redact(a.BiliJCT)),
		slog.String("buvid3", redact(a.Buvid3)),
	)
}

// MarshalJSON deliberately excludes all credential values.
func (a Account) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		DedeUserID string `json:"dede_user_id"`
		SESSDATA   string `json:"sessdata"`
		BiliJCT    string `json:"bili_jct"`
		Buvid3     string `json:"buvid3"`
	}{
		DedeUserID: redact(a.DedeUserID),
		SESSDATA:   redact(a.SESSDATA),
		BiliJCT:    redact(a.BiliJCT),
		Buvid3:     redact(a.Buvid3),
	})
}

func redact(value string) string {
	if value == "" {
		return "<empty>"
	}
	return redactedValue
}
