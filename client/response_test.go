package client_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Yuelioi/bpi-go"
)

type fixturePayload struct {
	Title string `json:"title"`
	AID   uint64 `json:"aid"`
}

func TestEnvelopeExtractsDataAndResultPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fixture   string
		wantTitle string
		wantAID   uint64
	}{
		{name: "data", fixture: "success.json", wantTitle: "fixture video", wantAID: 170001},
		{name: "result alias", fixture: "result-alias.json", wantTitle: "fixture result", wantAID: 170002},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			envelope, err := bpi.DecodeEnvelope[fixturePayload](fixture(t, test.fixture))
			if err != nil {
				t.Fatalf("DecodeEnvelope() error = %v", err)
			}
			payload, err := envelope.IntoPayload()
			if err != nil {
				t.Fatalf("IntoPayload() error = %v", err)
			}
			if payload.Title != test.wantTitle || payload.AID != test.wantAID {
				t.Fatalf("payload = %+v, want title %q and aid %d", payload, test.wantTitle, test.wantAID)
			}
		})
	}
}

func TestEnvelopeReturnsSemanticAPIErrors(t *testing.T) {
	t.Parallel()

	envelope, err := bpi.DecodeEnvelope[fixturePayload](fixture(t, "api-error.json"))
	if err != nil {
		t.Fatalf("DecodeEnvelope() error = %v", err)
	}
	_, err = envelope.IntoPayload()
	if !bpi.RequiresLogin(err) {
		t.Fatalf("RequiresLogin(%v) = false, want true", err)
	}
	var apiErr *bpi.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != -101 {
		t.Fatalf("error = %v, want APIError code -101", err)
	}
}

func TestEnvelopeDistinguishesMissingAndOptionalPayload(t *testing.T) {
	t.Parallel()

	missing, err := bpi.DecodeEnvelope[fixturePayload](fixture(t, "missing-data.json"))
	if err != nil {
		t.Fatalf("DecodeEnvelope(missing) error = %v", err)
	}
	if _, err := missing.IntoPayload(); !errors.Is(err, bpi.ErrMissingData) {
		t.Fatalf("IntoPayload() error = %v, want ErrMissingData", err)
	}

	optional, err := bpi.DecodeEnvelope[fixturePayload](fixture(t, "no-payload.json"))
	if err != nil {
		t.Fatalf("DecodeEnvelope(optional) error = %v", err)
	}
	payload, err := optional.IntoOptionalPayload()
	if err != nil {
		t.Fatalf("IntoOptionalPayload() error = %v", err)
	}
	if payload != nil {
		t.Fatalf("IntoOptionalPayload() = %+v, want nil", payload)
	}
}

func TestEnvelopeHandlesNullMessageAndErrnoAliases(t *testing.T) {
	t.Parallel()

	envelope, err := bpi.DecodeEnvelope[fixturePayload]([]byte(`{"code":0,"message":null,"data":{"title":"ok","aid":1}}`))
	if err != nil {
		t.Fatalf("DecodeEnvelope(null message) error = %v", err)
	}
	if envelope.Message != "" {
		t.Fatalf("Message = %q, want empty", envelope.Message)
	}

	errno, err := bpi.DecodeEnvelope[fixturePayload]([]byte(`{"errno":800501007,"msg":"user not login"}`))
	if err != nil {
		t.Fatalf("DecodeEnvelope(errno) error = %v", err)
	}
	_, err = errno.IntoPayload()
	if !bpi.RequiresLogin(err) {
		t.Fatalf("RequiresLogin(%v) = false, want true", err)
	}
}

func TestResponseDecodeErrorRetainsButDoesNotFormatBody(t *testing.T) {
	t.Parallel()

	body := []byte(`{"code":0,"data":{"aid":"private-response-marker"}}`)
	_, err := bpi.DecodeEnvelope[fixturePayload](body)
	var decodeErr *bpi.ResponseDecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("error = %T %v, want ResponseDecodeError", err, err)
	}
	if string(decodeErr.Body()) != string(body) {
		t.Fatalf("Body() = %q, want original body", decodeErr.Body())
	}
	formatted := fmt.Sprintf("%v %+v", err, err)
	serialized, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("marshal error: %v", marshalErr)
	}
	if strings.Contains(formatted, "private-response-marker") || strings.Contains(string(serialized), "private-response-marker") {
		t.Fatalf("decode error leaked response body: %s %s", formatted, serialized)
	}

	copyOfBody := decodeErr.Body()
	copyOfBody[0] = 'X'
	if string(decodeErr.Body()) != string(body) {
		t.Fatal("Body() returned mutable internal storage")
	}
}

func TestSemanticErrorHelpers(t *testing.T) {
	t.Parallel()

	if !bpi.RequiresVIP(&bpi.APIError{Code: -106}) {
		t.Fatal("RequiresVIP(-106) = false, want true")
	}
	if !bpi.IsPermissionError(&bpi.HTTPError{StatusCode: 403}) {
		t.Fatal("IsPermissionError(403) = false, want true")
	}
	if !bpi.IsRiskControl(&bpi.APIError{Code: -352}) {
		t.Fatal("IsRiskControl(-352) = false, want true")
	}
}

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "testdata", "fixtures", "envelope", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}
