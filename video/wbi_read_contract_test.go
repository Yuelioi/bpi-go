package video_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	bpi "github.com/Yuelioi/bpi-go"
	core "github.com/Yuelioi/bpi-go/client"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/contracttest"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/video"
)

const videoWBITestNavigationBody = `{"code":-101,"data":{"wbi_img":{"img_url":"https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png","sub_url":"https://i0.hdslb.com/bfs/wbi/ABCDEFGHIJKLMNOPQRSTUVWXYZ654321.png"}}}`

func TestVideoWBIReadBatchUsesPromotedContracts(t *testing.T) {
	t.Parallel()

	fixed := time.Unix(1_700_000_000, 0)
	var navCalls atomic.Int32
	client, err := bpi.NewClient(
		core.WithClock(func() time.Time { return fixed }),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			query := request.URL.Query()
			var fixture []byte
			switch request.URL.Path {
			case "/x/web-interface/nav":
				navCalls.Add(1)
				return testutil.JSONResponse(http.StatusOK, videoWBITestNavigationBody), nil
			case "/x/player/wbi/v2":
				assertWBIQuery(t, query, map[string]string{"bvid": "BV1xx411c7mD", "cid": "62131"})
				fixture = readFixture(t, "video", "player-read", "player-info-v2", "responses", "anonymous.success.json")
			case "/x/web-interface/wbi/index/top/feed/rcmd":
				assertWBIQuery(t, query, map[string]string{"fresh_type": "4", "ps": "12", "fresh_idx": "1", "fresh_idx_1h": "1", "brush": "1", "fetch_row": "1"})
				fixture = readFixture(t, "video", "player-read", "homepage-recommendations", "responses", "success.json")
			case "/x/web-interface/view/conclusion/get":
				assertWBIQuery(t, query, map[string]string{"bvid": "BV1xx411c7mD", "cid": "62131", "up_mid": "928123"})
				fixture = readFixture(t, "video", "player-read", "ai-summary", "responses", "success.json")
			default:
				t.Fatalf("unexpected request %s %s", request.Method, request.URL)
			}
			return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62_131)
	player, err := client.Video().PlayerInfoV2(context.Background(), video.PlayerInfoByBVID(bvid, cid))
	if err != nil || player.BVID != bvid || player.CID != cid || player.DMMask == nil {
		t.Fatalf("PlayerInfoV2() = %+v, %v", player, err)
	}
	recommendations, err := client.Video().HomepageRecommendations(context.Background(), video.NewHomepageRecommendationsParams())
	if err != nil || len(recommendations.Items) != 1 || recommendations.Items[0].TrackID != "sanitized-track-id" {
		t.Fatalf("HomepageRecommendations() = %+v, %v", recommendations, err)
	}
	mid, _ := ids.NewMID(928_123)
	summary, err := client.Video().AISummary(context.Background(), video.AISummaryByBVID(bvid, cid, mid))
	if err != nil || summary.Code != -1 || summary.ModelResult == nil {
		t.Fatalf("AISummary() = %+v, %v", summary, err)
	}
	if navCalls.Load() != 1 {
		t.Fatalf("nav calls = %d, want one cached key fetch", navCalls.Load())
	}
}

func TestVideoAISummaryDecodesPromotedAnonymousError(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.Path == "/x/web-interface/nav" {
			return testutil.JSONResponse(http.StatusOK, videoWBITestNavigationBody), nil
		}
		return testutil.JSONResponse(http.StatusOK, string(readFixture(t, "video", "player-read", "ai-summary", "responses", "anonymous.error.json"))), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62_131)
	mid, _ := ids.NewMID(928_123)
	_, err = client.Video().AISummary(context.Background(), video.AISummaryByBVID(bvid, cid, mid))
	var apiError *bpi.APIError
	if !errors.As(err, &apiError) || apiError.Code != -101 || !bpi.RequiresLogin(err) {
		t.Fatalf("AISummary() error = %v, want promoted login APIError", err)
	}
}

func assertWBIQuery(t *testing.T, query url.Values, unsigned map[string]string) {
	t.Helper()
	for key, value := range unsigned {
		if len(query[key]) != 1 || query[key][0] != value {
			t.Fatalf("query[%q] = %v, want %q", key, query[key], value)
		}
	}
	if query.Get("wts") != "1700000000" || len(query.Get("w_rid")) != 32 {
		t.Fatalf("WBI metadata = wts %q w_rid %q", query.Get("wts"), query.Get("w_rid"))
	}
	if len(query) != len(unsigned)+2 {
		t.Fatalf("signed query = %v, want %d fields", query, len(unsigned)+2)
	}
}

func readFixture(t *testing.T, parts ...string) []byte {
	return contracttest.Fixture(t, parts...)
}
