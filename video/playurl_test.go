package video_test

import (
	"encoding/json"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/video"
)

func TestPlayURLParamsMatchPromotedUnsignedQuery(t *testing.T) {
	t.Parallel()

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62131)
	params := video.PlayURLByBVID(bvid, cid).
		WithQuality(32).
		WithFormatFlags(16).
		WithFormatVersion(0)
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Encode() != "bvid=BV1xx411c7mD&cid=62131&fnval=16&fnver=0&platform=pc&qn=32" {
		t.Fatalf("query = %q, want promoted unsigned query", query.Encode())
	}
}

func TestPlayURLParamsUseAvidAndValidateCID(t *testing.T) {
	t.Parallel()

	aid, _ := ids.NewAID(170001)
	query, err := video.PlayURLByAID(aid, ids.CID(1)).EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Get("avid") != "170001" || query.Get("aid") != "" {
		t.Fatalf("query = %v, want avid", query)
	}
	if _, err := video.PlayURLByAID(aid, ids.CID(0)).EncodeQuery(); err == nil {
		t.Fatal("EncodeQuery(zero CID) error = nil")
	}
}

func TestPlayURLAcceptsNegativeResumeSentinels(t *testing.T) {
	t.Parallel()

	var payload video.PlayURL
	if err := json.Unmarshal([]byte(`{"last_play_time":-1000,"last_play_cid":-1000}`), &payload); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if payload.LastPlayTime != -1000 || payload.LastPlayCID != -1000 {
		t.Fatalf("resume sentinels = %d/%d", payload.LastPlayTime, payload.LastPlayCID)
	}
}
