package sign_test

import (
	"testing"

	"github.com/Yuelioi/bpi-go/internal/sign"
)

func TestBiliTicketSigningMatchesRustVector(t *testing.T) {
	t.Parallel()

	want := "a7da9d971f117aa2b439c4b6cc46c7afbba8ade9f3ca959578af1bcfb37ebd2f"
	digest, err := sign.HexSign("XgwSnGZ1p", 1_234_567_890)
	if err != nil {
		t.Fatalf("HexSign() error = %v", err)
	}
	if digest != want || sign.TicketHexSign(1_234_567_890) != want {
		t.Fatalf("ticket digest = %q, want %q", digest, want)
	}
	params := sign.TicketRequestParams(1_234_567_890, "csrf-token")
	if params["key_id"] != "ec02" || params["context[ts]"] != "1234567890" || params["csrf"] != "csrf-token" || params["hexsign"] != want {
		t.Fatalf("TicketRequestParams() = %#v, want Rust fields", params)
	}
}

func TestHexSignRejectsEmptyKey(t *testing.T) {
	t.Parallel()

	if _, err := sign.HexSign("", 1); err == nil {
		t.Fatal("HexSign(empty) error = nil, want error")
	}
}
