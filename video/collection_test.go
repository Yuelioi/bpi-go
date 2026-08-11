package video

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestCollectionParamsMatchPromotedQueries(t *testing.T) {
	t.Parallel()

	mid, _ := ids.NewMID(1_000_001)
	seasonID, _ := ids.NewSeasonID(4_294_056)
	seasons := NewSeasonsArchivesParams(mid, seasonID).WithSortReverse(false)
	query, err := seasons.EncodeQuery()
	if err != nil || query.Get("mid") != "1000001" || query.Get("season_id") != "4294056" || query.Get("page_num") != "1" || query.Get("page_size") != "20" || query.Get("sort_reverse") != "false" {
		t.Fatalf("seasons query = %v, %v", query, err)
	}
	home, err := NewHomeSeasonsSeriesParams(mid).EncodeQuery()
	if err != nil || home.Get("page_num") != "1" || home.Get("page_size") != "10" {
		t.Fatalf("home query = %v, %v", home, err)
	}
	listParams := NewSeasonsSeriesParams(mid)
	listParams, _ = listParams.WithPage(1)
	listParams, _ = listParams.WithPageSize(5)
	list, err := listParams.EncodeQuery()
	if err != nil || list.Get("page_num") != "1" || list.Get("page_size") != "5" {
		t.Fatalf("list query = %v, %v", list, err)
	}
	seriesID, _ := NewSeriesID(250_285)
	series, err := NewSeriesInfoParams(seriesID).EncodeQuery()
	if err != nil || series.Get("series_id") != "250285" {
		t.Fatalf("series query = %v, %v", series, err)
	}
	archivesParams := NewSeriesArchivesParams(mid, seriesID)
	archivesParams, _ = archivesParams.WithSort(CollectionArchiveAscending)
	archivesParams, _ = archivesParams.WithPage(1)
	archivesParams, _ = archivesParams.WithPageSize(10)
	archives, err := archivesParams.EncodeQuery()
	if err != nil || archives.Get("sort") != "asc" || archives.Get("pn") != "1" || archives.Get("ps") != "10" {
		t.Fatalf("archives query = %v, %v", archives, err)
	}
}

func TestCollectionParamsAndPageAliasesValidate(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := NewSeriesID(0); !errors.As(err, &parameterError) {
		t.Fatalf("NewSeriesID(0) error = %v, want ParameterError", err)
	}
	mid, _ := ids.NewMID(1)
	seriesID, _ := NewSeriesID(1)
	if _, err := NewSeriesArchivesParams(mid, seriesID).WithSort("random"); !errors.As(err, &parameterError) {
		t.Fatalf("WithSort(random) error = %v, want ParameterError", err)
	}
	if _, err := NewSeasonsSeriesParams(mid).WithPage(0); !errors.As(err, &parameterError) {
		t.Fatalf("WithPage(0) error = %v, want ParameterError", err)
	}

	var page CollectionPage
	if err := json.Unmarshal([]byte(`{"num":2,"size":10,"total":25}`), &page); err != nil {
		t.Fatalf("UnmarshalJSON(alias) error = %v", err)
	}
	if page.Page != 2 || page.Size != 10 || page.Total != 25 {
		t.Fatalf("alias page = %+v", page)
	}
	if err := json.Unmarshal([]byte(`{"page_num":3,"page_size":5,"total":11}`), &page); err != nil {
		t.Fatalf("UnmarshalJSON(canonical) error = %v", err)
	}
	if page.Page != 3 || page.Size != 5 || page.Total != 11 {
		t.Fatalf("canonical page = %+v", page)
	}
}
