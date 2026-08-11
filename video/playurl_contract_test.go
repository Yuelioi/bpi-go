package video_test

import (
	"context"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	bpi "github.com/Yuelioi/bpi-go"
	core "github.com/Yuelioi/bpi-go/client"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/contracttest"
	"github.com/Yuelioi/bpi-go/internal/sign"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/video"
)

func TestVideoPlayURLSignsPromotedContractAndDecodesFixture(t *testing.T) {
	t.Parallel()

	const navigationBody = `{"code":-101,"data":{"wbi_img":{"img_url":"https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png","sub_url":"https://i0.hdslb.com/bfs/wbi/ABCDEFGHIJKLMNOPQRSTUVWXYZ654321.png"}}}`
	fixed := time.Unix(1_700_000_000, 0)
	fixture := contracttest.Fixture(t, "video", "playurl", "play-url", "responses", "success.json")
	unsigned := map[string]string{
		"bvid": "BV1xx411c7mD", "cid": "62131", "qn": "32", "fnval": "16", "fnver": "0", "platform": "pc",
	}
	keys, _ := sign.NewWBIKeys("abcdefghijklmnopqrstuvwxyz123456", "ABCDEFGHIJKLMNOPQRSTUVWXYZ654321")
	wantSigned, _ := sign.SignWBIAt(unsigned, keys, uint64(fixed.Unix()))

	var navCalls atomic.Int32
	var playCalls atomic.Int32
	client, err := bpi.NewClient(
		core.WithClock(func() time.Time { return fixed }),
		bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			switch request.URL.Path {
			case "/x/web-interface/nav":
				navCalls.Add(1)
				return testutil.JSONResponse(http.StatusOK, navigationBody), nil
			case "/x/player/wbi/playurl":
				playCalls.Add(1)
				assertSignedQuery(t, request.URL.Query(), wantSigned)
				return testutil.JSONResponse(http.StatusOK, string(fixture)), nil
			default:
				t.Fatalf("unexpected request URL %s", request.URL)
				return nil, nil
			}
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62131)
	params := video.PlayURLByBVID(bvid, cid).WithQuality(32).WithFormatFlags(16).WithFormatVersion(0)
	for range 2 {
		payload, err := client.Video().PlayURL(context.Background(), params)
		if err != nil {
			t.Fatalf("PlayURL() error = %v", err)
		}
		if payload.Quality != 32 || payload.Format != "flv480" || payload.DASH == nil || len(payload.DASH.Video) == 0 || payload.DASH.Video[0].BaseURL != "https://example.invalid/bilibili/playurl/redacted.m4s" {
			t.Fatalf("PlayURL() = %+v, want promoted payload", payload)
		}
	}
	if navCalls.Load() != 1 || playCalls.Load() != 2 {
		t.Fatalf("calls nav/play = %d/%d, want 1/2", navCalls.Load(), playCalls.Load())
	}
}

func assertSignedQuery(t *testing.T, got url.Values, want map[string]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("signed query = %v, want %v", got, want)
	}
	for key, value := range want {
		if got.Get(key) != value {
			t.Fatalf("signed query[%q] = %q, want %q", key, got.Get(key), value)
		}
	}
}
