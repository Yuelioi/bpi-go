package danmaku_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	bpi "github.com/Yuelioi/bpi-go"
	core "github.com/Yuelioi/bpi-go/client"
	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go/danmaku"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestDanmakuJSONReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	historyCID, _ := ids.NewCID(772096113)
	cid, _ := ids.NewCID(413195701)
	bvid, _ := ids.NewBVID("BV1fK4y1t741")
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				type endpointSpec struct {
					folder, query string
					authenticated bool
				}
				specs := map[string]endpointSpec{
					"/x/v2/dm/history/index": {"history-dates", "month=2022-01&oid=772096113&type=1", true},
					"/x/v2/dm/ajax":          {"snapshot", "aid=BV1fK4y1t741", false},
					"/x/v2/dm/thumbup/stats": {"thumbup-stats", "ids=1932011031958944000&oid=413195701", false},
					"/x/dm/adv/state":        {"adv-state", "cid=413195701&mode=sp", true},
				}
				spec, ok := specs[request.URL.Path]
				if request.Method != http.MethodGet || !ok || request.URL.Query().Encode() != spec.query {
					t.Fatalf("unexpected request %s %s", request.Method, request.URL)
				}
				fileName := profile + ".success.json"
				if spec.authenticated && profile == "anonymous" {
					fileName = "anonymous.requires_login.json"
				}
				body := contracttest.Fixture(t, "danmaku", "json-read", spec.folder, "responses", fileName)
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			historyParams, _ := danmaku.NewHistoryDatesParams(historyCID, "2022-01")
			thumbParams, _ := danmaku.NewThumbupStatsParams(cid, 1932011031958944000)
			domain, ctx := client.Danmaku(), context.Background()
			history, historyErr := domain.HistoryDates(ctx, historyParams)
			snapshot, snapshotErr := domain.Snapshot(ctx, danmaku.NewSnapshotByBVID(bvid))
			stats, statsErr := domain.ThumbupStats(ctx, thumbParams)
			adv, advErr := domain.AdvState(ctx, danmaku.NewAdvStateParams(cid))
			if snapshotErr != nil || len(snapshot) != 0 || statsErr != nil || stats["1932011031958944000"].ID == "" {
				t.Fatalf("public reads failed: snapshot=%v stats=%v", snapshotErr, statsErr)
			}
			if profile == "anonymous" {
				if !bpi.RequiresLogin(historyErr) || !bpi.RequiresLogin(advErr) {
					t.Fatalf("private errors = history %v, adv %v", historyErr, advErr)
				}
				return
			}
			if historyErr != nil || history != nil || advErr != nil || !adv.Accept {
				t.Fatalf("authenticated reads = history %+v/%v adv %+v/%v", history, historyErr, adv, advErr)
			}
		})
	}
}

func TestDanmakuRawReadsUseAllProfiles(t *testing.T) {
	t.Parallel()
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path == "/x/web-interface/nav" {
					return testutil.JSONResponse(http.StatusOK, contracttest.WBINavigationBody), nil
				}
				type endpointSpec struct {
					batch, folder, query string
					authenticated, wbi   bool
				}
				specs := map[string]endpointSpec{
					"/x/v2/dm/wbi/web/seg.so":     {"non-json-read", "web-seg-wbi", "oid=16546&segment_index=1&type=1", false, true},
					"/x/v2/dm/web/view":           {"non-json-read", "web-view", "oid=16546&type=1", false, false},
					"/x/v2/dm/list/seg.so":        {"non-json-read", "mobile-seg", "oid=16546&segment_index=1&type=1", false, false},
					"/x/v2/dm/web/history/seg.so": {"non-json-read", "web-history-seg", "date=2022-01-01&oid=16546&type=1", true, false},
					"/x/v2/dm/history":            {"history-xml", "", "date=2022-01-01&oid=16546&type=1", true, false},
				}
				spec, ok := specs[request.URL.Path]
				if request.Method != http.MethodGet || !ok {
					t.Fatalf("unexpected request %s %s", request.Method, request.URL)
				}
				if spec.wbi {
					contracttest.AssertWBIFields(t, request.URL.Query(), map[string]string{"oid": "16546", "segment_index": "1", "type": "1"})
				} else if request.URL.Query().Encode() != spec.query {
					t.Fatalf("query = %q, want %q", request.URL.Query().Encode(), spec.query)
				}
				if spec.authenticated && profile == "anonymous" {
					body := contracttest.Fixture(t, "danmaku", spec.batch, spec.folder, "responses", "anonymous.requires_login.json")
					return testutil.JSONResponse(http.StatusOK, string(body)), nil
				}
				fileName := profile + ".success.json"
				fixture := contracttest.Fixture(t, "danmaku", spec.batch, spec.folder, "responses", fileName)
				body, contentType, decodeErr := testutil.DecodeBinaryProbeFixture(fixture)
				if decodeErr != nil {
					t.Fatalf("DecodeBinaryProbeFixture() error = %v", decodeErr)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {contentType}}, Body: io.NopCloser(bytes.NewReader(body))}, nil
			})}), contracttest.ProfileOption(profile), core.WithClock(func() time.Time { return time.Unix(1_700_000_000, 0) }))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			segment, _ := danmaku.NewSegmentParams(1, 16546, 1)
			view, _ := danmaku.NewWebViewParams(1, 16546)
			history, _ := danmaku.NewHistoryBytesParams(1, 16546, "2022-01-01")
			domain, ctx := client.Danmaku(), context.Background()
			wbiBody, wbiErr := domain.WebSegmentWBI(ctx, segment)
			viewBody, viewErr := domain.WebView(ctx, view)
			mobileBody, mobileErr := domain.MobileSegment(ctx, segment)
			historyBody, historyErr := domain.WebHistorySegment(ctx, history)
			xmlBody, xmlErr := domain.HistoryXMLBytes(ctx, history)
			if wbiErr != nil || len(wbiBody) == 0 || viewErr != nil || len(viewBody) == 0 || mobileErr != nil || len(mobileBody) == 0 {
				t.Fatalf("public raw reads failed: wbi=%v view=%v mobile=%v", wbiErr, viewErr, mobileErr)
			}
			if profile == "anonymous" {
				if !bpi.RequiresLogin(historyErr) || !bpi.RequiresLogin(xmlErr) {
					t.Fatalf("history errors = protobuf %v XML %v", historyErr, xmlErr)
				}
				return
			}
			if historyErr != nil || len(historyBody) == 0 || xmlErr != nil || len(xmlBody) == 0 {
				t.Fatalf("history reads failed: protobuf %v XML %v", historyErr, xmlErr)
			}
		})
	}
}

func TestDanmakuXMLReadsParsePromotedFixtures(t *testing.T) {
	t.Parallel()
	cid, _ := ids.NewCID(16546)
	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				folder := ""
				switch request.URL.Path {
				case "/x/v1/dm/list.so":
					folder = "list-so"
				case "/16546.xml":
					folder = "comment-xml"
				default:
					t.Fatalf("unexpected request %s", request.URL)
				}
				fixture := contracttest.Fixture(t, "danmaku", "xml-read", folder, "responses", "success.json")
				body, contentType, decodeErr := testutil.DecodeBinaryProbeFixture(fixture)
				if decodeErr != nil {
					t.Fatalf("DecodeBinaryProbeFixture() error = %v", decodeErr)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {contentType}}, Body: io.NopCloser(bytes.NewReader(body))}, nil
			})}), contracttest.ProfileOption(profile))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params := danmaku.NewXMLListParams(cid)
			apiXML, apiErr := client.Danmaku().XMLListSO(context.Background(), params)
			commentXML, commentErr := client.Danmaku().XMLList(context.Background(), params)
			if apiErr != nil || commentErr != nil || len(apiXML.Comments) != 307 || len(commentXML.Comments) != 307 || apiXML.Comments[0].Meta == nil {
				t.Fatalf("XML reads = api %d/%v comment %d/%v", len(apiXML.Comments), apiErr, len(commentXML.Comments), commentErr)
			}
		})
	}
}
