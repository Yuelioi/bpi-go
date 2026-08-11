package probe

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var expectedRiskCounts = map[string]int{
	"public-read":        115,
	"authenticated-read": 32,
	"private-read":       52,
	"login-session":      7,
}

func Audit(sourcePath, implementedPath, contractRoot, lockPath string) (AuditReport, error) {
	var err error
	contractRoot, err = filepath.Abs(contractRoot)
	if err != nil {
		return AuditReport{}, fmt.Errorf("probe: resolve contract root: %w", err)
	}
	var source SourceManifest
	if err := ReadJSON(sourcePath, &source); err != nil {
		return AuditReport{}, err
	}
	var implemented ImplementedManifest
	if err := ReadJSON(implementedPath, &implemented); err != nil {
		return AuditReport{}, err
	}
	var lock ContractLock
	if err := ReadJSON(lockPath, &lock); err != nil {
		return AuditReport{}, err
	}
	if source.SchemaVersion != SchemaVersion || implemented.SchemaVersion != SchemaVersion {
		return AuditReport{}, fmt.Errorf("probe: unsupported parity schema")
	}
	if source.Source.Commit != SourceCommit || implemented.SourceCommit != source.Source.Commit || lock.SourceCommit != source.Source.Commit {
		return AuditReport{}, fmt.Errorf("probe: source, implementation, and contract locks do not reference %s", SourceCommit)
	}
	if err := VerifyContractLock(contractRoot, lock); err != nil {
		return AuditReport{}, err
	}

	sourceContracts := make(map[string]SourceContract, source.Summary.PromotedContracts)
	riskCounts := make(map[string]int)
	domainNames := make(map[string]struct{}, len(source.Domains))
	for _, domain := range source.Domains {
		if domain.Name == "" {
			return AuditReport{}, fmt.Errorf("probe: source manifest contains an empty domain")
		}
		if _, duplicate := domainNames[domain.Name]; duplicate {
			return AuditReport{}, fmt.Errorf("probe: duplicate source domain %q", domain.Name)
		}
		domainNames[domain.Name] = struct{}{}
		for _, contract := range domain.Contracts {
			if _, duplicate := sourceContracts[contract.Name]; duplicate {
				return AuditReport{}, fmt.Errorf("probe: duplicate source contract %q", contract.Name)
			}
			if contract.Status != "promoted" || contract.Name == "" || contract.Risk == "" || contract.SourcePath == "" {
				return AuditReport{}, fmt.Errorf("probe: incomplete source contract %q", contract.Name)
			}
			if _, supported := expectedRiskCounts[contract.Risk]; !supported {
				return AuditReport{}, fmt.Errorf("probe: unsupported risk %q for %s", contract.Risk, contract.Name)
			}
			sourceContracts[contract.Name] = contract
			riskCounts[contract.Risk]++
		}
	}
	if len(domainNames) != source.Summary.DomainClients || len(sourceContracts) != source.Summary.PromotedContracts {
		return AuditReport{}, fmt.Errorf("probe: source summary counts do not match its inventory")
	}
	if len(domainNames) != 27 || len(sourceContracts) != 206 {
		return AuditReport{}, fmt.Errorf("probe: locked surface is %d domains/%d contracts, want 27/206", len(domainNames), len(sourceContracts))
	}
	for risk, want := range expectedRiskCounts {
		if riskCounts[risk] != want || source.Summary.RiskCounts[risk] != want {
			return AuditReport{}, fmt.Errorf("probe: risk %s count drifted from %d", risk, want)
		}
	}

	mappings := make(map[string]Mapping, len(implemented.Mappings))
	for _, mapping := range implemented.Mappings {
		if _, exists := sourceContracts[mapping.Contract]; !exists {
			return AuditReport{}, fmt.Errorf("probe: implementation references unknown contract %q", mapping.Contract)
		}
		if _, duplicate := mappings[mapping.Contract]; duplicate {
			return AuditReport{}, fmt.Errorf("probe: duplicate implementation mapping %q", mapping.Contract)
		}
		if mapping.DomainAccessor == "" || mapping.Method == "" || mapping.Params == "" || mapping.Model == "" {
			return AuditReport{}, fmt.Errorf("probe: implementation mapping %q has an empty Go symbol", mapping.Contract)
		}
		mappings[mapping.Contract] = mapping
	}
	if len(mappings) != len(sourceContracts) {
		return AuditReport{}, fmt.Errorf("probe: implemented mappings cover %d of %d contracts", len(mappings), len(sourceContracts))
	}

	contractFiles, err := findContractFiles(contractRoot)
	if err != nil {
		return AuditReport{}, err
	}
	if len(contractFiles) != len(sourceContracts) {
		return AuditReport{}, fmt.Errorf("probe: snapshot contains %d contract files, source contains %d", len(contractFiles), len(sourceContracts))
	}
	fixturePaths := make(map[string]struct{})
	for name, sourceContract := range sourceContracts {
		relative, ok := strings.CutPrefix(filepath.ToSlash(sourceContract.SourcePath), "tests/contracts/")
		if !ok {
			return AuditReport{}, fmt.Errorf("probe: source path for %s is outside tests/contracts", name)
		}
		path, err := joinWithin(contractRoot, filepath.FromSlash(relative))
		if err != nil {
			return AuditReport{}, err
		}
		var contract Contract
		if err := ReadJSON(path, &contract); err != nil {
			return AuditReport{}, err
		}
		if err := validateContract(path, contract, sourceContract, fixturePaths, contractRoot); err != nil {
			return AuditReport{}, err
		}
	}
	findings, err := AuditSanitization(contractRoot)
	if err != nil {
		return AuditReport{}, err
	}
	if len(findings) != 0 {
		first := findings[0]
		return AuditReport{}, fmt.Errorf("probe: sanitize audit found %d issue(s); first is %s %s: %s", len(findings), first.File, first.Path, first.Reason)
	}
	return AuditReport{
		Contracts:  len(sourceContracts),
		Fixtures:   len(fixturePaths),
		Mappings:   len(mappings),
		Files:      len(lock.Files),
		RiskCounts: riskCounts,
	}, nil
}

func validateContract(path string, contract Contract, source SourceContract, fixturePaths map[string]struct{}, contractRoot string) error {
	if contract.SchemaVersion != 2 || contract.Status != "promoted" {
		return fmt.Errorf("probe: %s is not a promoted schema-v2 contract", path)
	}
	if contract.Name != source.Name || contract.Risk != source.Risk {
		return fmt.Errorf("probe: contract identity mismatch in %s", path)
	}
	if contract.Module == "" || contract.Batch == "" || contract.Endpoint == "" || len(contract.Profiles) == 0 || contract.Sanitize == nil {
		return fmt.Errorf("probe: contract %s is missing required metadata", contract.Name)
	}
	if len(contract.Steps) == 0 {
		if contract.Request.Method != source.Method || contract.Request.URL != source.URL {
			return fmt.Errorf("probe: request drift for %s", contract.Name)
		}
		if contract.Request.Method != "GET" && contract.Request.Method != "POST" {
			return fmt.Errorf("probe: unsupported method %q for %s", contract.Request.Method, contract.Name)
		}
		if err := validateProbeURL(contract.Request.URL); err != nil {
			return fmt.Errorf("probe: %s: %w", contract.Name, err)
		}
	} else if contract.Risk != "login-session" || contract.Name != "login.qr.flow" {
		return fmt.Errorf("probe: unexpected multi-step contract %s", contract.Name)
	}
	profiles := make(map[string]struct{}, len(contract.Profiles))
	for _, profile := range contract.Profiles {
		if profile != "anonymous" && profile != "normal" && profile != "vip" {
			return fmt.Errorf("probe: unsupported profile %q for %s", profile, contract.Name)
		}
		if _, duplicate := profiles[profile]; duplicate {
			return fmt.Errorf("probe: duplicate profile %q for %s", profile, contract.Name)
		}
		profiles[profile] = struct{}{}
	}
	caseProfiles := make(map[string]struct{}, len(contract.Cases))
	for _, contractCase := range contract.Cases {
		if contractCase.Name == "" || contractCase.Profile == "" {
			return fmt.Errorf("probe: contract %s has an unnamed case", contract.Name)
		}
		if _, declared := profiles[contractCase.Profile]; !declared {
			return fmt.Errorf("probe: case profile %q is not declared by %s", contractCase.Profile, contract.Name)
		}
		if _, duplicate := caseProfiles[contractCase.Profile]; duplicate {
			return fmt.Errorf("probe: contract %s has duplicate case profile %q", contract.Name, contractCase.Profile)
		}
		caseProfiles[contractCase.Profile] = struct{}{}
		if contractCase.Response.Fixture == "" {
			if contractCase.Response.FixtureKind != "local_probe_blocked" {
				return fmt.Errorf("probe: case %s.%s has no fixture", contract.Name, contractCase.Name)
			}
			continue
		}
		fixturePath, err := joinWithin(filepath.Dir(path), filepath.FromSlash(contractCase.Response.Fixture))
		if err != nil {
			return fmt.Errorf("probe: %s fixture: %w", contract.Name, err)
		}
		if _, err := os.Stat(fixturePath); err != nil {
			return fmt.Errorf("probe: missing fixture for %s: %w", contract.Name, err)
		}
		relative, err := filepath.Rel(contractRoot, fixturePath)
		if err != nil {
			return err
		}
		fixturePaths[filepath.ToSlash(relative)] = struct{}{}
		if err := validateFixture(fixturePath, contractCase.Response); err != nil {
			return fmt.Errorf("probe: %s.%s: %w", contract.Name, contractCase.Name, err)
		}
	}
	if len(contract.Steps) == 0 && len(caseProfiles) != len(profiles) {
		return fmt.Errorf("probe: contract %s has %d profiles but %d cases", contract.Name, len(profiles), len(caseProfiles))
	}
	return nil
}

func validateFixture(path string, expected ContractResponse) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var value struct {
		Code       json.RawMessage `json:"code"`
		Errno      json.RawMessage `json:"errno"`
		Kind       string          `json:"kind"`
		BodyBase64 string          `json:"body_base64"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("invalid JSON fixture %s: %w", path, err)
	}
	if expected.APICode == nil {
		return nil
	}
	actual := fixtureInteger(value.Code)
	if actual == nil {
		actual = fixtureInteger(value.Errno)
	}
	if actual == nil && value.Kind == "binary" && value.BodyBase64 != "" {
		zero := 0
		actual = &zero
	}
	if actual == nil || *actual != *expected.APICode {
		return fmt.Errorf("fixture API code does not match expected %d", *expected.APICode)
	}
	return nil
}

func fixtureInteger(raw json.RawMessage) *int {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var value int
	if json.Unmarshal(raw, &value) != nil {
		return nil
	}
	return &value
}

func findContractFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && entry.Name() == "contract.json" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("probe: enumerate contracts: %w", err)
	}
	sort.Strings(paths)
	return paths, nil
}

func joinWithin(root, relative string) (string, error) {
	if filepath.IsAbs(relative) {
		return "", fmt.Errorf("absolute snapshot path is not allowed")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	target, err := filepath.Abs(filepath.Join(root, filepath.Clean(relative)))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("snapshot path escapes its root")
	}
	return target, nil
}

func validateProbeURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil {
		return fmt.Errorf("endpoint must be an absolute credential-free HTTPS URL")
	}
	host := strings.ToLower(parsed.Hostname())
	if host != "bilibili.com" && !strings.HasSuffix(host, ".bilibili.com") && host != "api.biliapi.net" {
		return fmt.Errorf("endpoint host %q is outside the approved Bilibili set", host)
	}
	return nil
}
