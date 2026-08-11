package ids_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestNumericConstructorsRejectZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func() error
	}{
		{name: "aid", call: func() error { _, err := ids.NewAID(0); return err }},
		{name: "audio", call: func() error { _, err := ids.NewAudioID(0); return err }},
		{name: "cid", call: func() error { _, err := ids.NewCID(0); return err }},
		{name: "mid", call: func() error { _, err := ids.NewMID(0); return err }},
		{name: "room", call: func() error { _, err := ids.NewRoomID(0); return err }},
		{name: "media", call: func() error { _, err := ids.NewMediaID(0); return err }},
		{name: "season", call: func() error { _, err := ids.NewSeasonID(0); return err }},
		{name: "episode", call: func() error { _, err := ids.NewEpisodeID(0); return err }},
		{name: "note", call: func() error { _, err := ids.NewNoteID(0); return err }},
		{name: "cvid", call: func() error { _, err := ids.NewCVID(0); return err }},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if err := test.call(); !errors.Is(err, ids.ErrInvalid) {
				t.Fatalf("constructor error = %v, want ErrInvalid", err)
			}
		})
	}
}

func TestNumericIDsFormatParseAndJSONRoundTrip(t *testing.T) {
	t.Parallel()

	aid, err := ids.NewAID(170001)
	if err != nil {
		t.Fatalf("NewAID() error = %v", err)
	}
	if aid.String() != "170001" || aid.Uint64() != 170001 {
		t.Fatalf("AID = %s/%d, want 170001", aid, aid.Uint64())
	}
	parsed, err := ids.ParseAID("170001")
	if err != nil || parsed != aid {
		t.Fatalf("ParseAID() = %v, %v; want %v, nil", parsed, err, aid)
	}
	data, err := json.Marshal(aid)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != "170001" {
		t.Fatalf("Marshal() = %s, want 170001", data)
	}
	var decoded ids.AID
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if decoded != aid {
		t.Fatalf("decoded = %v, want %v", decoded, aid)
	}
}

func TestNumericIDsMatchRustDisplayVectors(t *testing.T) {
	t.Parallel()

	audioID, err := ids.NewAudioID(13603)
	if err != nil {
		t.Fatalf("NewAudioID() error = %v", err)
	}
	episodeID, err := ids.NewEpisodeID(21265)
	if err != nil {
		t.Fatalf("NewEpisodeID() error = %v", err)
	}
	if audioID.String() != "13603" || episodeID.String() != "21265" {
		t.Fatalf("display vectors = %s/%s, want 13603/21265", audioID, episodeID)
	}
}

func TestNumericIDErrorsPreserveRustFieldNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		construct func() error
		wantField string
	}{
		{name: "aid", construct: func() error { _, err := ids.NewAID(0); return err }, wantField: "aid"},
		{name: "audio", construct: func() error { _, err := ids.NewAudioID(0); return err }, wantField: "sid"},
		{name: "season", construct: func() error { _, err := ids.NewSeasonID(0); return err }, wantField: "season_id"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.construct()
			var invalid *ids.InvalidError
			if !errors.As(err, &invalid) || invalid.Field != test.wantField {
				t.Fatalf("constructor error = %v, want field %q", err, test.wantField)
			}
		})
	}
}

func TestZeroValueNumericIDCannotBeSerialized(t *testing.T) {
	t.Parallel()

	if _, err := json.Marshal(ids.AID(0)); !errors.Is(err, ids.ErrInvalid) {
		t.Fatalf("Marshal(zero AID) error = %v, want ErrInvalid", err)
	}
}

func TestBVIDMatchesRustValidationAndRoundTrips(t *testing.T) {
	t.Parallel()

	bvid, err := ids.NewBVID("BV1bx411c7ux")
	if err != nil {
		t.Fatalf("ParseBVID() error = %v", err)
	}
	if bvid.String() != "BV1bx411c7ux" {
		t.Fatalf("BVID = %q, want BV1bx411c7ux", bvid)
	}
	data, err := json.Marshal(bvid)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	var decoded ids.BVID
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if decoded != bvid {
		t.Fatalf("decoded = %q, want %q", decoded, bvid)
	}

	for _, invalid := range []string{"", "av170001", "BV1!", "BV123 456789"} {
		if _, err := ids.ParseBVID(invalid); !errors.Is(err, ids.ErrInvalid) {
			t.Errorf("ParseBVID(%q) error = %v, want ErrInvalid", invalid, err)
		}
	}
}

func TestDynamicIDRejectsBlankValue(t *testing.T) {
	t.Parallel()

	if _, err := ids.NewDynamicID("   "); !errors.Is(err, ids.ErrInvalid) {
		t.Fatalf("ParseDynamicID(blank) error = %v, want ErrInvalid", err)
	}
}
