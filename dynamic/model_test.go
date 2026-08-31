package dynamic_test

import (
	"encoding/json"
	"testing"

	"github.com/Yuelioi/bpi-go/dynamic"
)

func TestDetailItemAcceptsLegacyNumericAndDeletedOriginalIDs(t *testing.T) {
	t.Parallel()

	var item dynamic.DetailItem
	if err := json.Unmarshal([]byte(`{
		"id_str":188029375348892,
		"basic":{},
		"modules":{},
		"orig":{
			"id_str":null,
			"basic":{},
			"modules":{},
			"orig":null,
			"type":"DYNAMIC_TYPE_NONE",
			"visible":false
		},
		"type":"DYNAMIC_TYPE_AV",
		"visible":true
	}`), &item); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if item.ID != "188029375348892" {
		t.Fatalf("ID = %q, want normalized numeric ID", item.ID)
	}
	if item.Original == nil || item.Original.ID != "" {
		t.Fatalf("Original = %+v, want retained tombstone with empty ID", item.Original)
	}
}

func TestDetailItemRejectsNonScalarID(t *testing.T) {
	t.Parallel()

	var item dynamic.DetailItem
	if err := json.Unmarshal([]byte(`{"id_str":{},"basic":{},"modules":{},"orig":null,"type":"DYNAMIC_TYPE_AV","visible":true}`), &item); err == nil {
		t.Fatal("Unmarshal(object ID) error = nil")
	}
}
