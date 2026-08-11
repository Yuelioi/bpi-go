package comment_test

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/comment"
)

func TestHotParamsMatchPromotedContract(t *testing.T) {
	t.Parallel()

	target, err := comment.NewTarget(1, 23199)
	if err != nil {
		t.Fatalf("NewTarget() error = %v", err)
	}
	params, err := comment.NewHotParams(target, 2_554_491_176)
	if err != nil {
		t.Fatalf("NewHotParams() error = %v", err)
	}
	params, _ = params.WithPage(1)
	params, _ = params.WithPageSize(5)
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Encode() != "oid=23199&pn=1&ps=5&root=2554491176&type=1" {
		t.Fatalf("query = %q, want promoted contract", query.Encode())
	}
}

func TestHotParamsRejectInvalidTargetRootAndPage(t *testing.T) {
	t.Parallel()

	if _, err := comment.NewTarget(0, 1); err == nil {
		t.Fatal("NewTarget(type=0) error = nil")
	}
	target, _ := comment.NewTarget(1, 1)
	if _, err := comment.NewHotParams(target, 0); err == nil {
		t.Fatal("NewHotParams(root=0) error = nil")
	}
	params, _ := comment.NewHotParams(target, 1)
	_, err := params.WithPageSize(0)
	var parameterErr *bpi.ParameterError
	if !errors.As(err, &parameterErr) || parameterErr.Field != "ps" {
		t.Fatalf("WithPageSize(0) error = %v, want ps ParameterError", err)
	}
}
