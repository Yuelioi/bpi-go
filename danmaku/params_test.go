package danmaku

import (
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestPromotedParamsEncodeQueries(t *testing.T) {
	t.Parallel()
	historyCID, _ := ids.NewCID(772096113)
	cid, _ := ids.NewCID(413195701)
	xmlCID, _ := ids.NewCID(16546)
	bvid, _ := ids.NewBVID("BV1fK4y1t741")
	history, _ := NewHistoryDatesParams(historyCID, "2022-01")
	thumbup, _ := NewThumbupStatsParams(cid, 1932011031958944000)
	view, _ := NewWebViewParams(1, 16546)
	bytesParams, _ := NewHistoryBytesParams(1, 16546, "2022-01-01")
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{"history", func() string { q, _ := history.EncodeQuery(); return q.Encode() }, "month=2022-01&oid=772096113&type=1"},
		{"snapshot", func() string { q, _ := NewSnapshotByBVID(bvid).EncodeQuery(); return q.Encode() }, "aid=BV1fK4y1t741"},
		{"thumbup", func() string { q, _ := thumbup.EncodeQuery(); return q.Encode() }, "ids=1932011031958944000&oid=413195701"},
		{"adv", func() string { q, _ := NewAdvStateParams(cid).EncodeQuery(); return q.Encode() }, "cid=413195701&mode=sp"},
		{"view", func() string { q, _ := view.EncodeQuery(); return q.Encode() }, "oid=16546&type=1"},
		{"history-bytes", func() string { q, _ := bytesParams.EncodeQuery(); return q.Encode() }, "date=2022-01-01&oid=16546&type=1"},
		{"xml", func() string { q, _ := NewXMLListParams(xmlCID).EncodeQuery(); return q.Encode() }, "oid=16546"},
	}
	for _, test := range tests {
		if got := test.got(); got != test.want {
			t.Errorf("%s query = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestDateAndIDListsAreValidated(t *testing.T) {
	t.Parallel()
	if _, err := NewHistoryBytesParams(1, 1, "2022-02-31"); err == nil {
		t.Fatal("invalid calendar date accepted")
	}
	if _, err := NewThumbupStatsParams(ids.CID(1)); err == nil {
		t.Fatal("empty danmaku ID list accepted")
	}
}
