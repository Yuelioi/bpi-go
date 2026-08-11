package bangumi_test

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/bangumi"
)

func TestTimelineParamsMatchPromotedContract(t *testing.T) {
	t.Parallel()

	params, err := bangumi.NewTimelineParams(bangumi.TimelineAnime, 3, 7)
	if err != nil {
		t.Fatalf("NewTimelineParams() error = %v", err)
	}
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if query.Encode() != "after=7&before=3&types=1" {
		t.Fatalf("query = %q, want promoted contract", query.Encode())
	}
}

func TestTimelineParamsRejectInvalidTypeAndOffsets(t *testing.T) {
	t.Parallel()

	for _, call := range []func() error{
		func() error { _, err := bangumi.NewTimelineParams(bangumi.TimelineType(2), 3, 7); return err },
		func() error { _, err := bangumi.NewTimelineParams(bangumi.TimelineAnime, 8, 7); return err },
		func() error { _, err := bangumi.NewTimelineParams(bangumi.TimelineAnime, 3, -1); return err },
	} {
		err := call()
		var parameterErr *bpi.ParameterError
		if !errors.As(err, &parameterErr) {
			t.Fatalf("error = %T %v, want ParameterError", err, err)
		}
	}
}
