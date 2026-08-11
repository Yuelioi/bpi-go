package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Yuelioi/bpi-go/internal/testutil"
)

func TestQueryAndFormPoliciesPreserveCallerValues(t *testing.T) {
	t.Parallel()

	query := url.Values{"page": {"2"}}
	request, err := NewQueryRequest(context.Background(), http.MethodGet, "https://api.bilibili.com/x/test?fixed=1", query)
	if err != nil {
		t.Fatalf("newQueryRequest() error = %v", err)
	}
	if request.URL.Query().Get("fixed") != "1" || request.URL.Query().Get("page") != "2" {
		t.Fatalf("query = %v, want fixed and page", request.URL.Query())
	}
	if query.Get("fixed") != "" {
		t.Fatal("newQueryRequest() mutated caller query")
	}

	form := url.Values{"keyword": {"space value"}}
	request, err = NewFormRequest(context.Background(), "https://api.bilibili.com/x/form", query, form)
	if err != nil {
		t.Fatalf("newFormRequest() error = %v", err)
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("read form body: %v", err)
	}
	if request.Header.Get("Content-Type") != "application/x-www-form-urlencoded" || string(body) != "keyword=space+value" {
		t.Fatalf("form request content-type=%q body=%q", request.Header.Get("Content-Type"), body)
	}
}

func TestMultipartPolicyStreamsFieldsAndFiles(t *testing.T) {
	t.Parallel()

	request, err := newMultipartRequest(
		context.Background(),
		"https://api.bilibili.com/x/upload",
		url.Values{"csrf": {"token"}},
		url.Values{"title": {"demo"}},
		[]multipartFile{{
			FieldName:   "cover",
			FileName:    "cover.txt",
			ContentType: "text/plain",
			Content:     strings.NewReader("file-body"),
		}},
	)
	if err != nil {
		t.Fatalf("newMultipartRequest() error = %v", err)
	}
	if request.URL.Query().Get("csrf") != "token" {
		t.Fatalf("query = %v, want csrf", request.URL.Query())
	}
	reader, err := request.MultipartReader()
	if err != nil {
		t.Fatalf("MultipartReader() error = %v", err)
	}
	parts := make(map[string]string)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("NextPart() error = %v", err)
		}
		body, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("read multipart part: %v", err)
		}
		parts[part.FormName()] = string(body)
		if part.FormName() == "cover" && (part.FileName() != "cover.txt" || part.Header.Get("Content-Type") != "text/plain") {
			t.Fatalf("file metadata filename=%q content-type=%q", part.FileName(), part.Header.Get("Content-Type"))
		}
	}
	if parts["title"] != "demo" || parts["cover"] != "file-body" {
		t.Fatalf("multipart parts = %#v", parts)
	}
}

func TestCSRFPolicyClonesValuesAndRequiresAccount(t *testing.T) {
	t.Parallel()

	client, err := NewClient(WithAccount(Account{DedeUserID: "42", SESSDATA: "session", BiliJCT: "secret-csrf", Buvid3: "buvid"}))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	original := url.Values{"key": {"value"}}
	values, err := client.csrfValues(original, "csrf", "csrf_token")
	if err != nil {
		t.Fatalf("csrfValues() error = %v", err)
	}
	if values.Get("csrf") != "secret-csrf" || values.Get("csrf_token") != "secret-csrf" || original.Get("csrf") != "" {
		t.Fatalf("csrf values=%v original=%v", values, original)
	}

	anonymous, _ := NewClient()
	if _, err := anonymous.csrfValues(nil); err == nil {
		t.Fatal("anonymous csrfValues() error = nil, want authentication error")
	}
}

func TestWBIPolicySignsSingleValueQueryWithoutMutation(t *testing.T) {
	t.Parallel()

	const navigationBody = `{"code":-101,"data":{"wbi_img":{"img_url":"https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png","sub_url":"https://i0.hdslb.com/bfs/wbi/ABCDEFGHIJKLMNOPQRSTUVWXYZ654321.png"}}}`
	var calls atomic.Int32
	client, err := NewClient(
		WithClock(func() time.Time { return time.Unix(1_700_000_000, 0) }),
		WithHTTPClient(&http.Client{Transport: testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			calls.Add(1)
			return testutil.JSONResponse(http.StatusOK, navigationBody), nil
		})}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	original := url.Values{"foo": {"value!'()*"}}
	signed, err := client.WBIValues(context.Background(), original)
	if err != nil {
		t.Fatalf("wbiValues() error = %v", err)
	}
	if signed.Get("foo") != "value" || signed.Get("wts") != "1700000000" || len(signed.Get("w_rid")) != 32 {
		t.Fatalf("signed values = %v", signed)
	}
	if original.Get("foo") != "value!'()*" || original.Get("w_rid") != "" {
		t.Fatal("wbiValues() mutated caller values")
	}
	if calls.Load() != 1 {
		t.Fatalf("navigation calls = %d, want 1", calls.Load())
	}
	if _, err := client.WBIValues(context.Background(), url.Values{"duplicate": {"one", "two"}}); err == nil {
		t.Fatal("wbiValues(duplicate) error = nil, want error")
	}
}
