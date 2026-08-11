package danmaku_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/contracttest"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/danmaku"
	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestDanmakuWebSegmentReturnsPromotedRawBinaryBodies(t *testing.T) {
	t.Parallel()

	for _, profile := range []string{"anonymous", "normal", "vip"} {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			t.Parallel()
			fixture := contracttest.Fixture(t, "danmaku", "non-json-read", "web-seg", "responses", profile+".success.json")
			want, contentType, err := testutil.DecodeBinaryProbeFixture(fixture)
			if err != nil {
				t.Fatalf("DecodeBinaryProbeFixture() error = %v", err)
			}
			client, err := bpi.NewClient(bpi.WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
				if request.Method != http.MethodGet || request.URL.Path != "/x/v2/dm/web/seg.so" || request.URL.Query().Encode() != "oid=16546&segment_index=1&type=1" {
					t.Errorf("request = %s %s", request.Method, request.URL)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {contentType}},
					Body:       io.NopCloser(bytes.NewReader(want)),
				}, nil
			})}))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}
			params, _ := danmaku.NewSegmentParams(1, 16546, 1)
			got, err := client.Danmaku().WebSegment(context.Background(), params)
			if err != nil {
				t.Fatalf("Danmaku().WebSegment() error = %v", err)
			}
			if !bytes.Equal(got, want) || len(got) == 0 || bytes.HasPrefix(got, []byte("{")) {
				t.Fatalf("WebSegment() returned %d bytes, want raw %d-byte protobuf", len(got), len(want))
			}
		})
	}
}
