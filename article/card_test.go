package article

import (
	"encoding/json"
	"testing"
)

func TestCardItemSelectsKnownAndUnknownVariants(t *testing.T) {
	t.Parallel()
	var video CardItem
	if err := json.Unmarshal([]byte(`{"aid":2,"bvid":"BV1xx411c7mD","cid":3,"owner":{"mid":2},"stat":{"aid":2}}`), &video); err != nil {
		t.Fatalf("decode video: %v", err)
	}
	if video.Kind() != CardVideo || video.Video.AID.Uint64() != 2 {
		t.Fatalf("video = %+v", video)
	}
	var unknown CardItem
	if err := json.Unmarshal([]byte(`{"future_card":{"id":1}}`), &unknown); err != nil {
		t.Fatalf("decode unknown: %v", err)
	}
	if unknown.Kind() != CardUnknown || len(unknown.Raw) == 0 {
		t.Fatalf("unknown = %+v", unknown)
	}
}
