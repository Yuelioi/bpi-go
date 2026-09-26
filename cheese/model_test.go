package cheese_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yuelioi/bpi-go/cheese"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", name, err)
	}
	return data
}

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

func TestCheeseDRMPlayURLPreservesEncryptionMetadata(t *testing.T) {
	t.Parallel()

	var envelope struct {
		Data cheese.PlayURL `json:"data"`
	}
	if err := json.Unmarshal(readFixture(t, "drm.anonymous.sanitized.json"), &envelope); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	play := envelope.Data
	if len(play.AcceptQualities) != 0 || play.AcceptFormat != "" || len(play.AcceptDescriptions) != 0 {
		t.Fatalf("missing accept_* fields should decode to zero values: %+v", play)
	}
	if play.IsDRM == nil || !*play.IsDRM {
		t.Fatalf("IsDRM = %v, want true", play.IsDRM)
	}
	if play.DRMType == nil || *play.DRMType != "bili_drm" {
		t.Fatalf("DRMType = %v, want bili_drm", play.DRMType)
	}
	if play.DRMTechType == nil || *play.DRMTechType != 3 {
		t.Fatalf("DRMTechType = %v, want 3", play.DRMTechType)
	}
	if play.HLS == nil || len(play.HLS.Video) == 0 || len(play.HLS.Audio) == 0 {
		t.Fatalf("HLS = %+v, want video and audio tracks", play.HLS)
	}
}

func TestCheesePreviewPlayURLPreservesDirectStreamAndPreviewFlag(t *testing.T) {
	t.Parallel()

	var envelope struct {
		Data cheese.PlayURL `json:"data"`
	}
	if err := json.Unmarshal(readFixture(t, "preview.anonymous.sanitized.json"), &envelope); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	play := envelope.Data
	if play.IsPreview == nil || *play.IsPreview != 1 {
		t.Fatalf("IsPreview = %v, want 1", play.IsPreview)
	}
	if play.HasPaid {
		t.Fatal("HasPaid = true, want false")
	}
	if play.DASH != nil || len(play.DURLs) != 1 || play.DURLs[0].URL != "https://example.invalid/preview.mp4" {
		t.Fatalf("preview streams = DASH:%v DURLs:%+v", play.DASH, play.DURLs)
	}
}

func TestCheeseCourseWithoutCouponStillDecodes(t *testing.T) {
	t.Parallel()

	var envelope struct {
		Data cheese.Course `json:"data"`
	}
	if err := json.Unmarshal(readFixture(t, "no-coupon.sanitized.json"), &envelope); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if envelope.Data.SeasonID != 877726892 || len(envelope.Data.Episodes) != 7 {
		t.Fatalf("course = season %d, episodes %d", envelope.Data.SeasonID, len(envelope.Data.Episodes))
	}
	if len(envelope.Data.Coupon) != 0 && string(envelope.Data.Coupon) != "null" {
		t.Fatalf("Coupon = %s, want missing or null", envelope.Data.Coupon)
	}
}
