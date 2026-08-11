package fav

import (
	"testing"

	"github.com/Yuelioi/bpi-go/ids"
)

func TestParamsEncodeContractQueries(t *testing.T) {
	t.Parallel()
	mediaID, _ := ids.NewMediaID(1052622027)
	mid, _ := ids.NewMID(7792521)
	resources, _ := NewResourceInfosParams("371494037:2")
	detail, _ := NewListDetailParams(mediaID).WithOrder("mtime")
	detail = detail.WithContentType(0)
	detail, _ = detail.WithPageSize(5)
	detail, _ = detail.WithPage(1)
	tests := []struct {
		name string
		got  func() string
		want string
	}{
		{"folder-info", func() string { q, _ := NewFolderInfoParams(mediaID).EncodeQuery(); return q.Encode() }, "media_id=1052622027"},
		{"created-list", func() string { q, _ := NewCreatedListParams(mid).EncodeQuery(); return q.Encode() }, "up_mid=7792521&web_location=333.1387"},
		{"collected-list", func() string { q, _ := NewCollectedListParams(mid).EncodeQuery(); return q.Encode() }, "platform=web&pn=1&ps=20&up_mid=7792521"},
		{"resource-infos", func() string { q, _ := resources.EncodeQuery(); return q.Encode() }, "platform=web&resources=371494037%3A2"},
		{"list-detail", func() string { q, _ := detail.EncodeQuery(); return q.Encode() }, "media_id=1052622027&order=mtime&platform=web&pn=1&ps=5&type=0"},
		{"resource-ids", func() string { q, _ := NewResourceIDsParams(mediaID).EncodeQuery(); return q.Encode() }, "media_id=1052622027&platform=web"},
	}
	for _, test := range tests {
		if got := test.got(); got != test.want {
			t.Errorf("%s query = %q, want %q", test.name, got, test.want)
		}
	}
}

func TestParamsRejectInvalidValues(t *testing.T) {
	t.Parallel()
	if _, err := NewResourceInfosParams("  "); err == nil {
		t.Fatal("NewResourceInfosParams(blank) error = nil")
	}
	if _, err := NewListDetailParams(ids.MediaID(1)).WithPage(0); err == nil {
		t.Fatal("WithPage(0) error = nil")
	}
	if _, err := NewCollectedListParams(ids.MID(1)).WithPageSize(0); err == nil {
		t.Fatal("WithPageSize(0) error = nil")
	}
}
