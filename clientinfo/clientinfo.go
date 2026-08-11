// Package clientinfo contains validated parameters and response models for
// Bilibili client-information endpoints.
package clientinfo

import (
	"net/netip"
	"net/url"
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

// IPParams configures clientinfo.ip. Its zero value asks the endpoint to
// inspect the calling client's address.
type IPParams struct {
	address netip.Addr
}

// NewIPParams creates parameters that omit the optional IP address.
func NewIPParams() IPParams { return IPParams{} }

// IPParamsFor creates parameters for an already parsed IPv4 or IPv6 address.
func IPParamsFor(address netip.Addr) (IPParams, error) {
	if !address.IsValid() || address.Zone() != "" {
		return IPParams{}, &bpierr.ParameterError{Field: "ip", Message: "value must be a valid IPv4 or IPv6 address"}
	}
	return IPParams{address: address}, nil
}

// WithIP returns a copy configured for the supplied IPv4 or IPv6 address.
func (p IPParams) WithIP(address string) (IPParams, error) {
	parsed, err := netip.ParseAddr(strings.TrimSpace(address))
	if err != nil || parsed.Zone() != "" {
		return IPParams{}, &bpierr.ParameterError{Field: "ip", Message: "value must be a valid IPv4 or IPv6 address"}
	}
	p.address = parsed
	return p, nil
}

// EncodeQuery returns a new URL query for the configured address.
func (p IPParams) EncodeQuery() (url.Values, error) {
	if !p.address.IsValid() {
		return url.Values{}, nil
	}
	return url.Values{"ip": {p.address.String()}}, nil
}

// IPInfo is the public geolocation information returned for an IP address.
// Pointer fields preserve the API's distinction between an absent value and
// an explicitly empty string.
type IPInfo struct {
	Country  *string `json:"country"`
	Province *string `json:"province"`
	City     *string `json:"city"`
	ISP      *string `json:"isp"`
	Address  *string `json:"addr"`
}
