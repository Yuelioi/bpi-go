package bangumi_test

import (
	"context"
	"net/http"
	"testing"

	bpi "github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/bangumi"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestBangumiRemainingPromotedReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder := ""
				var expected map[string]string
				switch request.URL.Path {
				case "/pgc/review/user":
					folder = "info/review-user"
					expected = map[string]string{"media_id": "28220978"}
				case "/pgc/view/web/season":
					if request.URL.Query().Get("season_id") != "" {
						folder = "info/season-detail-season"
						expected = map[string]string{"season_id": "1172"}
					} else {
						folder = "info/season-detail-episode"
						expected = map[string]string{"ep_id": "21265"}
					}
				case "/pgc/web/season/section":
					folder = "info/season-section"
					expected = map[string]string{"season_id": "1172"}
				case "/pgc/player/web/playurl":
					folder = "playurl"
					expected = map[string]string{"fnver": "0", "ep_id": "21265", "qn": "32", "fnval": "16"}
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				if request.Method != http.MethodGet || request.URL.Host != "api.bilibili.com" {
					t.Fatalf("request = %s %s", request.Method, request.URL)
				}
				contracttest.AssertQuery(t, request.URL.Query(), expected)
				body := contracttest.Fixture(t, "bangumi", folder, "responses", profile+".success.json")
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			mediaID, _ := ids.NewMediaID(28_220_978)
			seasonID, _ := ids.NewSeasonID(1_172)
			episodeID, _ := ids.NewEpisodeID(21_265)
			ctx := context.Background()
			domain := client.Bangumi()

			info, err := domain.Review(ctx, bangumi.NewInfoParams(mediaID))
			if err != nil || info.Media.MediaID != mediaID.Uint64() || info.Media.Title == "" {
				t.Fatalf("Review() = %+v, %v", info, err)
			}
			season, err := domain.SeasonDetail(ctx, seasonID)
			if err != nil || season.SeasonID != seasonID.Uint64() || len(season.Episodes) == 0 {
				t.Fatalf("SeasonDetail() = id %d episodes %d, %v", season.SeasonID, len(season.Episodes), err)
			}
			episode, err := domain.EpisodeDetail(ctx, episodeID)
			if err != nil || episode.SeasonID != seasonID.Uint64() || len(episode.Episodes) == 0 {
				t.Fatalf("EpisodeDetail() = id %d episodes %d, %v", episode.SeasonID, len(episode.Episodes), err)
			}
			sections, err := domain.Sections(ctx, bangumi.NewSectionsParams(seasonID))
			if err != nil || sections.Main.ID == 0 || len(sections.Main.Episodes) == 0 {
				t.Fatalf("Sections() = %+v, %v", sections.Main, err)
			}
			playParams := bangumi.PlayURLByEpisodeID(episodeID).WithQuality(bangumi.Quality480P).WithFormatFlags(bangumi.FormatDASH)
			play, err := domain.PlayURL(ctx, playParams)
			if err != nil || play.Quality != 32 || play.DASH == nil || len(play.DASH.Video) == 0 || play.DASH.Video[0].BaseURL == "" {
				t.Fatalf("PlayURL() = quality %d DASH %+v, %v", play.Quality, play.DASH, err)
			}
		})
	}
}
