package login

import (
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestParamsEncodePromotedQueries(t *testing.T) {
	t.Parallel()
	mid, _ := ids.NewMID(1000001)
	poll, _ := NewQRPollParams("sanitized-qrcode-key")
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{"notice", func() string { q, _ := NewNoticeParams(mid).EncodeQuery(); return q.Encode() }, "mid=1000001"},
		{"log", func() string { q, _ := NewLogParams().EncodeQuery(); return q.Encode() }, "jsonp=jsonp&web_location=333.33"},
		{"poll", func() string { q, _ := poll.EncodeQuery(); return q.Encode() }, "qrcode_key=sanitized-qrcode-key"},
	}
	for _, test := range tests {
		if got := test.got(); got != test.want {
			t.Errorf("%s query = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestQRPollRejectsBlankKey(t *testing.T) {
	t.Parallel()
	if _, err := NewQRPollParams("  "); err == nil {
		t.Fatal("NewQRPollParams(blank) error = nil")
	}
}
