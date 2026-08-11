// Package contracttest contains shared adapters for cross-package domain
// contract tests. It is not part of the SDK's production interface.
package contracttest

import (
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Yuelioi/bpi-go"
)

// WBINavigationBody is a deterministic navigation response used by WBI tests.
const WBINavigationBody = `{"code":-101,"data":{"wbi_img":{"img_url":"https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png","sub_url":"https://i0.hdslb.com/bfs/wbi/ABCDEFGHIJKLMNOPQRSTUVWXYZ654321.png"}}}`

type testingTB interface {
	Helper()
	Fatalf(format string, args ...any)
}

// Fixture reads one committed promoted-contract fixture from the repository.
func Fixture(t testingTB, parts ...string) []byte {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("locate contract-test helper")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(source), "..", ".."))
	pathParts := append([]string{repositoryRoot, "testdata", "contracts"}, parts...)
	data, err := os.ReadFile(filepath.Join(pathParts...))
	if err != nil {
		t.Fatalf("read contract fixture: %v", err)
	}
	return data
}

// ProfileOption returns deterministic explicit credentials for authenticated
// fixture profiles and a behavior-neutral option for anonymous profiles.
func ProfileOption(profile string) bpi.Option {
	if profile == "anonymous" {
		return bpi.WithTimeout(10 * time.Second)
	}
	return bpi.WithAccount(bpi.Account{
		DedeUserID: "1",
		SESSDATA:   "fixture-session",
		BiliJCT:    "fixture-csrf",
		Buvid3:     "fixture-buvid",
	})
}

// AssertQuery verifies an exact single-valued query.
func AssertQuery(t testingTB, got url.Values, want map[string]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("query = %v, want %v", got, want)
	}
	for key, value := range want {
		if len(got[key]) != 1 || got[key][0] != value {
			t.Fatalf("query[%q] = %v, want %q", key, got[key], value)
		}
	}
}

// AssertWBIFields verifies unsigned fields plus the generated WBI metadata.
func AssertWBIFields(t testingTB, query url.Values, fields map[string]string) {
	t.Helper()
	for key, value := range fields {
		if query.Get(key) != value {
			t.Fatalf("query[%q] = %q, want %q", key, query.Get(key), value)
		}
	}
	if query.Get("wts") == "" || len(query.Get("w_rid")) != 32 || len(query) != len(fields)+2 {
		t.Fatalf("signed query = %v", query)
	}
}
