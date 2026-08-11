package client_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/Yuelioi/bpi-go"
)

func TestAccountRepresentationsRedactCredentials(t *testing.T) {
	t.Parallel()

	account := completeAccount()
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	logger.Info("account", "value", account)

	jsonOutput, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("marshal account: %v", err)
	}
	representations := []string{
		account.String(),
		fmt.Sprintf("%v", account),
		fmt.Sprintf("%+v", account),
		string(jsonOutput),
		logOutput.String(),
	}
	for _, representation := range representations {
		for _, secret := range []string{"user-42", "session-secret", "csrf-secret", "buvid-secret"} {
			if strings.Contains(representation, secret) {
				t.Fatalf("representation leaked %q: %s", secret, representation)
			}
		}
	}
}

func TestAccountValidateRequiresCompleteCredentials(t *testing.T) {
	t.Parallel()

	if err := (bpi.Account{}).Validate(); err == nil {
		t.Fatal("Validate() error = nil, want incomplete-account error")
	}
	if err := completeAccount().Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func completeAccount() bpi.Account {
	return bpi.Account{
		DedeUserID: "user-42",
		SESSDATA:   "session-secret",
		BiliJCT:    "csrf-secret",
		Buvid3:     "buvid-secret",
	}
}
