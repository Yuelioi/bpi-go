package client

import (
	"bytes"
	"encoding/json"
)

// Envelope is the common Bilibili JSON response wrapper. Domain methods
// normally return its payload directly.
type Envelope[T any] struct {
	Code    int
	Data    *T
	Message string
	Status  bool
}

// UnmarshalJSON supports the response aliases observed in promoted contracts.
func (e *Envelope[T]) UnmarshalJSON(data []byte) error {
	var raw struct {
		Code    *int            `json:"code"`
		Errno   *int            `json:"errno"`
		Data    json.RawMessage `json:"data"`
		Result  json.RawMessage `json:"result"`
		Message *string         `json:"message"`
		Msg     *string         `json:"msg"`
		ShowMsg *string         `json:"showMsg"`
		Status  bool            `json:"status"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.Code = 0
	if raw.Code != nil {
		e.Code = *raw.Code
	} else if raw.Errno != nil {
		e.Code = *raw.Errno
	}
	e.Message = ""
	if raw.Message != nil {
		e.Message = *raw.Message
	} else if raw.Msg != nil {
		e.Message = *raw.Msg
	} else if raw.ShowMsg != nil {
		e.Message = *raw.ShowMsg
	}
	e.Status = raw.Status
	e.Data = nil

	payload := raw.Data
	if len(payload) == 0 {
		payload = raw.Result
	}
	if len(payload) == 0 || bytes.Equal(bytes.TrimSpace(payload), []byte("null")) {
		return nil
	}
	var value T
	if err := json.Unmarshal(payload, &value); err != nil {
		// Error envelopes sometimes carry partial or differently shaped data.
		// Preserve the semantic API error instead of masking it with a model
		// decode failure. Compatible non-zero payloads remain available through
		// IntoData for the few protocols that intentionally expose them.
		if e.Code != 0 {
			return nil
		}
		return err
	}
	e.Data = &value
	return nil
}

// DecodeEnvelope decodes a response and preserves a private copy of body when
// the model does not match.
func DecodeEnvelope[T any](body []byte) (Envelope[T], error) {
	var envelope Envelope[T]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return Envelope[T]{}, NewResponseDecodeError(err, body)
	}
	return envelope, nil
}

// EnsureSuccess converts a non-zero response code into APIError.
func (e Envelope[T]) EnsureSuccess() error {
	if e.Code == 0 {
		return nil
	}
	return &APIError{Code: e.Code, Message: e.Message}
}

// IntoPayload returns a required successful payload.
func (e Envelope[T]) IntoPayload() (T, error) {
	var zero T
	if err := e.EnsureSuccess(); err != nil {
		return zero, err
	}
	if e.Data == nil {
		return zero, ErrMissingData
	}
	return *e.Data, nil
}

// IntoOptionalPayload returns an optional successful payload.
func (e Envelope[T]) IntoOptionalPayload() (*T, error) {
	if err := e.EnsureSuccess(); err != nil {
		return nil, err
	}
	return e.Data, nil
}

// IntoData returns a required payload without checking the response code. It
// is reserved for interfaces such as anonymous navigation and QR polling that
// carry useful state alongside a non-zero business code.
func (e Envelope[T]) IntoData() (T, error) {
	var zero T
	if e.Data == nil {
		return zero, ErrMissingData
	}
	return *e.Data, nil
}
