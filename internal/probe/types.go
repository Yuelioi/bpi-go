// Package probe implements the offline audits, catalog generation, and
// explicitly gated read-only network runner used by cmd/bpi-probe.
package probe

import "encoding/json"

const (
	SchemaVersion = 1
	SourceCommit  = "94bcf43e46d6e11b55ec4d36d4848692cecd4213"
)

type SourceManifest struct {
	SchemaVersion int `json:"schema_version"`
	Source        struct {
		Repository   string `json:"repository"`
		Commit       string `json:"commit"`
		CargoVersion string `json:"cargo_version"`
	} `json:"source"`
	Summary struct {
		DomainClients     int            `json:"domain_clients"`
		PromotedContracts int            `json:"promoted_contracts"`
		RiskCounts        map[string]int `json:"risk_counts"`
	} `json:"summary"`
	Domains []SourceDomain `json:"domains"`
}

type SourceDomain struct {
	Name          string           `json:"name"`
	SourceMethods []string         `json:"source_methods"`
	Contracts     []SourceContract `json:"contracts"`
}

type SourceContract struct {
	Name       string `json:"name"`
	Risk       string `json:"risk"`
	Status     string `json:"status"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	SourcePath string `json:"source_path"`
}

type ImplementedManifest struct {
	SchemaVersion int       `json:"schema_version"`
	SourceCommit  string    `json:"source_commit"`
	Mappings      []Mapping `json:"mappings"`
}

type Mapping struct {
	Contract       string `json:"contract"`
	DomainAccessor string `json:"domain_accessor"`
	Method         string `json:"method"`
	Params         string `json:"params"`
	Model          string `json:"model"`
}

type Contract struct {
	SchemaVersion int               `json:"schema_version"`
	Name          string            `json:"name"`
	Module        string            `json:"module"`
	Batch         string            `json:"batch"`
	Endpoint      string            `json:"endpoint"`
	Risk          string            `json:"risk"`
	Status        string            `json:"status"`
	Profiles      []string          `json:"profiles"`
	Request       ContractRequest   `json:"request"`
	Cases         []ContractCase    `json:"cases"`
	Steps         []json.RawMessage `json:"steps"`
	Sanitize      *SanitizeSpec     `json:"sanitize"`
}

type ContractRequest struct {
	Method           string            `json:"method"`
	URL              string            `json:"url"`
	Query            map[string]string `json:"query"`
	Form             map[string]string `json:"form"`
	Body             json.RawMessage   `json:"body"`
	RequiredHeaders  []string          `json:"required_headers"`
	Headers          map[string]string `json:"headers"`
	Auth             ContractAuth      `json:"auth"`
	ResponseDecoding json.RawMessage   `json:"response_decoding"`
}

type ContractAuth struct {
	Profile  string   `json:"profile"`
	Requires []string `json:"requires"`
}

type ContractCase struct {
	Name     string           `json:"name"`
	Profile  string           `json:"profile"`
	Auth     ContractAuth     `json:"auth"`
	Response ContractResponse `json:"response"`
}

type ContractResponse struct {
	APICode     *int   `json:"api_code"`
	APICodeText string `json:"api_code_text"`
	HTTPStatus  *int   `json:"http_status"`
	Fixture     string `json:"fixture"`
	FixtureKind string `json:"fixture_kind"`
	Error       string `json:"error"`
	RustModel   string `json:"rust_model"`
}

type SanitizeSpec struct {
	Preset  []string                   `json:"preset"`
	Replace map[string]json.RawMessage `json:"replace"`
	Drop    []string                   `json:"drop"`
	Keep    []string                   `json:"keep"`
}

type ContractLock struct {
	SchemaVersion int          `json:"schema_version"`
	SourceCommit  string       `json:"source_commit"`
	Files         []LockedFile `json:"files"`
}

type LockedFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
}

type AuditReport struct {
	Contracts  int            `json:"contracts"`
	Fixtures   int            `json:"fixtures"`
	Mappings   int            `json:"mappings"`
	Files      int            `json:"files"`
	RiskCounts map[string]int `json:"risk_counts"`
}
