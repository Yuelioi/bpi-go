package bpi_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSourceManifestMatchesLockedPromotedSurface(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("parity/source.json")
	if err != nil {
		t.Fatalf("read parity/source.json: %v", err)
	}
	var manifest struct {
		SchemaVersion int `json:"schema_version"`
		Source        struct {
			Commit string `json:"commit"`
		} `json:"source"`
		Summary struct {
			DomainClients     int            `json:"domain_clients"`
			PromotedContracts int            `json:"promoted_contracts"`
			RiskCounts        map[string]int `json:"risk_counts"`
		} `json:"summary"`
		Domains []struct {
			Name      string `json:"name"`
			Contracts []struct {
				Name string `json:"name"`
			} `json:"contracts"`
		} `json:"domains"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("decode parity/source.json: %v", err)
	}

	if manifest.SchemaVersion != 1 {
		t.Fatalf("schema_version = %d, want 1", manifest.SchemaVersion)
	}
	if manifest.Source.Commit != "36cb1104befee33b4281c59a4365d42dfdde45a4" {
		t.Fatalf("source commit = %q, want locked bpi-rs commit", manifest.Source.Commit)
	}
	if manifest.Summary.DomainClients != 27 || len(manifest.Domains) != 27 {
		t.Fatalf("domain clients = %d/%d, want 27", manifest.Summary.DomainClients, len(manifest.Domains))
	}
	if manifest.Summary.PromotedContracts != 206 {
		t.Fatalf("promoted contracts = %d, want 206", manifest.Summary.PromotedContracts)
	}
	wantRiskCounts := map[string]int{
		"public-read":        115,
		"authenticated-read": 32,
		"private-read":       52,
		"login-session":      7,
	}
	for risk, want := range wantRiskCounts {
		if got := manifest.Summary.RiskCounts[risk]; got != want {
			t.Fatalf("risk count %s = %d, want %d", risk, got, want)
		}
	}

	seenDomains := make(map[string]struct{}, len(manifest.Domains))
	seenContracts := make(map[string]struct{}, manifest.Summary.PromotedContracts)
	for _, domain := range manifest.Domains {
		if _, exists := seenDomains[domain.Name]; exists {
			t.Fatalf("duplicate domain %q", domain.Name)
		}
		seenDomains[domain.Name] = struct{}{}
		for _, contract := range domain.Contracts {
			if _, exists := seenContracts[contract.Name]; exists {
				t.Fatalf("duplicate contract %q", contract.Name)
			}
			seenContracts[contract.Name] = struct{}{}
		}
	}
	if len(seenContracts) != 206 {
		t.Fatalf("unique contracts = %d, want 206", len(seenContracts))
	}
}

func TestImplementedMappingsReferenceLockedContractsAndCopiedFixtures(t *testing.T) {
	t.Parallel()

	sourceData, err := os.ReadFile("parity/source.json")
	if err != nil {
		t.Fatalf("read parity/source.json: %v", err)
	}
	var source struct {
		Source struct {
			Commit string `json:"commit"`
		} `json:"source"`
		Domains []struct {
			Contracts []struct {
				Name       string `json:"name"`
				Risk       string `json:"risk"`
				SourcePath string `json:"source_path"`
			} `json:"contracts"`
		} `json:"domains"`
	}
	if err := json.Unmarshal(sourceData, &source); err != nil {
		t.Fatalf("decode parity/source.json: %v", err)
	}
	type sourceContract struct {
		path string
		risk string
	}
	sourceContracts := make(map[string]sourceContract)
	for _, domain := range source.Domains {
		for _, contract := range domain.Contracts {
			sourceContracts[contract.Name] = sourceContract{path: contract.SourcePath, risk: contract.Risk}
		}
	}

	implementedData, err := os.ReadFile("parity/implemented.json")
	if err != nil {
		t.Fatalf("read parity/implemented.json: %v", err)
	}
	var implemented struct {
		SchemaVersion int    `json:"schema_version"`
		SourceCommit  string `json:"source_commit"`
		Mappings      []struct {
			Contract       string `json:"contract"`
			DomainAccessor string `json:"domain_accessor"`
			Method         string `json:"method"`
			Params         string `json:"params"`
			Model          string `json:"model"`
		} `json:"mappings"`
	}
	if err := json.Unmarshal(implementedData, &implemented); err != nil {
		t.Fatalf("decode parity/implemented.json: %v", err)
	}
	if implemented.SchemaVersion != 1 || implemented.SourceCommit != source.Source.Commit {
		t.Fatalf("implemented manifest lock = schema %d commit %q; want schema 1 commit %q", implemented.SchemaVersion, implemented.SourceCommit, source.Source.Commit)
	}

	seen := make(map[string]struct{}, len(implemented.Mappings))
	for _, mapping := range implemented.Mappings {
		if _, duplicate := seen[mapping.Contract]; duplicate {
			t.Fatalf("duplicate implemented mapping %q", mapping.Contract)
		}
		seen[mapping.Contract] = struct{}{}
		sourceContract, exists := sourceContracts[mapping.Contract]
		if !exists {
			t.Fatalf("implemented mapping %q is absent from locked source", mapping.Contract)
		}
		if mapping.DomainAccessor == "" || mapping.Method == "" || mapping.Params == "" || mapping.Model == "" {
			t.Fatalf("implemented mapping %q has an empty Go symbol", mapping.Contract)
		}
		relative := strings.TrimPrefix(filepath.ToSlash(sourceContract.path), "tests/contracts/")
		contractData, err := os.ReadFile(filepath.Join("testdata", "contracts", filepath.FromSlash(relative)))
		if err != nil {
			t.Fatalf("read copied contract for %q: %v", mapping.Contract, err)
		}
		var contract struct {
			Name   string `json:"name"`
			Risk   string `json:"risk"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal(contractData, &contract); err != nil {
			t.Fatalf("decode copied contract for %q: %v", mapping.Contract, err)
		}
		if contract.Name != mapping.Contract || contract.Status != "promoted" {
			t.Fatalf("copied contract = %q/%q, want %q/promoted", contract.Name, contract.Status, mapping.Contract)
		}
		if contract.Risk != sourceContract.risk {
			t.Fatalf("copied contract %q risk = %q, want locked risk %q", mapping.Contract, contract.Risk, sourceContract.risk)
		}
	}
}
