package opus

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestSpaceFeedParams(t *testing.T) {
	t.Parallel()

	mid, _ := ids.NewMID(1_000_001)
	params, err := NewSpaceFeedParams(mid)
	if err != nil {
		t.Fatalf("NewSpaceFeedParams() error = %v", err)
	}
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Get("host_mid") != "1000001" || query.Get("page") != "0" || query.Get("type") != "all" || query.Get("web_location") != "333.1387" {
		t.Fatalf("default query = %v", query)
	}

	params = params.WithPage(2)
	params, _ = params.WithOffset("offset-token")
	params, _ = params.WithKind(SpaceFeedArticle)
	query, err = params.EncodeQuery()
	if err != nil || query.Get("page") != "2" || query.Get("offset") != "offset-token" || query.Get("type") != "article" {
		t.Fatalf("custom query = %v, %v", query, err)
	}
}

func TestSpaceFeedParamsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := NewSpaceFeedParams(ids.MID(0)); !errors.As(err, &parameterError) {
		t.Fatalf("zero MID error = %v, want ParameterError", err)
	}
	mid, _ := ids.NewMID(1)
	params, _ := NewSpaceFeedParams(mid)
	if _, err := params.WithOffset("  "); !errors.As(err, &parameterError) {
		t.Fatalf("blank offset error = %v, want ParameterError", err)
	}
	if _, err := params.WithKind(SpaceFeedKind("unknown")); !errors.As(err, &parameterError) {
		t.Fatalf("unknown kind error = %v, want ParameterError", err)
	}
}
