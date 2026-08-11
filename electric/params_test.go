package electric

import (
	"testing"
	"time"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestParamsEncodePromotedQueries(t *testing.T) {
	t.Parallel()
	up, _ := ids.NewMID(1265680561)
	publicUP, _ := ids.NewMID(53456)
	bvid, _ := ids.NewBVID("BV1Dh411S7sS")
	recharge, _ := NewRechargeListParams(1, 10)
	page, _ := NewPaginationParams(1, 10)
	record, _ := NewChargeRecordParams(1, 1)
	rank, _ := NewMemberRankParams(up, 1, 10)
	remarks, _ := NewRemarkListParams(1, 10)
	detail, _ := NewRemarkDetailParams(1)
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{"month", func() string { q, _ := NewMonthUpListParams(publicUP).EncodeQuery(); return q.Encode() }, "up_mid=53456"},
		{"video", func() string { q, _ := NewVideoShowParams(publicUP).WithBVID(bvid).EncodeQuery(); return q.Encode() }, "bvid=BV1Dh411S7sS&mid=53456"},
		{"recharge", func() string { q, _ := recharge.EncodeQuery(); return q.Encode() }, "currentPage=1&customerId=10026&pageSize=10"},
		{"rank-recent", func() string { q, _ := page.EncodeQuery(); return q.Encode() }, "pn=1&ps=10"},
		{"charge-record", func() string { q, _ := record.EncodeQuery(); return q.Encode() }, "page=1&type=1"},
		{"up-mid", func() string { q, _ := NewUpMIDParams(up).EncodeQuery(); return q.Encode() }, "up_mid=1265680561"},
		{"member-rank", func() string { q, _ := rank.EncodeQuery(); return q.Encode() }, "pn=1&ps=10&up_mid=1265680561"},
		{"remarks", func() string { q, _ := remarks.EncodeQuery(); return q.Encode() }, "pn=1&ps=10"},
		{"remark-detail", func() string { q, _ := detail.EncodeQuery(); return q.Encode() }, "id=1"},
	}
	for _, test := range tests {
		if got := test.got(); got != test.want {
			t.Errorf("%s query = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestDateRangesRejectReverseOrder(t *testing.T) {
	t.Parallel()
	params, _ := NewRechargeListParams(1, 10)
	if _, err := params.WithDateRange(time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)); err == nil {
		t.Fatal("WithDateRange(reverse) error = nil")
	}
}
