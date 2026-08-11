package creativecenter_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/creativecenter"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestCreativeCenterReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	aid, _ := ids.NewAID(113602455409683)
	seasonID, _ := ids.NewSeasonID(4294056)
	sectionID, _ := ids.NewSeasonID(176088)
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				type endpointSpec struct {
					parts       []string
					query       string
					normalError string
				}
				specs := map[string]endpointSpec{
					"/x2/creative/web/seasons":             {[]string{"season", "list"}, "order=ctime&pn=1&ps=10&sort=desc", ""},
					"/x2/creative/web/season":              {[]string{"season", "info"}, "id=4294056", "normal.not_found.json"},
					"/x2/creative/web/season/aid":          {[]string{"season", "aid"}, "id=113602455409683", "normal.not_owner.json"},
					"/x2/creative/web/season/section":      {[]string{"season", "section"}, "id=176088", ""},
					"/x2/creative/web/archives/sp":         {[]string{"videos", "archives-list"}, "pn=1&ps=10", ""},
					"/x/web/archive/videos":                {[]string{"videos", "archive-videos"}, "aid=113602455409683", "normal.permission_denied.json"},
					"/x/web/index/stat":                    {[]string{"statistics", "up-stat"}, "", ""},
					"/x/web/data/archive_diagnose/compare": {[]string{"statistics", "archive-compare"}, "size=3", ""},
					"/x/web/data/article":                  {[]string{"statistics", "article-stat"}, "", ""},
					"/x/web/data/pandect":                  {[]string{"statistics", "video-trend"}, "type=1", ""},
					"/x/web/data/article/thirty":           {[]string{"statistics", "article-trend"}, "type=1", ""},
					"/x/web/data/playsource":               {[]string{"statistics", "play-source"}, "", ""},
					"/x/web/data/base":                     {[]string{"statistics", "viewer-data"}, "", ""},
					"/studio/up-rating/v3/rating/info":     {[]string{"railgun-read", "electromagnetic-info"}, "", ""},
				}
				spec, ok := specs[request.URL.Path]
				if request.Method != http.MethodGet || !ok || request.URL.Query().Encode() != spec.query {
					t.Fatalf("unexpected request %s %s", request.Method, request.URL)
				}
				fileName := profile + ".success.json"
				if profile == "anonymous" {
					fileName = "anonymous.requires_login.json"
				}
				if profile == "normal" && spec.normalError != "" {
					fileName = spec.normalError
				}
				parts := append([]string{"creativecenter"}, spec.parts...)
				parts = append(parts, "responses", fileName)
				body := contracttest.Fixture(t, parts...)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			seasonListParams, _ := creativecenter.NewSeasonListParams(1, 10)
			seasonListParams = seasonListParams.WithOrder(creativecenter.SeasonOrderCreated).WithSort(creativecenter.SortDescending)
			archivesParams, _ := creativecenter.NewArchivesListParams(1, 10)
			compareParams, _ := creativecenter.NewArchiveCompareParams().WithSize(3)
			domain, ctx := client.CreativeCenter(), context.Background()
			seasonList, seasonListErr := domain.SeasonList(ctx, seasonListParams)
			seasonInfo, seasonInfoErr := domain.SeasonInfo(ctx, creativecenter.NewSeasonInfoParams(seasonID))
			seasonByAID, seasonByAIDErr := domain.SeasonByAID(ctx, creativecenter.NewSeasonByAIDParams(aid))
			section, sectionErr := domain.SeasonSection(ctx, creativecenter.NewSectionParams(sectionID))
			archives, archivesErr := domain.ArchivesList(ctx, archivesParams)
			archiveVideos, archiveVideosErr := domain.ArchiveVideos(ctx, creativecenter.NewArchiveVideosParams(aid))
			upStat, upStatErr := domain.UpStat(ctx)
			compare, compareErr := domain.ArchiveCompare(ctx, compareParams)
			articleStat, articleStatErr := domain.ArticleStat(ctx)
			videoTrend, videoTrendErr := domain.VideoTrend(ctx, creativecenter.NewVideoTrendParams(creativecenter.VideoTrendPlay))
			articleTrend, articleTrendErr := domain.ArticleTrend(ctx, creativecenter.NewArticleTrendParams(creativecenter.ArticleTrendRead))
			playSource, playSourceErr := domain.PlaySource(ctx)
			viewer, viewerErr := domain.ViewerData(ctx)
			electromagnetic, electromagneticErr := domain.ElectromagneticInfo(ctx)
			allErrors := map[string]error{"season-list": seasonListErr, "season-info": seasonInfoErr, "season-aid": seasonByAIDErr, "section": sectionErr, "archives": archivesErr, "archive-videos": archiveVideosErr, "up-stat": upStatErr, "compare": compareErr, "article-stat": articleStatErr, "video-trend": videoTrendErr, "article-trend": articleTrendErr, "play-source": playSourceErr, "viewer": viewerErr, "electromagnetic": electromagneticErr}
			if profile == "anonymous" {
				for name, callErr := range allErrors {
					if !bpi.RequiresLogin(callErr) {
						t.Fatalf("%s error = %v, want login error", name, callErr)
					}
				}
				return
			}
			if profile == "normal" {
				var apiErr *bpi.APIError
				if !errors.As(seasonInfoErr, &apiErr) || apiErr.Code != -404 {
					t.Fatalf("SeasonInfo error = %v", seasonInfoErr)
				}
				if !errors.As(seasonByAIDErr, &apiErr) || apiErr.Code != 20103 {
					t.Fatalf("SeasonByAID error = %v", seasonByAIDErr)
				}
				if !bpi.IsPermissionError(archiveVideosErr) {
					t.Fatalf("ArchiveVideos error = %v", archiveVideosErr)
				}
			} else if seasonInfoErr != nil || seasonInfo.Season.ID == 0 || seasonByAIDErr != nil || seasonByAID.ID == 0 || archiveVideosErr != nil || len(archiveVideos.Videos) != 1 {
				t.Fatalf("vip-only successes failed: info=%v aid=%v videos=%v", seasonInfoErr, seasonByAIDErr, archiveVideosErr)
			}
			for name, callErr := range allErrors {
				if name == "season-info" || name == "season-aid" || name == "archive-videos" {
					continue
				}
				if callErr != nil {
					t.Fatalf("%s error = %v", name, callErr)
				}
			}
			if seasonList.Total != 1 || len(section.Episodes) != 1 || len(archives.Archives) != 1 || upStat.TotalClick != 0 || len(compare.List) != 1 || articleStat.View != 0 || len(videoTrend) != 1 || playSource != nil || len(viewer.ViewerArea.Fan) != 1 || electromagnetic.State != 0 {
				t.Fatalf("unexpected successful payload shapes")
			}
			if profile == "normal" && articleTrend != nil {
				t.Fatalf("normal article trend = %+v, want nil", articleTrend)
			}
			if profile == "vip" && len(articleTrend) != 1 {
				t.Fatalf("vip article trend = %+v", articleTrend)
			}
		})
	}
}
