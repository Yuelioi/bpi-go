package historytoview

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestListParamsRejectZeroPageSize(t *testing.T) {
	_, err := NewListParams().WithPageSize(0)
	var parameterError *bpierr.ParameterError
	if !errors.As(err, &parameterError) {
		t.Fatalf("error = %v, want ParameterError", err)
	}
}

func TestListParamsEncodeOptionalFilters(t *testing.T) {
	params, _ := NewListParams().WithPageSize(20)
	params = params.WithMax(1001).WithBusiness(BusinessArchive).WithViewAt(1_700_000_000).WithType(ListAll)
	query, err := params.EncodeQuery()
	if err != nil || query.Encode() != "business=archive&max=1001&ps=20&type=all&view_at=1700000000" {
		t.Fatalf("query = %v, %v", query, err)
	}
}
