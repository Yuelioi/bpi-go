package login_test

import (
	"encoding/json"
	"testing"

	"github.com/Yuelioi/bpi-go/login"
)

func TestNavNormalizesGuestIdentity(t *testing.T) {
	t.Parallel()

	var nav login.Nav
	if err := json.Unmarshal([]byte(`{"isLogin":false,"mid":0,"uname":" ","face":"","wbi_img":{"img_url":"img","sub_url":"sub"}}`), &nav); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if nav.IsLogin || nav.MID != nil || nav.Username != nil || nav.Face != nil {
		t.Fatalf("guest nav = %+v, want nil identity", nav)
	}
}
