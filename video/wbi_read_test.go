package video

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestWBIReadParams(t *testing.T) {
	t.Parallel()

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62_131)
	player, err := PlayerInfoByBVID(bvid, cid).EncodeQuery()
	if err != nil || player.Get("bvid") != bvid.String() || player.Get("cid") != cid.String() {
		t.Fatalf("player query = %v, %v", player, err)
	}

	recommendations, err := NewHomepageRecommendationsParams().EncodeQuery()
	if err != nil || recommendations.Get("fresh_type") != "4" || recommendations.Get("ps") != "12" || recommendations.Get("fresh_idx") != "1" || recommendations.Get("fresh_idx_1h") != "1" || recommendations.Get("brush") != "1" || recommendations.Get("fetch_row") != "1" {
		t.Fatalf("recommendations query = %v, %v", recommendations, err)
	}

	mid, _ := ids.NewMID(928_123)
	summary, err := AISummaryByBVID(bvid, cid, mid).EncodeQuery()
	if err != nil || summary.Get("bvid") != bvid.String() || summary.Get("cid") != cid.String() || summary.Get("up_mid") != mid.String() {
		t.Fatalf("summary query = %v, %v", summary, err)
	}
}

func TestWBIReadParamsValidateOptionsAndZeroValues(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := NewHomepageRecommendationsParams().WithPageSize(31); !errors.As(err, &parameterError) || parameterError.Field != "ps" {
		t.Fatalf("page size error = %v, want ps ParameterError", err)
	}
	if _, err := NewHomepageRecommendationsParams().WithFreshIndex(0); !errors.As(err, &parameterError) || parameterError.Field != "fresh_idx" {
		t.Fatalf("fresh index error = %v, want fresh_idx ParameterError", err)
	}
	if _, err := NewHomepageRecommendationsParams().WithFetchRow(0); !errors.As(err, &parameterError) || parameterError.Field != "fetch_row" {
		t.Fatalf("fetch row error = %v, want fetch_row ParameterError", err)
	}
	if _, err := (PlayerInfoParams{}).EncodeQuery(); !errors.As(err, &parameterError) {
		t.Fatalf("zero player params error = %v, want ParameterError", err)
	}
	if _, err := (AISummaryParams{}).EncodeQuery(); !errors.As(err, &parameterError) {
		t.Fatalf("zero summary params error = %v, want ParameterError", err)
	}
}
