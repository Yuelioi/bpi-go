// Package wallet contains private wallet read parameters and response models.
package wallet

import (
	"strings"

	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

type InfoParams struct {
	platformType uint32
	timestampMS  int64
	traceID      int64
	version      string
}

func InfoParamsAtMilliseconds(timestampMS int64) (InfoParams, error) {
	if timestampMS <= 0 {
		return InfoParams{}, &bpierr.ParameterError{Field: "timestamp", Message: "timestamp must be positive"}
	}
	return InfoParams{platformType: 3, timestampMS: timestampMS, traceID: timestampMS, version: "1.0"}, nil
}

func (p InfoParams) WithPlatformType(platformType uint32) (InfoParams, error) {
	if platformType == 0 {
		return InfoParams{}, &bpierr.ParameterError{Field: "platformType", Message: "platform type must be non-zero"}
	}
	p.platformType = platformType
	return p, nil
}

func (p InfoParams) WithTraceID(traceID int64) (InfoParams, error) {
	if traceID <= 0 {
		return InfoParams{}, &bpierr.ParameterError{Field: "traceId", Message: "trace ID must be positive"}
	}
	p.traceID = traceID
	return p, nil
}

func (p InfoParams) WithVersion(version string) (InfoParams, error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return InfoParams{}, &bpierr.ParameterError{Field: "version", Message: "version cannot be blank"}
	}
	p.version = version
	return p, nil
}

type infoRequest struct {
	CSRF         string `json:"csrf"`
	PlatformType uint32 `json:"platformType"`
	Timestamp    int64  `json:"timestamp"`
	TraceID      int64  `json:"traceId"`
	Version      string `json:"version"`
}

func (p InfoParams) RequestBody(csrf string) (any, error) {
	if csrf == "" {
		return nil, &bpierr.ParameterError{Field: "csrf", Message: "CSRF token cannot be empty"}
	}
	if p.platformType == 0 || p.timestampMS <= 0 || p.traceID <= 0 || p.version == "" {
		return nil, &bpierr.ParameterError{Field: "params", Message: "wallet parameters are not initialized"}
	}
	return infoRequest{CSRF: csrf, PlatformType: p.platformType, Timestamp: p.timestampMS, TraceID: p.traceID, Version: p.version}, nil
}

type Info struct {
	MID               int64   `json:"mid"`
	TotalBPCoin       float64 `json:"totalBp"`
	DefaultBPCoin     float64 `json:"defaultBp"`
	IOSBPCoin         float64 `json:"iosBp"`
	CouponBalance     float64 `json:"couponBalance"`
	AvailableBPCoin   float64 `json:"availableBp"`
	UnavailableBPCoin float64 `json:"unavailableBp"`
	UnavailableReason string  `json:"unavailableReason"`
	Tip               string  `json:"tip"`
	ShowClassBalance  int64   `json:"needShowClassBalance"`
}
