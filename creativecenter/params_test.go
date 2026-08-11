package creativecenter

import (
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestPromotedParamsEncodeQueries(t *testing.T) {
	t.Parallel()
	aid, _ := ids.NewAID(113602455409683)
	season, _ := ids.NewSeasonID(4294056)
	section, _ := ids.NewSeasonID(176088)
	seasonList, _ := NewSeasonListParams(1, 10)
	seasonList = seasonList.WithOrder(SeasonOrderCreated).WithSort(SortDescending)
	archives, _ := NewArchivesListParams(1, 10)
	compare, _ := NewArchiveCompareParams().WithSize(3)
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{"season-list", func() string { q, _ := seasonList.EncodeQuery(); return q.Encode() }, "order=ctime&pn=1&ps=10&sort=desc"},
		{"season-info", func() string { q, _ := NewSeasonInfoParams(season).EncodeQuery(); return q.Encode() }, "id=4294056"},
		{"season-aid", func() string { q, _ := NewSeasonByAIDParams(aid).EncodeQuery(); return q.Encode() }, "id=113602455409683"},
		{"section", func() string { q, _ := NewSectionParams(section).EncodeQuery(); return q.Encode() }, "id=176088"},
		{"archives", func() string { q, _ := archives.EncodeQuery(); return q.Encode() }, "pn=1&ps=10"},
		{"archive-videos", func() string { q, _ := NewArchiveVideosParams(aid).EncodeQuery(); return q.Encode() }, "aid=113602455409683"},
		{"compare", func() string { q, _ := compare.EncodeQuery(); return q.Encode() }, "size=3"},
		{"video-trend", func() string { q, _ := NewVideoTrendParams(VideoTrendPlay).EncodeQuery(); return q.Encode() }, "type=1"},
		{"article-trend", func() string { q, _ := NewArticleTrendParams(ArticleTrendRead).EncodeQuery(); return q.Encode() }, "type=1"},
	}
	for _, test := range tests {
		if got := test.got(); got != test.want {
			t.Errorf("%s query = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestPaginationRejectsZero(t *testing.T) {
	t.Parallel()
	if _, err := NewSeasonListParams(0, 10); err == nil {
		t.Fatal("zero page accepted")
	}
	if _, err := NewArchivesListParams(1, 0); err == nil {
		t.Fatal("zero page size accepted")
	}
}
