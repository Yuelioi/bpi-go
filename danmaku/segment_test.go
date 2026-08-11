package danmaku_test

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/danmaku"
)

func TestSegmentParamsMatchPromotedContract(t *testing.T) {
	t.Parallel()

	params, err := danmaku.NewSegmentParams(1, 16546, 1)
	if err != nil {
		t.Fatalf("NewSegmentParams() error = %v", err)
	}
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Encode() != "oid=16546&segment_index=1&type=1" {
		t.Fatalf("query = %q, want promoted contract", query.Encode())
	}
}

func TestSegmentParamsRejectInvalidRequiredValuesAndRange(t *testing.T) {
	t.Parallel()

	for _, call := range []func() error{
		func() error { _, err := danmaku.NewSegmentParams(0, 1, 1); return err },
		func() error { _, err := danmaku.NewSegmentParams(1, 0, 1); return err },
		func() error { _, err := danmaku.NewSegmentParams(1, 1, 0); return err },
		func() error {
			params, _ := danmaku.NewSegmentParams(1, 1, 1)
			_, err := params.WithRange(2, 1)
			return err
		},
	} {
		err := call()
		var parameterErr *bpi.ParameterError
		if !errors.As(err, &parameterErr) {
			t.Fatalf("error = %T %v, want ParameterError", err, err)
		}
	}
}
