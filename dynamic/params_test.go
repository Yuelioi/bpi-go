package dynamic

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedParameterQueries(t *testing.T) {
	all, err := NewAllParams().EncodeQuery()
	if err != nil || all.Get("features") != DefaultAllFeatures || all.Get("web_location") != "333.1365" {
		t.Fatalf("all query = %v, %v", all, err)
	}
	check, _ := NewCheckNewParams("0")
	query, err := check.EncodeQuery()
	if err != nil || query.Encode() != "update_baseline=0" {
		t.Fatalf("check query = %v, %v", query, err)
	}
	id, _ := ids.NewDynamicID("1099138163191840776")
	detail, err := NewDetailParams(id).EncodeQuery()
	if err != nil || detail.Get("features") != DefaultDetailFeatures || detail.Get("id") != id.String() {
		t.Fatalf("detail query = %v, %v", detail, err)
	}
	live, _ := NewLiveUsersParams().WithSize(1)
	query, err = live.EncodeQuery()
	if err != nil || query.Encode() != "size=1" {
		t.Fatalf("live query = %v, %v", query, err)
	}
	query, err = NewUpUsersParams().EncodeQuery()
	if err != nil || query.Encode() != "teenagers_mode=0" {
		t.Fatalf("up users query = %v, %v", query, err)
	}
}

func TestParamsRejectBlankAndZeroValues(t *testing.T) {
	_, blankErr := NewCheckNewParams("   ")
	_, sizeErr := NewLiveUsersParams().WithSize(0)
	_, idErr := NewItemParams("").EncodeQuery()
	for _, err := range []error{blankErr, sizeErr, idErr} {
		var parameterError *bpierr.ParameterError
		if !errors.As(err, &parameterError) {
			t.Fatalf("error = %v, want ParameterError", err)
		}
	}
}
