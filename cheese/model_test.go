package cheese_test

import (
	"encoding/json"
	"testing"

	"github.com/Yuelioi/bpi-go/cheese"
)

func TestCheesePlayURLUsesCurrentQualityDURL(t *testing.T) {
	t.Parallel()

	var play cheese.PlayURL
	if err := json.Unmarshal([]byte(`{
		"durl":[{"order":1,"length":1000,"size":100,"url":"https://example.invalid/current.mp4"}],
		"durls":[{"quality":32,"durl":[{"size":999}]}]
	}`), &play); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(play.DURLs) != 1 || play.DURLs[0].Size != 100 {
		t.Fatalf("DURLs = %+v, want current-quality durl", play.DURLs)
	}
}
