package videoranking

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedQueries(t *testing.T) {
	t.Parallel()

	popular := NewPopularListParams()
	popular, _ = popular.WithPage(1)
	popular, _ = popular.WithPageSize(2)
	query, err := popular.EncodeQuery()
	if err != nil || query.Get("pn") != "1" || query.Get("ps") != "2" {
		t.Fatalf("popular query = %v, %v", query, err)
	}
	ranking := NewRankingListParams()
	ranking, _ = ranking.WithRegionID(1)
	ranking, _ = ranking.WithType(RankingAll)
	query, err = ranking.EncodeQuery()
	if err != nil || query.Get("rid") != "1" || query.Get("type") != "all" {
		t.Fatalf("ranking query = %v, %v", query, err)
	}
	tag, _ := NewRegionTagDynamicParams(136, 10_026_108)
	tag, _ = tag.WithPage(1)
	tag, _ = tag.WithPageSize(2)
	query, err = tag.EncodeQuery()
	if err != nil || query.Get("rid") != "136" || query.Get("tag_id") != "10026108" || query.Get("pn") != "1" || query.Get("ps") != "2" {
		t.Fatalf("tag query = %v, %v", query, err)
	}
	rank, _ := NewRegionNewListRankParams(231, 2, "20260701", "20260703")
	rank, _ = rank.WithOrder(NewListRankClick)
	rank, _ = rank.WithPage(1)
	query, err = rank.EncodeQuery()
	if err != nil || query.Get("search_type") != "video" || query.Get("view_type") != "hot_rank" || query.Get("cate_id") != "231" || query.Get("pagesize") != "2" || query.Get("order") != "click" || query.Get("page") != "1" {
		t.Fatalf("new-list rank query = %v, %v", query, err)
	}
}

func TestParamsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := NewPopularSeriesParams(0); !errors.As(err, &parameterError) {
		t.Fatalf("zero series error = %v", err)
	}
	if _, err := NewRegionDynamicParams(0); !errors.As(err, &parameterError) {
		t.Fatalf("zero region error = %v", err)
	}
	if _, err := NewRegionNewListRankParams(1, 2, "", "20260703"); !errors.As(err, &parameterError) || parameterError.Field != "time_from" {
		t.Fatalf("blank from error = %v", err)
	}
	params, _ := NewRegionNewListRankParams(1, 2, "20260701", "20260703")
	if _, err := params.WithOrder("random"); !errors.As(err, &parameterError) || parameterError.Field != "order" {
		t.Fatalf("bad order error = %v", err)
	}
}
