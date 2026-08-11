package video_test

import (
	"context"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/video"
)

func TestVideoCollectionDomainUsesAllPromotedContracts(t *testing.T) {
	t.Parallel()

	var navCalls atomic.Int32
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		var fixture []byte
		switch request.URL.Path {
		case "/x/web-interface/nav":
			navCalls.Add(1)
			return testutil.JSONResponse(http.StatusOK, videoWBITestNavigationBody), nil
		case "/x/polymer/web-space/seasons_archives_list":
			assertWBIQueryFields(t, query, map[string]string{"mid": "1000001", "season_id": "4294056", "page_num": "1", "page_size": "20", "sort_reverse": "false"})
			fixture = readFixture(t, "video", "collection-read", "seasons-archives-list", "responses", "success.json")
		case "/x/polymer/web-space/home/seasons_series":
			assertWBIQueryFields(t, query, map[string]string{"mid": "1000001", "page_num": "1", "page_size": "10"})
			fixture = readFixture(t, "video", "collection-read", "home-seasons-series", "responses", "success.json")
		case "/x/polymer/web-space/seasons_series_list":
			assertWBIQueryFields(t, query, map[string]string{"mid": "1000001", "page_num": "1", "page_size": "5"})
			fixture = readFixture(t, "video", "collection-read", "seasons-series-list", "responses", "success.json")
		case "/x/series/series":
			if query.Get("series_id") != "250285" || query.Get("w_rid") != "" {
				t.Fatalf("series-info query = %v", query)
			}
			fixture = readFixture(t, "video", "collection-read", "series-info", "responses", "success.json")
		case "/x/series/archives":
			if query.Get("mid") != "1000001" || query.Get("series_id") != "250285" || query.Get("sort") != "asc" || query.Get("pn") != "1" || query.Get("ps") != "10" || query.Get("w_rid") != "" {
				t.Fatalf("series-archives query = %v", query)
			}
			fixture = readFixture(t, "video", "collection-read", "series-archives", "responses", "success.json")
		default:
			t.Fatalf("unexpected request %s %s", request.Method, request.URL)
		}
		return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	mid, _ := ids.NewMID(1_000_001)
	seasonID, _ := ids.NewSeasonID(4_294_056)
	seasons, err := client.Video().SeasonsArchives(context.Background(), video.NewSeasonsArchivesParams(mid, seasonID).WithSortReverse(false))
	if err != nil || seasons.Meta.SeasonID != seasonID || len(seasons.Archives) != 4 {
		t.Fatalf("SeasonsArchives() = %+v, %v", seasons, err)
	}
	home, err := client.Video().HomeSeasonsSeries(context.Background(), video.NewHomeSeasonsSeriesParams(mid))
	if err != nil || home.Items.Page.Page != 1 || len(home.Items.Series) != 1 {
		t.Fatalf("HomeSeasonsSeries() = %+v, %v", home, err)
	}
	listParams := video.NewSeasonsSeriesParams(mid)
	listParams, _ = listParams.WithPage(1)
	listParams, _ = listParams.WithPageSize(5)
	list, err := client.Video().SeasonsSeries(context.Background(), listParams)
	if err != nil || list.Items.Page.Size != 5 {
		t.Fatalf("SeasonsSeries() = %+v, %v", list, err)
	}
	seriesID, _ := video.NewSeriesID(250_285)
	series, err := client.Video().SeriesInfo(context.Background(), video.NewSeriesInfoParams(seriesID))
	if err != nil || series.Meta.SeriesID != seriesID {
		t.Fatalf("SeriesInfo() = %+v, %v", series, err)
	}
	archivesParams := video.NewSeriesArchivesParams(mid, seriesID)
	archivesParams, _ = archivesParams.WithSort(video.CollectionArchiveAscending)
	archivesParams, _ = archivesParams.WithPage(1)
	archivesParams, _ = archivesParams.WithPageSize(10)
	archives, err := client.Video().SeriesArchives(context.Background(), archivesParams)
	if err != nil || archives.Page.Page != 1 || archives.Page.Size != 10 || len(archives.Archives) != 2 {
		t.Fatalf("SeriesArchives() = %+v, %v", archives, err)
	}
	if navCalls.Load() != 1 {
		t.Fatalf("nav calls = %d, want one", navCalls.Load())
	}
}

func assertWBIQueryFields(t *testing.T, query url.Values, fields map[string]string) {
	t.Helper()
	for key, value := range fields {
		if query.Get(key) != value {
			t.Fatalf("query[%q] = %q, want %q", key, query.Get(key), value)
		}
	}
	if query.Get("wts") == "" || len(query.Get("w_rid")) != 32 || len(query) != len(fields)+2 {
		t.Fatalf("signed query = %v", query)
	}
}
