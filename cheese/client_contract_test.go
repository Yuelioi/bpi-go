package cheese_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/cheese"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestCheesePromotedReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder := ""
				var expected map[string]string
				switch request.URL.Path {
				case "/pugv/view/web/ep/list":
					folder = "info/ep-list"
					expected = map[string]string{"season_id": "556", "ps": "50", "pn": "1"}
				case "/pugv/view/web/season":
					if request.URL.Query().Get("season_id") != "" {
						folder = "info/season-detail-season"
						expected = map[string]string{"season_id": "556"}
					} else {
						folder = "info/season-detail-episode"
						expected = map[string]string{"ep_id": "20767"}
					}
				case "/pugv/player/web/playurl":
					folder = "playurl"
					expected = map[string]string{"avid": "997984154", "ep_id": "163956", "cid": "1183682680", "fnver": "0", "qn": "32", "fnval": "16"}
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				if request.Method != http.MethodGet || request.URL.Host != "api.bilibili.com" {
					t.Fatalf("request = %s %s", request.Method, request.URL)
				}
				contracttest.AssertQuery(t, request.URL.Query(), expected)
				body := contracttest.Fixture(t, "cheese", folder, "responses", profile+".success.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			seasonID, _ := ids.NewSeasonID(556)
			episodeID, _ := ids.NewEpisodeID(20_767)
			aid, _ := ids.NewAID(997_984_154)
			playEpisodeID, _ := ids.NewEpisodeID(163_956)
			cid, _ := ids.NewCID(1_183_682_680)
			domain := client.Cheese()
			ctx := context.Background()

			listParams, _ := cheese.NewEpisodeListParams(seasonID).WithPageSize(50)
			listParams, _ = listParams.WithPage(1)
			list, err := domain.EpisodeList(ctx, listParams)
			if err != nil || list.Page.Total == 0 || len(list.Items) == 0 {
				t.Fatalf("EpisodeList() = page %+v items %d, %v", list.Page, len(list.Items), err)
			}
			season, err := domain.SeasonDetail(ctx, seasonID)
			if err != nil || season.SeasonID != seasonID.Uint64() || len(season.Episodes) == 0 || season.Title == "" {
				t.Fatalf("SeasonDetail() = id %d episodes %d title %q, %v", season.SeasonID, len(season.Episodes), season.Title, err)
			}
			episode, err := domain.EpisodeDetail(ctx, episodeID)
			if err != nil || episode.SeasonID != seasonID.Uint64() || len(episode.Episodes) == 0 {
				t.Fatalf("EpisodeDetail() = id %d episodes %d, %v", episode.SeasonID, len(episode.Episodes), err)
			}
			playParams := cheese.NewPlayURLParams(aid, playEpisodeID, cid).
				WithQuality(cheese.Quality480P).
				WithFormatFlags(cheese.FormatDASH)
			play, err := domain.PlayURL(ctx, playParams)
			if err != nil || play.Quality != 32 || play.DASH == nil || len(play.DASH.Video) == 0 || play.DASH.Video[0].BaseURL == "" {
				t.Fatalf("PlayURL() = quality %d DASH %+v, %v", play.Quality, play.DASH, err)
			}
		})
	}
}

func TestCheeseRejectsZeroDetailIDsBeforeTransport(t *testing.T) {
	client, err := bpi.NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if _, err := client.Cheese().SeasonDetail(context.Background(), 0); err == nil {
		t.Fatal("SeasonDetail() error = nil")
	}
	if _, err := client.Cheese().EpisodeDetail(context.Background(), 0); err == nil {
		t.Fatal("EpisodeDetail() error = nil")
	}
}
