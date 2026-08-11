package webwidget

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestParams(t *testing.T) {
	t.Parallel()

	headerQuery, err := NewHeaderPageParams().EncodeQuery()
	if err != nil || headerQuery.Get("resource_id") != "142" {
		t.Fatalf("header query = %v, %v", headerQuery, err)
	}
	header, err := NewHeaderPageParams().WithResourceID(143)
	if err != nil {
		t.Fatalf("WithResourceID() error = %v", err)
	}
	headerQuery, _ = header.EncodeQuery()
	if headerQuery.Get("resource_id") != "143" {
		t.Fatalf("custom header query = %v", headerQuery)
	}

	region, err := NewRegionBannerParams(1005)
	if err != nil {
		t.Fatalf("NewRegionBannerParams() error = %v", err)
	}
	regionQuery, err := region.EncodeQuery()
	if err != nil || regionQuery.Get("region_id") != "1005" {
		t.Fatalf("region query = %v, %v", regionQuery, err)
	}
}

func TestParamsRejectZeroIDs(t *testing.T) {
	t.Parallel()

	var parameterError *bpierr.ParameterError
	if _, err := NewHeaderPageParams().WithResourceID(0); !errors.As(err, &parameterError) {
		t.Fatalf("WithResourceID(0) error = %v, want ParameterError", err)
	}
	if _, err := NewRegionBannerParams(0); !errors.As(err, &parameterError) {
		t.Fatalf("NewRegionBannerParams(0) error = %v, want ParameterError", err)
	}
	if _, err := (RegionBannerParams{}).EncodeQuery(); !errors.As(err, &parameterError) {
		t.Fatalf("zero RegionBannerParams error = %v, want ParameterError", err)
	}
}

func TestHeaderDataRejectsMalformedSplitLayer(t *testing.T) {
	t.Parallel()

	var header HeaderData
	err := json.Unmarshal([]byte(`{"split_layer":"not-json"}`), &header)
	if err == nil {
		t.Fatal("UnmarshalJSON() error = nil, want malformed split_layer error")
	}
}
