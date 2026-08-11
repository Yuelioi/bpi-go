package misc

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestShortLinkDefaultForm(t *testing.T) {
	aid, _ := ids.NewAID(10001)
	form, err := NewShortLinkParams(aid).EncodeForm()
	if err != nil {
		t.Fatalf("EncodeForm() error = %v", err)
	}
	if got := form.Encode(); got != "build=6114514&buvid=qwq&oid=10001&platform=unix&share_channel=COPY&share_id=main.ugc-video-detail.0.0.pv&share_mode=4" {
		t.Fatalf("form = %q", got)
	}
}

func TestShortLinkRejectsBlankBuvid(t *testing.T) {
	aid, _ := ids.NewAID(1)
	_, err := NewShortLinkParams(aid).WithBuvid(" ")
	var parameterError *bpierr.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("error = %v, want ParameterError", err)
	}
}
