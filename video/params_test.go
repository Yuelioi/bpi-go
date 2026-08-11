package video_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/video"
)

func TestVideoInfoParamsEncodeAIDAndBVID(t *testing.T) {
	t.Parallel()

	aid, _ := ids.NewAID(170001)
	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	tests := []struct {
		name  string
		query func() (url.Values, error)
		key   string
		value string
	}{
		{name: "view aid", query: video.ViewByAID(aid).EncodeQuery, key: "aid", value: "170001"},
		{name: "view bvid", query: video.ViewByBVID(bvid).EncodeQuery, key: "bvid", value: "BV1xx411c7mD"},
		{name: "pagelist bvid", query: video.PageListByBVID(bvid).EncodeQuery, key: "bvid", value: "BV1xx411c7mD"},
		{name: "desc aid", query: video.DescByAID(aid).EncodeQuery, key: "aid", value: "170001"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			query, err := test.query()
			if err != nil {
				t.Fatalf("EncodeQuery() error = %v", err)
			}
			if len(query) != 1 || query.Get(test.key) != test.value {
				t.Fatalf("query = %v, want %s=%s", query, test.key, test.value)
			}
		})
	}
}

func TestVideoDetailParamsMatchPromotedContract(t *testing.T) {
	t.Parallel()

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	query, err := video.DetailByBVID(bvid).WithElectric(false).EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Get("bvid") != "BV1xx411c7mD" || query.Get("need_elec") != "0" {
		t.Fatalf("query = %v, want promoted detail query", query)
	}
}

func TestVideoParamsRejectZeroAndForgedIDs(t *testing.T) {
	t.Parallel()

	queries := []func() (url.Values, error){
		video.ViewByAID(ids.AID(0)).EncodeQuery,
		video.DetailByBVID(ids.BVID("forged")).EncodeQuery,
		video.PageListByAID(ids.AID(0)).EncodeQuery,
		video.DescByBVID(ids.BVID("forged")).EncodeQuery,
	}
	for _, query := range queries {
		_, err := query()
		var parameterErr *bpi.ParameterError
		if !errors.As(err, &parameterErr) {
			t.Fatalf("EncodeQuery() error = %T %v, want ParameterError", err, err)
		}
	}
}
