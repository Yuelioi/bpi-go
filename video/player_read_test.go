package video

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPlayerReadParams(t *testing.T) {
	t.Parallel()

	bvid, _ := ids.NewBVID("BV1xx411c7mD")
	cid, _ := ids.NewCID(62_131)

	online, err := OnlineTotalByBVID(bvid, cid).EncodeQuery()
	if err != nil || online.Get("bvid") != bvid.String() || online.Get("cid") != "62131" {
		t.Fatalf("online query = %v, %v", online, err)
	}
	related, err := RelatedByBVID(bvid).EncodeQuery()
	if err != nil || related.Get("bvid") != bvid.String() {
		t.Fatalf("related query = %v, %v", related, err)
	}
	tags, err := TagsByBVID(bvid).WithCID(cid).EncodeQuery()
	if err != nil || tags.Get("bvid") != bvid.String() || tags.Get("cid") != "62131" {
		t.Fatalf("tags query = %v, %v", tags, err)
	}

	aid, _ := ids.NewAID(114_347_430_905_959)
	interactiveParams, err := InteractiveInfoByAID(aid, 1_273_647)
	if err != nil {
		t.Fatalf("InteractiveInfoByAID() error = %v", err)
	}
	interactive, err := interactiveParams.EncodeQuery()
	if err != nil || interactive.Get("aid") != aid.String() || interactive.Get("graph_version") != "1273647" {
		t.Fatalf("interactive query = %v, %v", interactive, err)
	}
}

func TestPlayerReadParamsRejectZeroValues(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := (OnlineTotalParams{}).EncodeQuery(); !errors.As(err, &parameterError) {
		t.Fatalf("zero online params error = %v, want ParameterError", err)
	}
	aid, _ := ids.NewAID(1)
	if _, err := InteractiveInfoByAID(aid, 0); !errors.As(err, &parameterError) || parameterError.Field != "graph_version" {
		t.Fatalf("zero graph version error = %v, want graph_version ParameterError", err)
	}
	params, _ := InteractiveInfoByAID(aid, 1)
	if _, err := params.WithEdgeID(0); !errors.As(err, &parameterError) || parameterError.Field != "edge_id" {
		t.Fatalf("zero edge error = %v, want edge_id ParameterError", err)
	}
}
