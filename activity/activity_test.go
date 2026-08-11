package activity_test

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/activity"
	"github.com/Yuelioi/bpi-go/ids"
)

func TestInfoParamsMatchPromotedContract(t *testing.T) {
	t.Parallel()

	params, err := activity.NewInfoParams(4_017_552)
	if err != nil {
		t.Fatalf("NewInfoParams() error = %v", err)
	}
	bvid, _ := ids.NewBVID("BV1mKY4e8ELy")
	params, err = params.WithBVID(bvid)
	if err != nil {
		t.Fatalf("WithBVID() error = %v", err)
	}
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Get("sid") != "4017552" || query.Get("bvid") != "BV1mKY4e8ELy" {
		t.Fatalf("query = %v, want promoted contract values", query)
	}
}

func TestInfoParamsRejectZeroAndInvalidBVID(t *testing.T) {
	t.Parallel()

	if _, err := activity.NewInfoParams(0); err == nil {
		t.Fatal("NewInfoParams(0) error = nil")
	}
	params, _ := activity.NewInfoParams(1)
	if _, err := params.WithBVID(ids.BVID("invalid")); err == nil {
		t.Fatal("WithBVID(invalid) error = nil")
	}
}

func TestListParamsMatchRustDefaultsAndContract(t *testing.T) {
	t.Parallel()

	defaults, err := (activity.ListParams{}).EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery(defaults) error = %v", err)
	}
	if defaults.Encode() != "http=3&mold=0&plat=1%2C3&pn=1&ps=15" {
		t.Fatalf("default query = %q", defaults.Encode())
	}
	params, err := activity.NewListParams().WithPageSize(1)
	if err != nil {
		t.Fatalf("WithPageSize() error = %v", err)
	}
	contract, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery(contract) error = %v", err)
	}
	if contract.Get("plat") != "1,3" || contract.Get("ps") != "1" {
		t.Fatalf("contract query = %v", contract)
	}
	_, err = activity.NewListParams().WithPage(0)
	var parameterErr *bpi.ParameterError
	if !errors.As(err, &parameterErr) || parameterErr.Field != "pn" {
		t.Fatalf("WithPage(0) error = %v, want pn ParameterError", err)
	}
}
