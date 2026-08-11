package video_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/testutil"
	"github.com/Yuelioi/bpi-go/video"
)

func TestVideoInfoMethodsMatchPromotedContractsAndFixtures(t *testing.T) {
	t.Parallel()

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	tests := []struct {
		name    string
		path    string
		query   string
		fixture string
		assert  func(t *testing.T, client *bpi.Client)
	}{
		{
			name: "view", path: "/x/web-interface/view", query: "bvid=BV1xx411c7mD", fixture: "view",
			assert: func(t *testing.T, client *bpi.Client) {
				payload, err := client.Video().View(context.Background(), video.ViewByBVID(bvid))
				if err != nil {
					t.Fatalf("View() error = %v", err)
				}
				if payload.BVID != bvid || payload.CID.Uint64() != 62131 || payload.Picture != "https://i0.hdslb.com/bfs/archive/sanitized.jpg" {
					t.Fatalf("View() = %+v, want promoted payload", payload)
				}
			},
		},
		{
			name: "detail", path: "/x/web-interface/view/detail", query: "bvid=BV1xx411c7mD&need_elec=0", fixture: "detail",
			assert: func(t *testing.T, client *bpi.Client) {
				payload, err := client.Video().Detail(context.Background(), video.DetailByBVID(bvid).WithElectric(false))
				if err != nil {
					t.Fatalf("Detail() error = %v", err)
				}
				if payload.View.BVID != bvid || len(payload.Related) != 1 || len(payload.Card) == 0 || len(payload.Reply) == 0 {
					t.Fatalf("Detail() = %+v, want promoted payload", payload)
				}
			},
		},
		{
			name: "pagelist", path: "/x/player/pagelist", query: "bvid=BV1xx411c7mD", fixture: "pagelist",
			assert: func(t *testing.T, client *bpi.Client) {
				payload, err := client.Video().PageList(context.Background(), video.PageListByBVID(bvid))
				if err != nil {
					t.Fatalf("PageList() error = %v", err)
				}
				if len(payload) != 1 || payload[0].CID.Uint64() != 62131 {
					t.Fatalf("PageList() = %+v, want promoted page", payload)
				}
			},
		},
		{
			name: "desc", path: "/x/web-interface/archive/desc", query: "bvid=BV1xx411c7mD", fixture: "desc",
			assert: func(t *testing.T, client *bpi.Client) {
				payload, err := client.Video().Desc(context.Background(), video.DescByBVID(bvid))
				if err != nil {
					t.Fatalf("Desc() error = %v", err)
				}
				if payload != "www" {
					t.Fatalf("Desc() = %q, want www", payload)
				}
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			body := contracttest.Fixture(t, "video", "info-read", test.fixture, "responses", "success.json")
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.Path != test.path || request.URL.RawQuery != test.query {
					t.Errorf("request = %s %s?%s, want GET %s?%s", request.Method, request.URL.Path, request.URL.RawQuery, test.path, test.query)
				}
				return testutil.JSONResponse(http.StatusOK, string(body)), nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			test.assert(t, client)
		})
	}
}

func TestVideoParamsFailBeforeNetwork(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return testutil.JSONResponse(http.StatusOK, `{"code":0}`), nil
	})}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if _, err := client.Video().View(context.Background(), video.ViewByBVID(ids.BVID("forged"))); err == nil {
		t.Fatal("View(forged) error = nil")
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid params performed %d network calls", calls.Load())
	}
}
