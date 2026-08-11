package cheese

import (
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestEpisodeListParamsEncodeQuery(t *testing.T) {
	seasonID, _ := ids.NewSeasonID(556)
	params, err := NewEpisodeListParams(seasonID).WithPageSize(50)
	if err != nil {
		t.Fatalf("WithPageSize() error = %v", err)
	}
	params, err = params.WithPage(1)
	if err != nil {
		t.Fatalf("WithPage() error = %v", err)
	}
	query, err := params.EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if got := query.Encode(); got != "pn=1&ps=50&season_id=556" {
		t.Fatalf("query = %q", got)
	}
}

func TestPlayURLParamsEncodeQuery(t *testing.T) {
	aid, _ := ids.NewAID(997_984_154)
	episodeID, _ := ids.NewEpisodeID(163_956)
	cid, _ := ids.NewCID(1_183_682_680)
	query, err := NewPlayURLParams(aid, episodeID, cid).
		WithQuality(Quality480P).
		WithFormatFlags(FormatDASH).
		EncodeQuery()
	if err != nil {
		t.Fatalf("EncodeQuery() error = %v", err)
	}
	if got := query.Encode(); got != "avid=997984154&cid=1183682680&ep_id=163956&fnval=16&fnver=0&qn=32" {
		t.Fatalf("query = %q", got)
	}
}

func TestCheeseParamsRejectInvalidValues(t *testing.T) {
	var parameterError *bpierr.ParameterError
	if _, err := NewEpisodeListParams(0).EncodeQuery(); !errors.As(err, &parameterError) {
		t.Fatalf("EpisodeListParams error = %v", err)
	}
	parameterError = nil
	if _, err := NewEpisodeListParams(1).WithPageSize(0); !errors.As(err, &parameterError) {
		t.Fatalf("WithPageSize(0) error = %v", err)
	}
	parameterError = nil
	if _, err := NewPlayURLParams(0, 0, 0).EncodeQuery(); !errors.As(err, &parameterError) {
		t.Fatalf("PlayURLParams error = %v", err)
	}
}
