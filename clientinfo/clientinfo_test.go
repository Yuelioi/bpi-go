package clientinfo

import (
	"errors"
	"net/netip"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestIPParams(t *testing.T) {
	t.Parallel()

	query, err := NewIPParams().EncodeQuery()
	if err != nil || len(query) != 0 {
		t.Fatalf("default query = %v, %v; want empty", query, err)
	}

	params, err := IPParamsFor(netip.MustParseAddr("2001:4860:4860::8888"))
	if err != nil {
		t.Fatalf("IPParamsFor() error = %v", err)
	}
	query, err = params.EncodeQuery()
	if err != nil || query.Get("ip") != "2001:4860:4860::8888" {
		t.Fatalf("IPv6 query = %v, %v", query, err)
	}

	params, err = NewIPParams().WithIP(" 8.8.8.8 ")
	if err != nil {
		t.Fatalf("WithIP() error = %v", err)
	}
	query, _ = params.EncodeQuery()
	if query.Get("ip") != "8.8.8.8" {
		t.Fatalf("IPv4 query = %v", query)
	}
}

func TestIPParamsRejectInvalidAddress(t *testing.T) {
	t.Parallel()

	_, err := NewIPParams().WithIP("not-an-ip")
	var parameterError *bpierr.ParameterError
	if !errors.As(err, &parameterError) || parameterError.Field != "ip" {
		t.Fatalf("WithIP() error = %v, want ip ParameterError", err)
	}
	if _, err := IPParamsFor(netip.Addr{}); !errors.As(err, &parameterError) {
		t.Fatalf("IPParamsFor() error = %v, want ParameterError", err)
	}
	if _, err := NewIPParams().WithIP("fe80::1%eth0"); !errors.As(err, &parameterError) {
		t.Fatalf("scoped IPv6 error = %v, want ParameterError", err)
	}
}
