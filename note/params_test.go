package note

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedQueries(t *testing.T) {
	aid, _ := ids.NewAID(338_677_252)
	query, err := NewIsForbidParams(aid).EncodeQuery()
	if err != nil || query.Encode() != "aid=338677252" {
		t.Fatalf("is-forbid query = %v, %v", query, err)
	}
	query, err = NewPublicArchiveListParams(aid).EncodeQuery()
	if err != nil || query.Encode() != "oid=338677252&oid_type=0&pn=1&ps=10" {
		t.Fatalf("public archive query = %v, %v", query, err)
	}
}

func TestPaginationRejectsZero(t *testing.T) {
	_, err := NewPagination().WithPage(0)
	var parameterError *bpierr.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("error = %v, want ParameterError", err)
	}
}
