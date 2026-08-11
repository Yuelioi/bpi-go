package comment

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestReadParamsMatchPromotedQueries(t *testing.T) {
	target, _ := NewTarget(1, 23199)
	list, _ := NewListParams(target).WithPage(1)
	list, _ = list.WithPageSize(5)
	list, _ = list.WithSort(SortByTime)
	list = list.WithoutHot(false)
	query, err := list.EncodeQuery()
	if err != nil || query.Encode() != "nohot=0&oid=23199&pn=1&ps=5&sort=0&type=1" {
		t.Fatalf("list query = %q, %v", query.Encode(), err)
	}
	replies, _ := NewRepliesParams(target, 2_554_491_176)
	replies, _ = replies.WithPage(1)
	replies, _ = replies.WithPageSize(5)
	query, err = replies.EncodeQuery()
	if err != nil || query.Encode() != "oid=23199&pn=1&ps=5&root=2554491176&type=1" {
		t.Fatalf("replies query = %q, %v", query.Encode(), err)
	}
	query, err = NewCountParams(target).EncodeQuery()
	if err != nil || query.Encode() != "oid=23199&type=1" {
		t.Fatalf("count query = %q, %v", query.Encode(), err)
	}
}

func TestListParamsRejectInvalidPageSize(t *testing.T) {
	target, _ := NewTarget(1, 1)
	_, err := NewListParams(target).WithPageSize(21)
	var parameterError *bpierr.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("error = %v, want ParameterError", err)
	}
}
