package user

import (
	"encoding/json"
	"testing"
)

func TestFlexibleMemberIDsAndTransparentNotice(t *testing.T) {
	t.Parallel()
	var card CardProfile
	if err := json.Unmarshal([]byte(`{"card":{"mid":"2","name":"user","face":"face"},"following":false}`), &card); err != nil {
		t.Fatalf("decode card: %v", err)
	}
	if card.Card.MID.Uint64() != 2 {
		t.Fatalf("card MID = %d", card.Card.MID)
	}
	var lookup NameToUID
	if err := json.Unmarshal([]byte(`{"uid_list":[{"name":"one","uid":"2"},{"name":"two","uid":3}]}`), &lookup); err != nil {
		t.Fatalf("decode lookup: %v", err)
	}
	if len(lookup.Items) != 2 || lookup.Items[1].MID.Uint64() != 3 {
		t.Fatalf("lookup = %+v", lookup)
	}
	var notice SpaceNotice
	if err := json.Unmarshal([]byte(`"hello"`), &notice); err != nil || notice.Content != "hello" {
		t.Fatalf("notice = %+v, %v", notice, err)
	}
}
