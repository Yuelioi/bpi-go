package probe

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestCommittedSnapshotAuditAndCatalog(t *testing.T) {
	root := filepath.Join("..", "..")
	report, err := Audit(
		filepath.Join(root, "parity", "source.json"),
		filepath.Join(root, "parity", "implemented.json"),
		filepath.Join(root, "testdata", "contracts"),
		filepath.Join(root, "parity", "contracts.lock.json"),
	)
	if err != nil {
		t.Fatalf("Audit() error = %v", err)
	}
	if report.Contracts != 206 || report.Mappings != 206 || report.Files != 642 || report.RiskCounts["public-read"] != 115 {
		t.Fatalf("Audit() = %+v", report)
	}
	first, err := GenerateCatalog(filepath.Join(root, "parity", "source.json"), filepath.Join(root, "parity", "implemented.json"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateCatalog(filepath.Join(root, "parity", "source.json"), filepath.Join(root, "parity", "implemented.json"))
	if err != nil {
		t.Fatal(err)
	}
	if first != second || !strings.Contains(first, "对齐进度：**206/206 条契约**") || !strings.Contains(first, "`LiveClient.WebHeartBeat`") {
		t.Fatal("catalog generation is incomplete or non-deterministic")
	}
}

func TestContractLockIsDeterministic(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "contracts")
	first, err := BuildContractLock(root, SourceCommit)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildContractLock(root, SourceCommit)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) || len(first.Files) != 642 {
		t.Fatalf("contract locks differ or have wrong size: %d/%d", len(first.Files), len(second.Files))
	}
}

func TestSanitizationAuditRejectsCredentials(t *testing.T) {
	root := t.TempDir()
	responses := filepath.Join(root, "sample", "responses")
	if err := os.MkdirAll(responses, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(responses, "leak.json")
	if err := os.WriteFile(path, []byte(`{"code":0,"data":{"bili_jct":"real-secret","note":"SESSDATA=also-secret"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	findings, err := AuditSanitization(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) < 2 {
		t.Fatalf("findings = %+v, want credential field and Cookie assignment", findings)
	}
}

func TestBatchRunRequiresNetworkAndReadOnlyGates(t *testing.T) {
	t.Setenv("BPI_PROBE", "")
	_, err := RunBatch(context.Background(), BatchConfig{ContractRoot: t.TempDir(), Profiles: []string{"anonymous"}, ReadOnly: true})
	if err == nil || !strings.Contains(err.Error(), "BPI_PROBE=1") {
		t.Fatalf("RunBatch() error = %v", err)
	}
	t.Setenv("BPI_PROBE", "1")
	_, err = RunBatch(context.Background(), BatchConfig{ContractRoot: t.TempDir(), Profiles: []string{"anonymous"}})
	if err == nil || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("RunBatch() error = %v", err)
	}
}

func TestBatchClientsUseExplicitCookieHeaders(t *testing.T) {
	t.Parallel()

	clients, err := buildBatchClients(BatchConfig{
		CookieHeaders: map[string]string{
			"normal": "SESSDATA=session; bili_jct=csrf; device_cookie=preserved",
		},
	}, []string{"normal"}, time.Now)
	if err != nil {
		t.Fatalf("buildBatchClients() error = %v", err)
	}
	client := clients["normal"].client
	if !client.HasLoginCookies() {
		t.Fatal("normal client has no login Cookie")
	}
	if csrf, err := client.CSRF(); err != nil || csrf != "csrf" {
		t.Fatalf("CSRF() = %q, %v; want csrf, nil", csrf, err)
	}

	_, err = buildBatchClients(BatchConfig{}, []string{"vip"}, time.Now)
	if err == nil || !strings.Contains(err.Error(), "Cookie header is required") {
		t.Fatalf("buildBatchClients(missing Cookie) error = %v", err)
	}
}

func TestBatchRunProducesBodyFreeSummaryAndSignsWBI(t *testing.T) {
	t.Setenv("BPI_PROBE", "1")
	root := t.TempDir()
	contractDir := filepath.Join(root, "video", "probe")
	if err := os.MkdirAll(contractDir, 0o755); err != nil {
		t.Fatal(err)
	}
	contract := Contract{
		SchemaVersion: 2,
		Name:          "video.probe",
		Module:        "video",
		Batch:         "test",
		Endpoint:      "probe",
		Risk:          "public-read",
		Status:        "promoted",
		Profiles:      []string{"anonymous"},
		Request: ContractRequest{
			Method: http.MethodGet,
			URL:    "https://api.bilibili.com/x/test/wbi",
			Query:  map[string]string{"id": "1"},
			Auth:   ContractAuth{Requires: []string{"wbi"}},
		},
		Cases: []ContractCase{{Name: "anonymous", Profile: "anonymous", Response: ContractResponse{APICode: intPointer(0)}}},
	}
	data, err := json.Marshal(contract)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contractDir, "contract.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		calls.Add(1)
		switch request.URL.Path {
		case "/x/web-interface/nav":
			return jsonResponse(`{"code":-101,"data":{"wbi_img":{"img_url":"https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png","sub_url":"https://i0.hdslb.com/bfs/wbi/ABCDEFGHIJKLMNOPQRSTUVWXYZ654321.png"}}}`), nil
		case "/x/test/wbi":
			query := request.URL.Query()
			if query.Get("id") != "1" || query.Get("wts") != "1700000000" || len(query.Get("w_rid")) != 32 {
				t.Fatalf("signed query = %v", query)
			}
			return jsonResponse(`{"code":0,"data":{"private_note":"must-never-enter-summary"}}`), nil
		default:
			t.Fatalf("unexpected request %s", request.URL)
			return nil, nil
		}
	})}
	fixed := time.Unix(1_700_000_000, 0)
	summary, err := RunBatch(context.Background(), BatchConfig{
		ContractRoot: root,
		Profiles:     []string{"anonymous"},
		ReadOnly:     true,
		HTTPClient:   httpClient,
		Now:          func() time.Time { return fixed },
	})
	if err != nil {
		t.Fatalf("RunBatch() error = %v", err)
	}
	if calls.Load() != 2 || summary.Total != 1 || summary.Passed != 1 || summary.Results[0].Outcome != "passed" {
		t.Fatalf("summary = %+v, calls = %d", summary, calls.Load())
	}
	encoded, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "must-never-enter-summary") || summary.Results[0].SHA256 == "" || summary.Results[0].Bytes == 0 {
		t.Fatalf("unsafe or incomplete summary: %s", encoded)
	}
}

func intPointer(value int) *int { return &value }

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
