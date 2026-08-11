package article

import (
	"errors"
	"net/url"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedQueries(t *testing.T) {
	t.Parallel()
	id, _ := ids.NewCVID(2)
	cards, err := NewCardsParams("av2,cv1,cv2")
	if err != nil {
		t.Fatalf("NewCardsParams() error = %v", err)
	}
	articles, err := NewArticlesInfoParams(207_146)
	if err != nil {
		t.Fatalf("NewArticlesInfoParams() error = %v", err)
	}
	tests := []struct {
		name   string
		encode func() (url.Values, error)
		want   string
	}{
		{"info", NewInfoParams(id).EncodeQuery, "id=2"},
		{"view", NewViewParams(id).EncodeQuery, "gaia_source=main_web&id=2"},
		{"cards", cards.EncodeQuery, "ids=av2%2Ccv1%2Ccv2&web_location=333.1305"},
		{"articles", articles.EncodeQuery, "id=207146"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			query, err := test.encode()
			if err != nil || query.Encode() != test.want {
				t.Fatalf("EncodeQuery() = %q, %v; want %q", query.Encode(), err, test.want)
			}
		})
	}
}

func TestPromotedParametersRejectInvalidValues(t *testing.T) {
	t.Parallel()
	id, _ := ids.NewCVID(2)
	tests := []func() error{
		func() error { _, err := NewInfoParams(0).EncodeQuery(); return err },
		func() error { _, err := NewViewParams(0).EncodeQuery(); return err },
		func() error { _, err := NewViewParams(id).WithGaiaSource(" "); return err },
		func() error { _, err := NewCardsParams(" "); return err },
		func() error { _, err := NewArticlesInfoParams(0); return err },
	}
	for index, run := range tests {
		var parameterError *bpierr.ParameterError
		if err := run(); !errors.As(err, &parameterError) {
			t.Fatalf("case %d error = %v, want ParameterError", index, err)
		}
	}
}
