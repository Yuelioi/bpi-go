package live

import (
	"net/url"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestPromotedRequestParameters(t *testing.T) {
	roomID, _ := ids.NewRoomID(23_174_842)
	streamRoomID, _ := ids.NewRoomID(14_073_662)
	danmuRoomID, _ := ids.NewRoomID(21_733_448)
	emoticonRoomID, _ := ids.NewRoomID(14_047)
	moderationRoomID, _ := ids.NewRoomID(3_818_081)
	anchorID, _ := ids.NewMID(504_140_200)
	moderationAnchorID, _ := ids.NewMID(1_000_001)

	stream, err := NewStreamParams(streamRoomID).WithPlatform("web")
	if err != nil {
		t.Fatal(err)
	}
	stream, err = stream.WithQN(10_000)
	if err != nil {
		t.Fatal(err)
	}
	blind, err := NewBlindGiftInfoParams(32_251)
	if err != nil {
		t.Fatal(err)
	}
	follow := NewFollowUpListParams()
	follow, err = follow.WithPage(1)
	if err != nil {
		t.Fatal(err)
	}
	follow, err = follow.WithPageSize(2)
	if err != nil {
		t.Fatal(err)
	}
	follow, err = follow.WithIgnoreRecord(1)
	if err != nil {
		t.Fatal(err)
	}
	follow = follow.WithHitAB(true)
	replay := NewReplayListParams()
	replay, err = replay.WithPage(1)
	if err != nil {
		t.Fatal(err)
	}
	replay, err = replay.WithPageSize(2)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		got  func() (url.Values, error)
		want string
	}{
		{"room-info", NewRoomInfoParams(roomID).EncodeQuery, "room_id=23174842"},
		{"stream", stream.EncodeQuery, "cid=14073662&platform=web&qn=10000"},
		{"room-gift-list", NewRoomGiftListParams(roomID).EncodeQuery, "platform=web&room_id=23174842"},
		{"blind-gift", blind.EncodeQuery, "gift_id=32251"},
		{"danmu-info", NewDanmuInfoParams(danmuRoomID).EncodeQuery, "id=21733448&type=0"},
		{"emoticons", NewEmoticonsParams(emoticonRoomID).EncodeQuery, "platform=pc&room_id=14047"},
		{"lottery", NewLotteryInfoParams(roomID).EncodeQuery, "roomid=23174842"},
		{"medals", NewMyMedalsParams().EncodeQuery, "page=1&page_size=10"},
		{"follow-up", follow.EncodeQuery, "hit_ab=true&ignoreRecord=1&page=1&page_size=2"},
		{"follow-up-web", NewFollowUpWebListParams().WithHitAB(false).EncodeQuery, "hit_ab=false"},
		{"replay", replay.EncodeQuery, "page=1&page_size=2"},
		{"guard", NewGuardListParams(roomID, anchorID).EncodeQuery, "page=1&page_size=20&roomid=23174842&ruid=504140200&typ=5"},
		{"heartbeat", NewWebHeartBeatParams(roomID).EncodeQuery, "hb=NjB8MjMxNzQ4NDJ8MXww&pf=web"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values, err := test.got()
			if err != nil {
				t.Fatalf("EncodeQuery() error = %v", err)
			}
			if got := values.Encode(); got != test.want {
				t.Fatalf("EncodeQuery() = %q, want %q", got, test.want)
			}
		})
	}

	csrf := "fixture-csrf"
	silent := NewSilentUsersParams(moderationRoomID)
	form, err := silent.EncodeForm(csrf)
	if err != nil || form.Encode() != "csrf=fixture-csrf&csrf_token=fixture-csrf&pn=1&ps=10&room_id=3818081" {
		t.Fatalf("SilentUsers form = %q, %v", form.Encode(), err)
	}
	banned := NewBannedUsersParams(moderationRoomID, moderationAnchorID)
	query, err := banned.EncodeQuery(csrf)
	if err != nil || query.Encode() != "anchor_id=1000001&csrf=fixture-csrf&csrf_token=fixture-csrf&mobi_app=android&platform=android&pn=1&ps=10&spmid=444.8.0.0&visit_id=" {
		t.Fatalf("BannedUsers query = %q, %v", query.Encode(), err)
	}
	shield := NewShieldKeywordsParams(moderationRoomID)
	form, err = shield.EncodeForm(csrf)
	if err != nil || form.Encode() != "csrf=fixture-csrf&csrf_token=fixture-csrf&mobi_app=android&platform=android&room_id=3818081&spmid=444.8.0.0&visit_id=" {
		t.Fatalf("ShieldKeywords form = %q, %v", form.Encode(), err)
	}
	for name, referer := range map[string]func() (string, error){
		"silent": silent.Referer,
		"banned": banned.Referer,
		"shield": shield.Referer,
	} {
		got, err := referer()
		if err != nil || got != "https://live.bilibili.com/3818081" {
			t.Fatalf("%s Referer() = %q, %v", name, got, err)
		}
	}
}

func TestPromotedRequestParametersRejectInvalidValues(t *testing.T) {
	zeroRoom := ids.RoomID(0)
	if _, err := NewRoomInfoParams(zeroRoom).EncodeQuery(); err == nil {
		t.Fatal("zero room ID was accepted")
	}
	if _, err := NewBlindGiftInfoParams(0); err == nil {
		t.Fatal("zero gift ID was accepted")
	}
	if _, err := NewStreamParams(zeroRoom).WithPlatform(" "); err == nil {
		t.Fatal("blank platform was accepted")
	}
	if _, err := NewMyMedalsParams().WithPage(0); err == nil {
		t.Fatal("zero page was accepted")
	}
	if _, err := NewFollowUpListParams().WithIgnoreRecord(2); err == nil {
		t.Fatal("unsupported ignoreRecord was accepted")
	}
	if _, err := NewWebHeartBeatParams(zeroRoom).EncodeQuery(); err == nil {
		t.Fatal("zero heartbeat room ID was accepted")
	}
}
