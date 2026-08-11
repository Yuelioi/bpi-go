package probe

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/internal/sign"
)

type BatchConfig struct {
	ContractRoot  string
	CookieHeaders map[string]string
	Profiles      []string
	ReadOnly      bool
	HTTPClient    *http.Client
	Now           func() time.Time
}

type BatchSummary struct {
	SchemaVersion int           `json:"schema_version"`
	SourceCommit  string        `json:"source_commit"`
	StartedAt     time.Time     `json:"started_at"`
	FinishedAt    time.Time     `json:"finished_at"`
	ReadOnly      bool          `json:"read_only"`
	Profiles      []string      `json:"profiles"`
	Total         int           `json:"total"`
	Passed        int           `json:"passed"`
	Failed        int           `json:"failed"`
	Skipped       int           `json:"skipped"`
	Results       []ProbeResult `json:"results"`
}

type ProbeResult struct {
	Contract   string `json:"contract"`
	Profile    string `json:"profile"`
	Risk       string `json:"risk"`
	HTTPStatus int    `json:"http_status,omitempty"`
	APICode    *int   `json:"api_code,omitempty"`
	Bytes      int    `json:"bytes,omitempty"`
	SHA256     string `json:"sha256,omitempty"`
	Outcome    string `json:"outcome"`
}

type BatchError struct{ Failed int }

func (e *BatchError) Error() string { return fmt.Sprintf("probe: %d live probe(s) failed", e.Failed) }

type batchClient struct {
	client *bpi.Client
	keys   *sign.WBIKeys
	now    func() time.Time
}

func RunBatch(ctx context.Context, config BatchConfig) (BatchSummary, error) {
	if ctx == nil {
		return BatchSummary{}, fmt.Errorf("probe: context cannot be nil")
	}
	if os.Getenv("BPI_PROBE") != "1" {
		return BatchSummary{}, fmt.Errorf("probe: set BPI_PROBE=1 to enable network probes")
	}
	if !config.ReadOnly {
		return BatchSummary{}, fmt.Errorf("probe: only the explicit read-only gate is supported")
	}
	profiles, err := normalizeProfiles(config.Profiles)
	if err != nil {
		return BatchSummary{}, err
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	started := now().UTC()
	summary := BatchSummary{
		SchemaVersion: SchemaVersion,
		SourceCommit:  SourceCommit,
		StartedAt:     started,
		ReadOnly:      true,
		Profiles:      profiles,
	}
	clients, err := buildBatchClients(config, profiles, now)
	if err != nil {
		return summary, err
	}
	paths, err := findContractFiles(config.ContractRoot)
	if err != nil {
		return summary, err
	}
	for _, path := range paths {
		var contract Contract
		if err := ReadJSON(path, &contract); err != nil {
			return summary, err
		}
		if !isReadOnlyRisk(contract.Risk) {
			summary.Skipped++
			continue
		}
		for _, profile := range profiles {
			contractCase, ok := caseForProfile(contract.Cases, profile)
			if !ok {
				continue
			}
			summary.Total++
			result := executeProbe(ctx, clients[profile], contract, contractCase)
			summary.Results = append(summary.Results, result)
			if result.Outcome == "passed" {
				summary.Passed++
			} else {
				summary.Failed++
			}
		}
	}
	summary.FinishedAt = now().UTC()
	if summary.Failed != 0 {
		return summary, &BatchError{Failed: summary.Failed}
	}
	return summary, nil
}

func buildBatchClients(config BatchConfig, profiles []string, now func() time.Time) (map[string]*batchClient, error) {
	clients := make(map[string]*batchClient, len(profiles))
	for _, profile := range profiles {
		options := make([]bpi.Option, 0, 2)
		if config.HTTPClient != nil {
			options = append(options, bpi.WithHTTPClient(config.HTTPClient))
		}
		if profile != "anonymous" {
			cookieHeader := config.CookieHeaders[profile]
			if strings.TrimSpace(cookieHeader) == "" {
				return nil, fmt.Errorf("probe: Cookie header is required for profile %s", profile)
			}
			options = append(options, bpi.WithCookie(cookieHeader))
		}
		client, err := bpi.NewClient(options...)
		if err != nil {
			return nil, err
		}
		clients[profile] = &batchClient{client: client, now: now}
	}
	return clients, nil
}

func executeProbe(ctx context.Context, client *batchClient, contract Contract, contractCase ContractCase) ProbeResult {
	result := ProbeResult{Contract: contract.Name, Profile: contractCase.Profile, Risk: contract.Risk}
	request, err := buildContractRequest(ctx, client, contract.Request)
	if err != nil {
		result.Outcome = "request_error"
		return result
	}
	response, err := client.client.Do(ctx, request, "probe."+contract.Name+"."+contractCase.Profile)
	if err != nil {
		var httpError *bpi.HTTPError
		if errors.As(err, &httpError) {
			result.HTTPStatus = httpError.StatusCode
			if contractCase.Response.HTTPStatus != nil && httpError.StatusCode == *contractCase.Response.HTTPStatus {
				result.Outcome = "passed"
				return result
			}
			result.Outcome = "http_mismatch"
			return result
		}
		result.Outcome = "transport_error"
		return result
	}
	result.HTTPStatus = response.StatusCode
	result.Bytes = len(response.Body)
	digest := sha256.Sum256(response.Body)
	result.SHA256 = hex.EncodeToString(digest[:])
	result.APICode = responseAPICode(response.Body)
	wantStatus := http.StatusOK
	if contractCase.Response.HTTPStatus != nil {
		wantStatus = *contractCase.Response.HTTPStatus
	}
	if result.HTTPStatus != wantStatus {
		result.Outcome = "http_mismatch"
		return result
	}
	if contractCase.Response.APICode != nil {
		actualCode := result.APICode
		if actualCode == nil && !json.Valid(response.Body) {
			zero := 0
			actualCode = &zero
			result.APICode = actualCode
		}
		if actualCode == nil || *actualCode != *contractCase.Response.APICode {
			result.Outcome = "api_mismatch"
			return result
		}
	}
	result.Outcome = "passed"
	return result
}

func buildContractRequest(ctx context.Context, client *batchClient, contract ContractRequest) (*http.Request, error) {
	if err := validateProbeURL(contract.URL); err != nil {
		return nil, err
	}
	variables := map[string]string{"csrf": ""}
	if csrf, err := client.client.CSRF(); err == nil {
		variables["csrf"] = csrf
	}
	query, err := resolveValues(contract.Query, variables)
	if err != nil {
		return nil, err
	}
	if requires(contract.Auth, "wbi") {
		query, err = client.signWBI(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	var body io.Reader
	contentType := ""
	if len(contract.Form) != 0 {
		form, err := resolveValues(contract.Form, variables)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(form.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else if len(bytes.TrimSpace(contract.Body)) != 0 && !bytes.Equal(bytes.TrimSpace(contract.Body), []byte("null")) {
		resolved, err := resolveJSON(contract.Body, variables)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(resolved)
		contentType = "application/json"
	}
	request, err := http.NewRequestWithContext(ctx, contract.Method, contract.URL, body)
	if err != nil {
		return nil, fmt.Errorf("probe: create request: %w", err)
	}
	request.URL.RawQuery = query.Encode()
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	for key, value := range contract.Headers {
		resolved, err := resolveTemplate(value, variables)
		if err != nil {
			return nil, err
		}
		if strings.ContainsAny(key+resolved, "\r\n") {
			return nil, fmt.Errorf("probe: contract header contains a newline")
		}
		request.Header.Set(key, resolved)
	}
	return request, nil
}

func (c *batchClient) signWBI(ctx context.Context, values url.Values) (url.Values, error) {
	if c.keys == nil {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.bilibili.com/x/web-interface/nav", nil)
		if err != nil {
			return nil, err
		}
		response, err := c.client.Do(ctx, request, "probe.sign.wbi.navigation")
		if err != nil {
			return nil, err
		}
		envelope, err := bpi.DecodeEnvelope[struct {
			WBIImage struct {
				ImageURL string `json:"img_url"`
				SubURL   string `json:"sub_url"`
			} `json:"wbi_img"`
		}](response.Body)
		if err != nil {
			return nil, err
		}
		payload, err := envelope.IntoData()
		if err != nil {
			return nil, err
		}
		keys, err := sign.WBIKeysFromURLs(payload.WBIImage.ImageURL, payload.WBIImage.SubURL)
		if err != nil {
			return nil, err
		}
		c.keys = &keys
	}
	params := make(map[string]string, len(values))
	for key, candidates := range values {
		if len(candidates) != 1 {
			return nil, fmt.Errorf("probe: WBI query %q must have exactly one value", key)
		}
		params[key] = candidates[0]
	}
	signed, err := sign.SignWBIAt(params, *c.keys, uint64(c.now().Unix()))
	if err != nil {
		return nil, err
	}
	result := make(url.Values, len(signed))
	for key, value := range signed {
		result.Set(key, value)
	}
	return result, nil
}

func resolveValues(input map[string]string, variables map[string]string) (url.Values, error) {
	result := make(url.Values, len(input))
	for key, value := range input {
		resolved, err := resolveTemplate(value, variables)
		if err != nil {
			return nil, err
		}
		result.Set(key, resolved)
	}
	return result, nil
}

func resolveJSON(input json.RawMessage, variables map[string]string) ([]byte, error) {
	var value any
	if err := json.Unmarshal(input, &value); err != nil {
		return nil, fmt.Errorf("probe: decode request body: %w", err)
	}
	resolved, err := resolveJSONValue(value, variables)
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(resolved)
	if err != nil {
		return nil, fmt.Errorf("probe: encode request body: %w", err)
	}
	return result, nil
}

func resolveJSONValue(value any, variables map[string]string) (any, error) {
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			resolved, err := resolveJSONValue(child, variables)
			if err != nil {
				return nil, err
			}
			current[key] = resolved
		}
		return current, nil
	case []any:
		for index, child := range current {
			resolved, err := resolveJSONValue(child, variables)
			if err != nil {
				return nil, err
			}
			current[index] = resolved
		}
		return current, nil
	case string:
		return resolveTemplate(current, variables)
	default:
		return value, nil
	}
}

func resolveTemplate(value string, variables map[string]string) (string, error) {
	for key, replacement := range variables {
		value = strings.ReplaceAll(value, "${"+key+"}", replacement)
	}
	if strings.Contains(value, "${") {
		return "", fmt.Errorf("probe: unresolved contract template in %q", value)
	}
	return value, nil
}

func responseAPICode(body []byte) *int {
	var envelope struct {
		Code  *int `json:"code"`
		Errno *int `json:"errno"`
	}
	if json.Unmarshal(body, &envelope) != nil {
		return nil
	}
	if envelope.Code != nil {
		return envelope.Code
	}
	return envelope.Errno
}

func requires(auth ContractAuth, requirement string) bool {
	for _, candidate := range auth.Requires {
		if candidate == requirement {
			return true
		}
	}
	return false
}

func caseForProfile(cases []ContractCase, profile string) (ContractCase, bool) {
	for _, candidate := range cases {
		if candidate.Profile == profile {
			return candidate, true
		}
	}
	return ContractCase{}, false
}

func normalizeProfiles(profiles []string) ([]string, error) {
	if len(profiles) == 0 {
		profiles = []string{"anonymous"}
	}
	seen := make(map[string]struct{}, len(profiles))
	result := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		profile = strings.ToLower(strings.TrimSpace(profile))
		if profile != "anonymous" && profile != "normal" && profile != "vip" {
			return nil, fmt.Errorf("probe: unsupported profile %q", profile)
		}
		if _, duplicate := seen[profile]; duplicate {
			continue
		}
		seen[profile] = struct{}{}
		result = append(result, profile)
	}
	sort.Slice(result, func(i, j int) bool {
		order := map[string]int{"anonymous": 0, "normal": 1, "vip": 2}
		return order[result[i]] < order[result[j]]
	})
	return result, nil
}

func ParseProfiles(value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("probe: profiles cannot be empty")
	}
	return normalizeProfiles(strings.Split(value, ","))
}

func isReadOnlyRisk(risk string) bool {
	return risk == "public-read" || risk == "authenticated-read" || risk == "private-read"
}
