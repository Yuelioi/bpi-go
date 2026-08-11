package sign_test

import (
	"strings"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/sign"
)

func FuzzSignWBIAt(f *testing.F) {
	f.Add("foo", "bar", uint64(1_700_000_000))
	f.Add("keyword", "space value!'()*", uint64(0))
	f.Add("unicode", "哔哩哔哩", ^uint64(0))

	keys, err := sign.NewWBIKeys("abcdefghijklmnopqrstuvwxyz123456", "ABCDEFGHIJKLMNOPQRSTUVWXYZ654321")
	if err != nil {
		f.Fatalf("NewWBIKeys() error = %v", err)
	}
	f.Fuzz(func(t *testing.T, key, value string, timestamp uint64) {
		params := map[string]string{key: value}
		signed, err := sign.SignWBIAt(params, keys, timestamp)
		if err != nil {
			t.Fatalf("SignWBIAt() error = %v", err)
		}
		if params[key] != value || len(params) != 1 {
			t.Fatal("SignWBIAt() mutated caller parameters")
		}
		if len(signed["w_rid"]) != 32 || signed["wts"] == "" {
			t.Fatalf("signed metadata = %#v, want w_rid and wts", signed)
		}
		if key != "w_rid" && key != "wts" && strings.ContainsAny(signed[key], "!'()*") {
			t.Fatalf("signed value %q retained filtered characters", signed[key])
		}
	})
}
