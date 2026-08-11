package testutil

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// DecodeBinaryProbeFixture decodes a promoted bpi-rs probe-binary-body
// fixture into the raw HTTP response bytes and content type.
func DecodeBinaryProbeFixture(data []byte) ([]byte, string, error) {
	var fixture struct {
		BodyBase64  string `json:"body_base64"`
		ContentType string `json:"content_type"`
		Encoding    string `json:"encoding"`
		Kind        string `json:"kind"`
		Length      int    `json:"length"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, "", fmt.Errorf("decode binary fixture: %w", err)
	}
	if fixture.Kind != "binary" || fixture.Encoding != "base64" || fixture.Length < 0 {
		return nil, "", fmt.Errorf("invalid binary fixture metadata")
	}
	body, err := base64.StdEncoding.DecodeString(fixture.BodyBase64)
	if err != nil {
		return nil, "", fmt.Errorf("decode binary fixture body: %w", err)
	}
	if len(body) != fixture.Length {
		return nil, "", fmt.Errorf("binary fixture length is %d, want %d", len(body), fixture.Length)
	}
	return body, fixture.ContentType, nil
}
