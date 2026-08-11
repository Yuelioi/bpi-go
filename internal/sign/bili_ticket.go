package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

const (
	ticketKeyID = "ec02"
	ticketKey   = "XgwSnGZ1p"
)

// HexSign returns the lowercase HMAC-SHA256 digest for a Bili-ticket
// timestamp message.
func HexSign(key string, timestamp uint64) (string, error) {
	if key == "" {
		return "", &InvalidError{Field: "key", Message: "HMAC key cannot be empty"}
	}
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte("ts" + strconv.FormatUint(timestamp, 10)))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// TicketHexSign signs a timestamp using the Bili-ticket Web key.
func TicketHexSign(timestamp uint64) string {
	value, _ := HexSign(ticketKey, timestamp)
	return value
}

// TicketRequestParams builds the required Bili-ticket request parameters.
func TicketRequestParams(timestamp uint64, csrf string) map[string]string {
	return map[string]string{
		"key_id":      ticketKeyID,
		"hexsign":     TicketHexSign(timestamp),
		"context[ts]": strconv.FormatUint(timestamp, 10),
		"csrf":        csrf,
	}
}
