// Package ids provides validated Bilibili identifier types.
package ids

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// ErrInvalid is matched by all identifier validation errors.
var ErrInvalid = errors.New("invalid Bilibili identifier")

// InvalidError describes an invalid identifier value.
type InvalidError struct {
	Field   string
	Message string
}

func (e *InvalidError) Error() string {
	return fmt.Sprintf("ids: invalid %s: %s", e.Field, e.Message)
}

func (e *InvalidError) Unwrap() error { return ErrInvalid }

type (
	// AID is a Bilibili numeric AV video identifier.
	AID uint64
	// AudioID is a Bilibili audio song identifier.
	AudioID uint64
	// CID is a Bilibili video page/content identifier.
	CID uint64
	// MID is a Bilibili member/user identifier.
	MID uint64
	// RoomID is a Bilibili live-room identifier.
	RoomID uint64
	// MediaID is a Bilibili media identifier.
	MediaID uint64
	// SeasonID is a Bilibili season identifier.
	SeasonID uint64
	// EpisodeID is a Bilibili episode identifier.
	EpisodeID uint64
	// NoteID is a Bilibili note identifier.
	NoteID uint64
	// CVID is a Bilibili note/article CV identifier.
	CVID uint64
)

func NewAID(value uint64) (AID, error)             { return newNumeric[AID]("aid", value) }
func NewAudioID(value uint64) (AudioID, error)     { return newNumeric[AudioID]("sid", value) }
func NewCID(value uint64) (CID, error)             { return newNumeric[CID]("cid", value) }
func NewMID(value uint64) (MID, error)             { return newNumeric[MID]("mid", value) }
func NewRoomID(value uint64) (RoomID, error)       { return newNumeric[RoomID]("room_id", value) }
func NewMediaID(value uint64) (MediaID, error)     { return newNumeric[MediaID]("media_id", value) }
func NewSeasonID(value uint64) (SeasonID, error)   { return newNumeric[SeasonID]("season_id", value) }
func NewEpisodeID(value uint64) (EpisodeID, error) { return newNumeric[EpisodeID]("ep_id", value) }
func NewNoteID(value uint64) (NoteID, error)       { return newNumeric[NoteID]("note_id", value) }
func NewCVID(value uint64) (CVID, error)           { return newNumeric[CVID]("cvid", value) }

func ParseAID(value string) (AID, error)             { return parseNumeric[AID]("aid", value) }
func ParseAudioID(value string) (AudioID, error)     { return parseNumeric[AudioID]("sid", value) }
func ParseCID(value string) (CID, error)             { return parseNumeric[CID]("cid", value) }
func ParseMID(value string) (MID, error)             { return parseNumeric[MID]("mid", value) }
func ParseRoomID(value string) (RoomID, error)       { return parseNumeric[RoomID]("room_id", value) }
func ParseMediaID(value string) (MediaID, error)     { return parseNumeric[MediaID]("media_id", value) }
func ParseSeasonID(value string) (SeasonID, error)   { return parseNumeric[SeasonID]("season_id", value) }
func ParseEpisodeID(value string) (EpisodeID, error) { return parseNumeric[EpisodeID]("ep_id", value) }
func ParseNoteID(value string) (NoteID, error)       { return parseNumeric[NoteID]("note_id", value) }
func ParseCVID(value string) (CVID, error)           { return parseNumeric[CVID]("cvid", value) }

func (id AID) Uint64() uint64       { return uint64(id) }
func (id AudioID) Uint64() uint64   { return uint64(id) }
func (id CID) Uint64() uint64       { return uint64(id) }
func (id MID) Uint64() uint64       { return uint64(id) }
func (id RoomID) Uint64() uint64    { return uint64(id) }
func (id MediaID) Uint64() uint64   { return uint64(id) }
func (id SeasonID) Uint64() uint64  { return uint64(id) }
func (id EpisodeID) Uint64() uint64 { return uint64(id) }
func (id NoteID) Uint64() uint64    { return uint64(id) }
func (id CVID) Uint64() uint64      { return uint64(id) }

func (id AID) String() string       { return formatNumeric(id) }
func (id AudioID) String() string   { return formatNumeric(id) }
func (id CID) String() string       { return formatNumeric(id) }
func (id MID) String() string       { return formatNumeric(id) }
func (id RoomID) String() string    { return formatNumeric(id) }
func (id MediaID) String() string   { return formatNumeric(id) }
func (id SeasonID) String() string  { return formatNumeric(id) }
func (id EpisodeID) String() string { return formatNumeric(id) }
func (id NoteID) String() string    { return formatNumeric(id) }
func (id CVID) String() string      { return formatNumeric(id) }

func (id AID) Validate() error       { return validateNumeric("aid", uint64(id)) }
func (id AudioID) Validate() error   { return validateNumeric("sid", uint64(id)) }
func (id CID) Validate() error       { return validateNumeric("cid", uint64(id)) }
func (id MID) Validate() error       { return validateNumeric("mid", uint64(id)) }
func (id RoomID) Validate() error    { return validateNumeric("room_id", uint64(id)) }
func (id MediaID) Validate() error   { return validateNumeric("media_id", uint64(id)) }
func (id SeasonID) Validate() error  { return validateNumeric("season_id", uint64(id)) }
func (id EpisodeID) Validate() error { return validateNumeric("ep_id", uint64(id)) }
func (id NoteID) Validate() error    { return validateNumeric("note_id", uint64(id)) }
func (id CVID) Validate() error      { return validateNumeric("cvid", uint64(id)) }

func (id AID) MarshalText() ([]byte, error)       { return marshalNumericText("aid", id) }
func (id AudioID) MarshalText() ([]byte, error)   { return marshalNumericText("sid", id) }
func (id CID) MarshalText() ([]byte, error)       { return marshalNumericText("cid", id) }
func (id MID) MarshalText() ([]byte, error)       { return marshalNumericText("mid", id) }
func (id RoomID) MarshalText() ([]byte, error)    { return marshalNumericText("room_id", id) }
func (id MediaID) MarshalText() ([]byte, error)   { return marshalNumericText("media_id", id) }
func (id SeasonID) MarshalText() ([]byte, error)  { return marshalNumericText("season_id", id) }
func (id EpisodeID) MarshalText() ([]byte, error) { return marshalNumericText("ep_id", id) }
func (id NoteID) MarshalText() ([]byte, error)    { return marshalNumericText("note_id", id) }
func (id CVID) MarshalText() ([]byte, error)      { return marshalNumericText("cvid", id) }

func (id *AID) UnmarshalText(text []byte) error     { return unmarshalNumericText("aid", text, id) }
func (id *AudioID) UnmarshalText(text []byte) error { return unmarshalNumericText("sid", text, id) }
func (id *CID) UnmarshalText(text []byte) error     { return unmarshalNumericText("cid", text, id) }
func (id *MID) UnmarshalText(text []byte) error     { return unmarshalNumericText("mid", text, id) }
func (id *RoomID) UnmarshalText(text []byte) error  { return unmarshalNumericText("room_id", text, id) }
func (id *MediaID) UnmarshalText(text []byte) error {
	return unmarshalNumericText("media_id", text, id)
}
func (id *SeasonID) UnmarshalText(text []byte) error {
	return unmarshalNumericText("season_id", text, id)
}
func (id *EpisodeID) UnmarshalText(text []byte) error { return unmarshalNumericText("ep_id", text, id) }
func (id *NoteID) UnmarshalText(text []byte) error    { return unmarshalNumericText("note_id", text, id) }
func (id *CVID) UnmarshalText(text []byte) error      { return unmarshalNumericText("cvid", text, id) }

func (id AID) MarshalJSON() ([]byte, error)       { return marshalNumericJSON("aid", id) }
func (id AudioID) MarshalJSON() ([]byte, error)   { return marshalNumericJSON("sid", id) }
func (id CID) MarshalJSON() ([]byte, error)       { return marshalNumericJSON("cid", id) }
func (id MID) MarshalJSON() ([]byte, error)       { return marshalNumericJSON("mid", id) }
func (id RoomID) MarshalJSON() ([]byte, error)    { return marshalNumericJSON("room_id", id) }
func (id MediaID) MarshalJSON() ([]byte, error)   { return marshalNumericJSON("media_id", id) }
func (id SeasonID) MarshalJSON() ([]byte, error)  { return marshalNumericJSON("season_id", id) }
func (id EpisodeID) MarshalJSON() ([]byte, error) { return marshalNumericJSON("ep_id", id) }
func (id NoteID) MarshalJSON() ([]byte, error)    { return marshalNumericJSON("note_id", id) }
func (id CVID) MarshalJSON() ([]byte, error)      { return marshalNumericJSON("cvid", id) }

func (id *AID) UnmarshalJSON(data []byte) error     { return unmarshalNumericJSON("aid", data, id) }
func (id *AudioID) UnmarshalJSON(data []byte) error { return unmarshalNumericJSON("sid", data, id) }
func (id *CID) UnmarshalJSON(data []byte) error     { return unmarshalNumericJSON("cid", data, id) }
func (id *MID) UnmarshalJSON(data []byte) error     { return unmarshalNumericJSON("mid", data, id) }
func (id *RoomID) UnmarshalJSON(data []byte) error  { return unmarshalNumericJSON("room_id", data, id) }
func (id *MediaID) UnmarshalJSON(data []byte) error {
	return unmarshalNumericJSON("media_id", data, id)
}
func (id *SeasonID) UnmarshalJSON(data []byte) error {
	return unmarshalNumericJSON("season_id", data, id)
}
func (id *EpisodeID) UnmarshalJSON(data []byte) error { return unmarshalNumericJSON("ep_id", data, id) }
func (id *NoteID) UnmarshalJSON(data []byte) error    { return unmarshalNumericJSON("note_id", data, id) }
func (id *CVID) UnmarshalJSON(data []byte) error      { return unmarshalNumericJSON("cvid", data, id) }

func newNumeric[T ~uint64](field string, value uint64) (T, error) {
	if err := validateNumeric(field, value); err != nil {
		return 0, err
	}
	return T(value), nil
}

func parseNumeric[T ~uint64](field, value string) (T, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, &InvalidError{Field: field, Message: "identifier must be an unsigned decimal integer"}
	}
	return newNumeric[T](field, parsed)
}

func validateNumeric(field string, value uint64) error {
	if value == 0 {
		return &InvalidError{Field: field, Message: "identifier must be non-zero"}
	}
	return nil
}

func formatNumeric[T ~uint64](value T) string {
	return strconv.FormatUint(uint64(value), 10)
}

func marshalNumericText[T ~uint64](field string, value T) ([]byte, error) {
	if err := validateNumeric(field, uint64(value)); err != nil {
		return nil, err
	}
	return []byte(formatNumeric(value)), nil
}

func unmarshalNumericText[T ~uint64](field string, text []byte, destination *T) error {
	if destination == nil {
		return &InvalidError{Field: field, Message: "cannot unmarshal into a nil identifier"}
	}
	value, err := parseNumeric[T](field, string(text))
	if err != nil {
		return err
	}
	*destination = value
	return nil
}

func marshalNumericJSON[T ~uint64](field string, value T) ([]byte, error) {
	if err := validateNumeric(field, uint64(value)); err != nil {
		return nil, err
	}
	return json.Marshal(uint64(value))
}

func unmarshalNumericJSON[T ~uint64](field string, data []byte, destination *T) error {
	if destination == nil {
		return &InvalidError{Field: field, Message: "cannot unmarshal into a nil identifier"}
	}
	var raw uint64
	if err := json.Unmarshal(data, &raw); err != nil {
		return &InvalidError{Field: field, Message: "identifier must be an unsigned JSON integer"}
	}
	value, err := newNumeric[T](field, raw)
	if err != nil {
		return err
	}
	*destination = value
	return nil
}
