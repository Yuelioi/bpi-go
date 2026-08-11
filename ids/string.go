package ids

import (
	"encoding/json"
	"strings"
)

type (
	// BVID is a validated Bilibili BV video identifier.
	BVID string
	// DynamicID is a validated Bilibili dynamic-item identifier.
	DynamicID string
)

// NewBVID validates and returns a BVID.
func NewBVID(value string) (BVID, error) { return ParseBVID(value) }

// ParseBVID validates and returns a BVID.
func ParseBVID(value string) (BVID, error) {
	if !strings.HasPrefix(value, "BV") {
		return "", &InvalidError{Field: "bvid", Message: "identifier must start with 'BV'"}
	}
	if len(value) < 12 || !isASCIIAlphanumeric(value) {
		return "", &InvalidError{Field: "bvid", Message: "identifier must contain at least 12 ASCII alphanumeric characters"}
	}
	return BVID(value), nil
}

// NewDynamicID validates and returns a DynamicID.
func NewDynamicID(value string) (DynamicID, error) { return ParseDynamicID(value) }

// ParseDynamicID validates and returns a DynamicID.
func ParseDynamicID(value string) (DynamicID, error) {
	if strings.TrimSpace(value) == "" {
		return "", &InvalidError{Field: "dynamic_id", Message: "identifier cannot be blank"}
	}
	return DynamicID(value), nil
}

func (id BVID) String() string      { return string(id) }
func (id DynamicID) String() string { return string(id) }

func (id BVID) Validate() error {
	_, err := ParseBVID(string(id))
	return err
}

func (id DynamicID) Validate() error {
	_, err := ParseDynamicID(string(id))
	return err
}

func (id BVID) MarshalText() ([]byte, error) { return marshalStringText(id.Validate(), string(id)) }
func (id DynamicID) MarshalText() ([]byte, error) {
	return marshalStringText(id.Validate(), string(id))
}

func (id *BVID) UnmarshalText(text []byte) error {
	if id == nil {
		return &InvalidError{Field: "bvid", Message: "cannot unmarshal into a nil identifier"}
	}
	value, err := ParseBVID(string(text))
	if err != nil {
		return err
	}
	*id = value
	return nil
}

func (id *DynamicID) UnmarshalText(text []byte) error {
	if id == nil {
		return &InvalidError{Field: "dynamic_id", Message: "cannot unmarshal into a nil identifier"}
	}
	value, err := ParseDynamicID(string(text))
	if err != nil {
		return err
	}
	*id = value
	return nil
}

func (id BVID) MarshalJSON() ([]byte, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(id))
}

func (id DynamicID) MarshalJSON() ([]byte, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(id))
}

func (id *BVID) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return &InvalidError{Field: "bvid", Message: "identifier must be a JSON string"}
	}
	return id.UnmarshalText([]byte(value))
}

func (id *DynamicID) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return &InvalidError{Field: "dynamic_id", Message: "identifier must be a JSON string"}
	}
	return id.UnmarshalText([]byte(value))
}

func marshalStringText(validation error, value string) ([]byte, error) {
	if validation != nil {
		return nil, validation
	}
	return []byte(value), nil
}

func isASCIIAlphanumeric(value string) bool {
	for _, character := range []byte(value) {
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') {
			continue
		}
		return false
	}
	return true
}
