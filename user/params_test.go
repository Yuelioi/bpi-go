package user

import (
	"errors"
	"net/url"
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/internal/bpierr"
)

func TestPromotedParameterQueries(t *testing.T) {
	t.Parallel()
	mid, _ := ids.NewMID(2)
	other, _ := ids.NewMID(3)

	cards, err := NewCardsParams(mid, other)
	if err != nil {
		t.Fatalf("NewCardsParams() error = %v", err)
	}
	infos, err := NewInfosParams(mid, other)
	if err != nil {
		t.Fatalf("NewInfosParams() error = %v", err)
	}
	names, err := NewNameToUIDParams(" LexBurner ", "某科学")
	if err != nil {
		t.Fatalf("NewNameToUIDParams() error = %v", err)
	}
	bangumi := NewBangumiFollowListParams(mid)
	followings := NewFollowingsParams(mid)
	followings, _ = followings.WithOrderType("attention")
	followings, _ = followings.WithPageSize(20)
	followings, _ = followings.WithPage(1)
	followers := NewFollowersParams(mid)
	followers, _ = followers.WithPageSize(20)
	followers, _ = followers.WithPage(1)
	uploaded := NewUploadedVideosParams(mid)

	tests := []struct {
		name   string
		encode func() (url.Values, error)
		want   url.Values
	}{
		{"card", NewCardParams(mid).WithPhoto(CardPhotoInclude).EncodeQuery, url.Values{"mid": {"2"}, "photo": {"true"}}},
		{"cards", cards.EncodeQuery, url.Values{"uids": {"2,3"}}},
		{"infos", infos.EncodeQuery, url.Values{"uids": {"2,3"}}},
		{"bangumi", bangumi.EncodeQuery, url.Values{"vmid": {"2"}, "type": {"1"}, "pn": {"1"}, "ps": {"15"}}},
		{"followings", followings.EncodeQuery, url.Values{"vmid": {"2"}, "order_type": {"attention"}, "ps": {"20"}, "pn": {"1"}}},
		{"followers", followers.EncodeQuery, url.Values{"vmid": {"2"}, "ps": {"20"}, "pn": {"1"}}},
		{"name-to-uid", names.EncodeQuery, url.Values{"names": {"LexBurner,某科学"}}},
		{"uploaded-videos", uploaded.EncodeQuery, url.Values{"mid": {"2"}, "order": {"pubdate"}, "tid": {"0"}, "pn": {"1"}, "ps": {"30"}}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := test.encode()
			if err != nil {
				t.Fatalf("EncodeQuery() error = %v", err)
			}
			if got.Encode() != test.want.Encode() {
				t.Fatalf("EncodeQuery() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestParameterValidationRejectsInvalidZeroValues(t *testing.T) {
	t.Parallel()
	mid, _ := ids.NewMID(2)
	tests := []struct {
		name string
		run  func() error
	}{
		{"empty cards", func() error { _, err := NewCardsParams(); return err }},
		{"empty infos", func() error { _, err := NewInfosParams(); return err }},
		{"blank name", func() error { _, err := NewNameToUIDParams(" "); return err }},
		{"zero bangumi page", func() error { _, err := NewBangumiFollowListParams(mid).WithPage(0); return err }},
		{"large bangumi page size", func() error { _, err := NewBangumiFollowListParams(mid).WithPageSize(31); return err }},
		{"zero following page", func() error { _, err := NewFollowingsParams(mid).WithPage(0); return err }},
		{"zero follower size", func() error { _, err := NewFollowersParams(mid).WithPageSize(0); return err }},
		{"zero uploaded page", func() error { _, err := NewUploadedVideosParams(mid).WithPage(0); return err }},
		{"invalid uploaded order", func() error { _, err := ParseUploadedVideoOrder("newest"); return err }},
		{"zero member", func() error { _, err := NewCardParams(0).EncodeQuery(); return err }},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var parameterError *bpierr.ParameterError
			if err := test.run(); !errors.As(err, &parameterError) {
				t.Fatalf("error = %v, want ParameterError", err)
			}
		})
	}
}
