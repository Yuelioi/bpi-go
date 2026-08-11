package videoranking_test

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/videoranking"
)

func TestVideoRankingDomainUsesAllPromotedContractsAndProfiles(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			var navCalls atomic.Int32
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				query := request.URL.Query()
				if request.URL.Path == "/x/web-interface/nav" {
					navCalls.Add(1)
					return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
				}
				folder := map[string]string{
					"/x/web-interface/popular":             "popular-list",
					"/x/web-interface/popular/precious":    "popular-precious",
					"/x/web-interface/popular/series/list": "popular-series-list",
					"/x/web-interface/popular/series/one":  "popular-series-one",
					"/x/web-interface/ranking/v2":          "ranking-list",
					"/x/web-interface/dynamic/region":      "region-dynamic",
					"/x/web-interface/newlist":             "region-newlist",
					"/x/web-interface/newlist_rank":        "region-newlist-rank",
					"/x/web-interface/dynamic/tag":         "region-tag-dynamic",
				}[request.URL.Path]
				if folder == "" {
					t.Fatalf("unexpected request %s", request.URL)
				}
				if folder == "popular-series-one" {
					contracttest.AssertWBIFields(t, query, map[string]string{"number": "1"})
				} else if query.Get("w_rid") != "" {
					t.Fatalf("unexpected WBI query for %s: %v", folder, query)
				}
				kind := "success"
				if folder == "region-dynamic" {
					kind = "error"
				}
				fixture := contracttest.Fixture(t, "video_ranking", "read", folder, "responses", profile+"."+kind+".json")
				return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			ctx := context.Background()

			popularParams := videoranking.NewPopularListParams()
			popularParams, _ = popularParams.WithPage(1)
			popularParams, _ = popularParams.WithPageSize(2)
			popular, err := client.VideoRanking().PopularList(ctx, popularParams)
			if err != nil || popular.Items == nil {
				t.Fatalf("PopularList() = %+v, %v", popular, err)
			}
			precious, err := client.VideoRanking().Precious(ctx)
			if err != nil || precious.Title == "" || precious.Items == nil {
				t.Fatalf("Precious() = %+v, %v", precious, err)
			}
			seriesList, err := client.VideoRanking().PopularSeriesList(ctx)
			if err != nil || len(seriesList.Items) == 0 {
				t.Fatalf("PopularSeriesList() = %+v, %v", seriesList, err)
			}
			seriesParams, _ := videoranking.NewPopularSeriesParams(1)
			series, err := client.VideoRanking().PopularSeries(ctx, seriesParams)
			if err != nil || series.Config.Number != 1 || series.Items == nil {
				t.Fatalf("PopularSeries() = %+v, %v", series, err)
			}
			rankingParams := videoranking.NewRankingListParams()
			rankingParams, _ = rankingParams.WithRegionID(1)
			rankingParams, _ = rankingParams.WithType(videoranking.RankingAll)
			ranking, err := client.VideoRanking().RankingList(ctx, rankingParams)
			if err != nil || ranking.Items == nil {
				t.Fatalf("RankingList() = %+v, %v", ranking, err)
			}

			dynamicParams, _ := videoranking.NewRegionDynamicParams(21)
			dynamicParams, _ = dynamicParams.WithPage(1)
			dynamicParams, _ = dynamicParams.WithPageSize(2)
			_, err = client.VideoRanking().RegionDynamic(ctx, dynamicParams)
			var apiError *bpi.APIError
			if !errors.As(err, &apiError) || apiError.Code != -404 {
				t.Fatalf("RegionDynamic() error = %v", err)
			}

			newListParams, _ := videoranking.NewRegionNewListParams(231)
			newListParams, _ = newListParams.WithPage(1)
			newListParams, _ = newListParams.WithPageSize(2)
			newListParams, _ = newListParams.WithType(1)
			newList, err := client.VideoRanking().RegionNewList(ctx, newListParams)
			if err != nil || newList.Archives == nil || newList.Page.Size != 2 {
				t.Fatalf("RegionNewList() = %+v, %v", newList, err)
			}
			tagParams, _ := videoranking.NewRegionTagDynamicParams(136, 10_026_108)
			tagParams, _ = tagParams.WithPage(1)
			tagParams, _ = tagParams.WithPageSize(2)
			tag, err := client.VideoRanking().RegionTagDynamic(ctx, tagParams)
			if err != nil || tag.Archives == nil || tag.Page.Size != 2 {
				t.Fatalf("RegionTagDynamic() = %+v, %v", tag, err)
			}
			rankParams, _ := videoranking.NewRegionNewListRankParams(231, 2, "20260701", "20260703")
			rankParams, _ = rankParams.WithOrder(videoranking.NewListRankClick)
			rankParams, _ = rankParams.WithPage(1)
			newRank, err := client.VideoRanking().RegionNewListRank(ctx, rankParams)
			if err != nil || newRank.Result == nil || newRank.Page != 1 {
				t.Fatalf("RegionNewListRank() = %+v, %v", newRank, err)
			}
			if navCalls.Load() != 1 {
				t.Fatalf("nav calls = %d, want one", navCalls.Load())
			}
		})
	}
}
