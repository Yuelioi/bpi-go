package audio

import (
	"errors"
	"net/url"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedQueries(t *testing.T) {
	t.Parallel()
	sid, _ := ids.NewAudioID(13_603)
	streamSID, _ := ids.NewAudioID(15_664)
	page, _ := NewPageParams(1, 2)
	collection, _ := NewCollectionInfoParams(15_967_839)
	rank, _ := NewRankListParams(76)
	tests := []struct {
		name   string
		encode func() (url.Values, error)
		want   string
	}{
		{"song", NewSongParams(sid).EncodeQuery, "sid=13603"},
		{"page", page.EncodeQuery, "pn=1&ps=2"},
		{"collection", collection.EncodeQuery, "sid=15967839"},
		{"stream web", NewStreamURLWebParams(sid).EncodeQuery, "privilege=2&quality=2&sid=13603"},
		{"stream", NewStreamURLParams(streamSID, QualityHigh).EncodeQuery, "mid=2&platform=android&privilege=2&quality=2&songid=15664"},
		{"period", func() (url.Values, error) { return NewRankPeriodParams(RankOriginal).EncodeQuery("token") }, "csrf=token&list_type=2"},
		{"rank", func() (url.Values, error) { return rank.EncodeQuery("token") }, "csrf=token&list_id=76"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			query, err := test.encode()
			if err != nil || query.Encode() != test.want {
				t.Fatalf("EncodeQuery() = %q, %v; want %q", query.Encode(), err, test.want)
			}
		})
	}
}

func TestPromotedParametersRejectInvalidValues(t *testing.T) {
	t.Parallel()
	sid, _ := ids.NewAudioID(13_603)
	tests := []func() error{
		func() error { _, err := NewSongParams(0).EncodeQuery(); return err },
		func() error { _, err := NewPageParams(0, 1); return err },
		func() error { _, err := NewPageParams(1, 0); return err },
		func() error { _, err := NewCollectionInfoParams(0); return err },
		func() error { _, err := NewStreamURLWebParams(sid).WithPrivilege(0); return err },
		func() error { _, err := NewStreamURLParams(sid, Quality(9)).EncodeQuery(); return err },
		func() error { _, err := NewStreamURLParams(sid, QualityHigh).WithPlatform(" "); return err },
		func() error { _, err := NewCustomRankListType(0); return err },
		func() error { _, err := NewRankListParams(0); return err },
	}
	for index, run := range tests {
		var parameterError *bpierr.ParameterError
		if err := run(); !errors.As(err, &parameterError) {
			t.Fatalf("case %d error = %v, want ParameterError", index, err)
		}
	}
}
