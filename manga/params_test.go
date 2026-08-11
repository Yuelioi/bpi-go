package manga

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestCouponsParamsDefaultsAndValidation(t *testing.T) {
	params := NewCouponsParams()
	body, err := params.RequestBody()
	if err != nil {
		t.Fatalf("RequestBody() error = %v", err)
	}
	request := body.(couponsRequest)
	if request.Page != 1 || request.PageSize != 20 || !request.NotExpired || request.TabType != 1 || request.Type != 0 {
		t.Fatalf("body = %+v", request)
	}
	for _, err := range []error{
		func() error { _, err := params.WithPage(0); return err }(),
		func() error { _, err := params.WithPageSize(0); return err }(),
		func() error { _, err := params.WithPageSize(101); return err }(),
	} {
		var parameterError *bpierr.ParameterError
		if !errors.As(err, &parameterError) {
			t.Fatalf("error = %v, want ParameterError", err)
		}
	}
}
