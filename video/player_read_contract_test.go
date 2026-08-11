package video_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/video"
)

func TestVideoPlayerReadBatchUsesPromotedContracts(t *testing.T) {
	t.Parallel()

	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		query := request.URL.Query()
		var body []byte
		switch request.URL.Path {
		case "/x/player/online/total":
			if query.Get("bvid") != "BV1xx411c7mD" || query.Get("cid") != "62131" {
				t.Errorf("online query = %v", query)
			}
			body = contracttest.Fixture(t, "video", "player-read", "online-total", "responses", "success.json")
		case "/x/web-interface/archive/related":
			if query.Get("bvid") != "BV1xx411c7mD" {
				t.Errorf("related query = %v", query)
			}
			body = contracttest.Fixture(t, "video", "player-read", "related-videos", "responses", "success.json")
		case "/x/web-interface/view/detail/tag":
			if query.Get("bvid") != "BV1xx411c7mD" || query.Get("cid") != "62131" {
				t.Errorf("tags query = %v", query)
			}
			body = contracttest.Fixture(t, "video", "player-read", "tags", "responses", "success.json")
		case "/x/stein/edgeinfo_v2":
			if query.Get("aid") != "114347430905959" || query.Get("graph_version") != "1273647" {
				t.Errorf("interactive query = %v", query)
			}
			body = contracttest.Fixture(t, "video", "player-read", "interactive-info", "responses", "success.json")
		default:
			t.Fatalf("unexpected request %s %s", request.Method, request.URL)
		}
		if request.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", request.Method)
		}
		return testutil.JSONResponse(http.StatusOK, string(body)), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62_131)
	online, err := client.Video().OnlineTotal(context.Background(), video.OnlineTotalByBVID(bvid, cid))
	if err != nil || online.Total != "8" || !online.ShowSwitch.Total {
		t.Fatalf("OnlineTotal() = %+v, %v", online, err)
	}
	related, err := client.Video().RelatedVideos(context.Background(), video.RelatedByBVID(bvid))
	if err != nil || len(related) != 1 || related[0].BVID == "" {
		t.Fatalf("RelatedVideos() = %+v, %v", related, err)
	}
	tags, err := client.Video().Tags(context.Background(), video.TagsByBVID(bvid).WithCID(cid))
	if err != nil || tags == nil || len(tags) != 0 {
		t.Fatalf("Tags() = %+v, %v; want non-nil empty promoted payload", tags, err)
	}
	aid, _ := ids.NewAID(114_347_430_905_959)
	interactiveParams, _ := video.InteractiveInfoByAID(aid, 1_273_647)
	interactive, err := client.Video().InteractiveVideoInfo(context.Background(), interactiveParams)
	if err != nil || interactive.EdgeID != 1 || interactive.Edges == nil || len(interactive.Edges.Questions) != 1 {
		t.Fatalf("InteractiveVideoInfo() = %+v, %v", interactive, err)
	}
}
