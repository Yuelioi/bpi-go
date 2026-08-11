package sign_test

import (
	"testing"

	"github.com/Yuelioi/bpi-go/internal/sign"
)

const (
	imageKey = "abcdefghijklmnopqrstuvwxyz123456"
	subKey   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ654321"
)

func TestMixinKeyMatchesRustVector(t *testing.T) {
	t.Parallel()

	keys, err := sign.NewWBIKeys(imageKey, subKey)
	if err != nil {
		t.Fatalf("NewWBIKeys() error = %v", err)
	}
	mixin, err := sign.MixinKey(keys)
	if err != nil {
		t.Fatalf("MixinKey() error = %v", err)
	}
	if mixin != "OPscVixApSk66dND2LfRBjKt43oHmGJn" {
		t.Fatalf("MixinKey() = %q, want Rust vector", mixin)
	}
}

func TestSignWBIAtMatchesRustVector(t *testing.T) {
	t.Parallel()

	keys, _ := sign.NewWBIKeys(imageKey, subKey)
	signed, err := sign.SignWBIAt(map[string]string{
		"foo": "value!'()*",
		"bar": "space value",
	}, keys, 1_700_000_000)
	if err != nil {
		t.Fatalf("SignWBIAt() error = %v", err)
	}
	if signed["foo"] != "value" || signed["wts"] != "1700000000" {
		t.Fatalf("signed params = %#v, want filtered value and timestamp", signed)
	}
	if signed["w_rid"] != "e0bbf3c23838f46f9b7b2785d19dd01e" {
		t.Fatalf("w_rid = %q, want Rust vector", signed["w_rid"])
	}
	if query := sign.EncodeWBIQuery(map[string]string{"bar": "space value"}); query != "bar=space%20value" {
		t.Fatalf("EncodeWBIQuery() = %q, want RFC 3986 space", query)
	}
}

func TestWBIKeysFromURLs(t *testing.T) {
	t.Parallel()

	keys, err := sign.WBIKeysFromURLs(
		"https://i0.hdslb.com/bfs/wbi/abc123.png",
		"https://i0.hdslb.com/bfs/wbi/def456.png",
	)
	if err != nil {
		t.Fatalf("WBIKeysFromURLs() error = %v", err)
	}
	if keys.Image() != "abc123" || keys.Sub() != "def456" {
		t.Fatalf("keys = %q/%q, want abc123/def456", keys.Image(), keys.Sub())
	}
	if _, err := sign.WBIKeysFromURLs("https://example.com/no-extension", "https://example.com/sub.png"); err == nil {
		t.Fatal("WBIKeysFromURLs(malformed) error = nil, want error")
	}
}
