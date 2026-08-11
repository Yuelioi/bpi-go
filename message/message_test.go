package message

import "testing"

func TestParamsEncodeContractDefaults(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{"unread-count", func() string { q, _ := NewUnreadCountParams().EncodeQuery(); return q.Encode() }, "build=0&mobi_app=web"},
		{"reply-feed", func() string { q, _ := NewReplyFeedParams().EncodeQuery(); return q.Encode() }, "build=0&mobi_app=web&platform=web&web_location="},
		{"single-unread", func() string { q, _ := NewSingleUnreadParams().EncodeQuery(); return q.Encode() }, "build=0&mobi_app=web&show_dustbin=0&show_unfollow_list=0&unread_type=0"},
	}
	for _, test := range tests {
		if got := test.got(); got != test.want {
			t.Errorf("%s query = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestParamsRejectInvalidValues(t *testing.T) {
	t.Parallel()
	if _, err := NewReplyFeedParams().WithStartID(0); err == nil {
		t.Fatal("WithStartID(0) error = nil")
	}
	if _, err := NewSingleUnreadParams().WithCustomUnreadType(0); err == nil {
		t.Fatal("WithCustomUnreadType(0) error = nil")
	}
	if _, err := NewUnreadCountParams().WithMobiApp("  "); err == nil {
		t.Fatal("WithMobiApp(blank) error = nil")
	}
}
