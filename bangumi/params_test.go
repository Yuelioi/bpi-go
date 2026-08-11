package bangumi

import (
	"errors"
	"net/url"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedReadQueries(t *testing.T) {
	t.Parallel()
	mediaID, _ := ids.NewMediaID(28_220_978)
	seasonID, _ := ids.NewSeasonID(1_172)
	episodeID, _ := ids.NewEpisodeID(21_265)
	play := PlayURLByEpisodeID(episodeID).WithQuality(Quality480P).WithFormatFlags(FormatDASH)
	tests := []struct {
		name   string
		encode func() (url.Values, error)
		want   string
	}{
		{"info", NewInfoParams(mediaID).EncodeQuery, "media_id=28220978"},
		{"season detail", DetailBySeasonID(seasonID).EncodeQuery, "season_id=1172"},
		{"episode detail", DetailByEpisodeID(episodeID).EncodeQuery, "ep_id=21265"},
		{"sections", NewSectionsParams(seasonID).EncodeQuery, "season_id=1172"},
		{"play URL", play.EncodeQuery, "ep_id=21265&fnval=16&fnver=0&qn=32"},
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

func TestPromotedReadParamsRejectInvalidIDs(t *testing.T) {
	t.Parallel()
	tests := []func() error{
		func() error { _, err := NewInfoParams(0).EncodeQuery(); return err },
		func() error { _, err := (DetailParams{}).EncodeQuery(); return err },
		func() error { _, err := NewSectionsParams(0).EncodeQuery(); return err },
		func() error { _, err := (PlayURLParams{}).EncodeQuery(); return err },
		func() error { _, err := PlayURLByEpisodeID(1).WithQuality(0).EncodeQuery(); return err },
	}
	for index, run := range tests {
		var parameterError *bpierr.ParameterError
		if err := run(); !errors.As(err, &parameterError) {
			t.Fatalf("case %d error = %v, want ParameterError", index, err)
		}
	}
}
