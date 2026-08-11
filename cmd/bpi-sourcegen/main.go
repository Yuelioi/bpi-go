// Command bpi-sourcegen creates a deterministic inventory of the bpi-rs
// domain-client and promoted-contract surface used by the Go port.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const schemaVersion = 1

var (
	methodPattern  = regexp.MustCompile(`(?m)^\s*pub async fn\s+([A-Za-z0-9_]+)\s*\(`)
	versionPattern = regexp.MustCompile(`(?m)^version\s*=\s*"([^"]+)"`)
)

type options struct {
	rustRoot        string
	manifestPath    string
	lockPath        string
	contractsPath   string
	expectCommit    string
	expectDomains   int
	expectContracts int
	check           bool
	syncContracts   bool
}

type sourceInfo struct {
	Repository   string `json:"repository"`
	Commit       string `json:"commit"`
	CargoVersion string `json:"cargo_version"`
}

type summary struct {
	DomainClients     int            `json:"domain_clients"`
	PromotedContracts int            `json:"promoted_contracts"`
	RiskCounts        map[string]int `json:"risk_counts"`
}

type contractRef struct {
	Name       string `json:"name"`
	Risk       string `json:"risk"`
	Status     string `json:"status"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	SourcePath string `json:"source_path"`
}

type domain struct {
	Name          string        `json:"name"`
	SourceMethods []string      `json:"source_methods"`
	Contracts     []contractRef `json:"contracts"`
}

type manifest struct {
	SchemaVersion int        `json:"schema_version"`
	Source        sourceInfo `json:"source"`
	Summary       summary    `json:"summary"`
	Domains       []domain   `json:"domains"`
}

type sourceLock struct {
	SchemaVersion int    `json:"schema_version"`
	Repository    string `json:"repository"`
	Commit        string `json:"commit"`
	CargoVersion  string `json:"cargo_version"`
}

type contractDocument struct {
	Name    string `json:"name"`
	Module  string `json:"module"`
	Risk    string `json:"risk"`
	Status  string `json:"status"`
	Request struct {
		Method string `json:"method"`
		URL    string `json:"url"`
	} `json:"request"`
}

func main() {
	var opts options
	flag.StringVar(&opts.rustRoot, "rust", "../bpi-rs", "path to the bpi-rs checkout")
	flag.StringVar(&opts.manifestPath, "out", "parity/source.json", "manifest output path")
	flag.StringVar(&opts.lockPath, "lock", "parity/source.lock.json", "source lock output path")
	flag.StringVar(&opts.contractsPath, "contracts", "testdata/contracts", "committed contract snapshot path")
	flag.StringVar(&opts.expectCommit, "expect-commit", "", "required bpi-rs commit")
	flag.IntVar(&opts.expectDomains, "expect-domains", 0, "required domain-client count")
	flag.IntVar(&opts.expectContracts, "expect-contracts", 0, "required promoted-contract count")
	flag.BoolVar(&opts.check, "check", false, "verify generated files and contract snapshot without writing")
	flag.BoolVar(&opts.syncContracts, "sync-contracts", false, "copy changed or missing contract snapshot files (never removes stale files)")
	flag.Parse()

	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "bpi-sourcegen:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	if opts.check && opts.syncContracts {
		return errors.New("-check and -sync-contracts cannot be used together")
	}
	rustRoot, err := filepath.Abs(opts.rustRoot)
	if err != nil {
		return fmt.Errorf("resolve Rust root: %w", err)
	}
	if err := ensureRustSourceClean(rustRoot); err != nil {
		return err
	}

	commit, err := gitOutput(rustRoot, "rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("read Rust commit: %w", err)
	}
	if opts.expectCommit != "" && commit != opts.expectCommit {
		return fmt.Errorf("rust commit is %s, expected %s", commit, opts.expectCommit)
	}

	cargoVersion, repository, err := cargoMetadata(filepath.Join(rustRoot, "Cargo.toml"))
	if err != nil {
		return err
	}
	domains, err := readDomains(rustRoot)
	if err != nil {
		return err
	}
	riskCounts, contractCount, err := readContracts(rustRoot, domains)
	if err != nil {
		return err
	}

	if opts.expectDomains != 0 && len(domains) != opts.expectDomains {
		return fmt.Errorf("found %d domain clients, expected %d", len(domains), opts.expectDomains)
	}
	if opts.expectContracts != 0 && contractCount != opts.expectContracts {
		return fmt.Errorf("found %d promoted contracts, expected %d", contractCount, opts.expectContracts)
	}

	source := sourceInfo{
		Repository:   repository,
		Commit:       commit,
		CargoVersion: cargoVersion,
	}
	result := manifest{
		SchemaVersion: schemaVersion,
		Source:        source,
		Summary: summary{
			DomainClients:     len(domains),
			PromotedContracts: contractCount,
			RiskCounts:        riskCounts,
		},
		Domains: domains,
	}
	lock := sourceLock{
		SchemaVersion: schemaVersion,
		Repository:    source.Repository,
		Commit:        source.Commit,
		CargoVersion:  source.CargoVersion,
	}

	if opts.check {
		if err := checkJSON(opts.manifestPath, result); err != nil {
			return fmt.Errorf("check manifest: %w", err)
		}
		if err := checkJSON(opts.lockPath, lock); err != nil {
			return fmt.Errorf("check lock: %w", err)
		}
		if opts.contractsPath != "" {
			if err := compareContractTrees(filepath.Join(rustRoot, "tests", "contracts"), opts.contractsPath); err != nil {
				return err
			}
		}
		fmt.Printf("checked %s, %s, and %s: %d domains, %d promoted contracts\n",
			opts.manifestPath, opts.lockPath, opts.contractsPath, len(domains), contractCount)
		return nil
	}

	if err := writeJSON(opts.manifestPath, result); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	if err := writeJSON(opts.lockPath, lock); err != nil {
		return fmt.Errorf("write lock: %w", err)
	}
	if opts.syncContracts {
		if strings.TrimSpace(opts.contractsPath) == "" {
			return errors.New("-sync-contracts requires a non-empty -contracts path")
		}
		if err := syncContractTrees(filepath.Join(rustRoot, "tests", "contracts"), opts.contractsPath); err != nil {
			return err
		}
	}

	fmt.Printf("wrote %s and %s: %d domains, %d promoted contracts\n",
		opts.manifestPath, opts.lockPath, len(domains), contractCount)
	return nil
}

func ensureRustSourceClean(root string) error {
	output, err := gitOutput(root, "status", "--porcelain", "--", "Cargo.toml", "src", "tests/contracts", "docs/api-index.md")
	if err != nil {
		return fmt.Errorf("inspect Rust source status: %w", err)
	}
	if output != "" {
		return fmt.Errorf("rust source paths have uncommitted changes:\n%s", output)
	}
	return nil
}

func cargoMetadata(path string) (version string, repository string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("read Cargo.toml: %w", err)
	}
	versionMatch := versionPattern.FindSubmatch(data)
	if len(versionMatch) != 2 {
		return "", "", errors.New("cargo.toml package version was not found")
	}
	repositoryPattern := regexp.MustCompile(`(?m)^repository\s*=\s*"([^"]+)"`)
	repositoryMatch := repositoryPattern.FindSubmatch(data)
	if len(repositoryMatch) != 2 {
		return "", "", errors.New("cargo.toml repository was not found")
	}
	return string(versionMatch[1]), string(repositoryMatch[1]), nil
}

func readDomains(root string) ([]domain, error) {
	entries, err := os.ReadDir(filepath.Join(root, "src"))
	if err != nil {
		return nil, fmt.Errorf("read Rust src: %w", err)
	}

	var domains []domain
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		clientPath := filepath.Join(root, "src", entry.Name(), "client.rs")
		data, err := os.ReadFile(clientPath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", clientPath, err)
		}

		matches := methodPattern.FindAllSubmatch(data, -1)
		methods := make([]string, 0, len(matches))
		seen := make(map[string]struct{}, len(matches))
		for _, match := range matches {
			method := string(match[1])
			if _, exists := seen[method]; exists {
				return nil, fmt.Errorf("duplicate public async method %s.%s", entry.Name(), method)
			}
			seen[method] = struct{}{}
			methods = append(methods, method)
		}
		sort.Strings(methods)
		domains = append(domains, domain{Name: entry.Name(), SourceMethods: methods})
	}
	sort.Slice(domains, func(i, j int) bool { return domains[i].Name < domains[j].Name })
	return domains, nil
}

func readContracts(root string, domains []domain) (map[string]int, int, error) {
	byName := make(map[string]*domain, len(domains))
	for i := range domains {
		byName[domains[i].Name] = &domains[i]
	}

	riskCounts := make(map[string]int)
	seenContracts := make(map[string]string)
	contractsRoot := filepath.Join(root, "tests", "contracts")
	err := filepath.WalkDir(contractsRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != "contract.json" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var document contractDocument
		if err := json.Unmarshal(data, &document); err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}
		if document.Name == "" || document.Module == "" || document.Risk == "" || document.Status == "" {
			return fmt.Errorf("contract %s is missing identity metadata", path)
		}
		if document.Status != "promoted" {
			return fmt.Errorf("contract %s has non-promoted status %q", path, document.Status)
		}
		if previous, exists := seenContracts[document.Name]; exists {
			return fmt.Errorf("duplicate contract name %s in %s and %s", document.Name, previous, path)
		}
		domain, exists := byName[document.Module]
		if !exists {
			return fmt.Errorf("contract %s references unknown domain %s", document.Name, document.Module)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		domain.Contracts = append(domain.Contracts, contractRef{
			Name:       document.Name,
			Risk:       document.Risk,
			Status:     document.Status,
			Method:     document.Request.Method,
			URL:        document.Request.URL,
			SourcePath: filepath.ToSlash(relative),
		})
		seenContracts[document.Name] = path
		riskCounts[document.Risk]++
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("read contracts: %w", err)
	}

	for i := range domains {
		sort.Slice(domains[i].Contracts, func(a, b int) bool {
			return domains[i].Contracts[a].Name < domains[i].Contracts[b].Name
		})
	}
	return riskCounts, len(seenContracts), nil
}

func gitOutput(root string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func writeJSON(path string, value any) error {
	data, err := renderJSON(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func checkJSON(path string, value any) error {
	want, err := renderJSON(value)
	if err != nil {
		return err
	}
	got, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Equal(got, want) {
		return fmt.Errorf("%s is stale", path)
	}
	return nil
}

func renderJSON(value any) ([]byte, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func compareContractTrees(sourceRoot, destinationRoot string) error {
	source, err := contractTreeFiles(sourceRoot)
	if err != nil {
		return err
	}
	destination, err := contractTreeFiles(destinationRoot)
	if err != nil {
		return err
	}
	if len(source) != len(destination) {
		return fmt.Errorf("contract snapshot has %d files, source has %d", len(destination), len(source))
	}
	for relative, sourceData := range source {
		destinationData, ok := destination[relative]
		if !ok {
			return fmt.Errorf("contract snapshot is missing %s", relative)
		}
		if !bytes.Equal(sourceData, destinationData) {
			return fmt.Errorf("contract snapshot differs at %s", relative)
		}
	}
	return nil
}

func syncContractTrees(sourceRoot, destinationRoot string) error {
	sourceRoot, err := filepath.Abs(sourceRoot)
	if err != nil {
		return err
	}
	destinationRoot, err = filepath.Abs(destinationRoot)
	if err != nil {
		return err
	}
	if sourceRoot == destinationRoot || destinationRoot == filepath.VolumeName(destinationRoot)+string(filepath.Separator) {
		return errors.New("unsafe contract snapshot destination")
	}
	source, err := contractTreeFiles(sourceRoot)
	if err != nil {
		return err
	}
	for relative, data := range source {
		target := filepath.Join(destinationRoot, filepath.FromSlash(relative))
		resolved, err := filepath.Abs(target)
		if err != nil || !strings.HasPrefix(resolved, destinationRoot+string(filepath.Separator)) {
			return fmt.Errorf("contract path escapes destination: %s", relative)
		}
		if err := os.MkdirAll(filepath.Dir(resolved), 0o755); err != nil {
			return err
		}
		if current, err := os.ReadFile(resolved); err == nil && bytes.Equal(current, data) {
			continue
		}
		if err := os.WriteFile(resolved, data, 0o644); err != nil {
			return err
		}
	}
	destination, err := contractTreeFiles(destinationRoot)
	if err != nil {
		return err
	}
	for relative := range destination {
		if _, exists := source[relative]; !exists {
			return fmt.Errorf("stale contract snapshot file %s was not removed; review and delete it explicitly", relative)
		}
	}
	return compareContractTrees(sourceRoot, destinationRoot)
}

func contractTreeFiles(root string) (map[string][]byte, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]byte)
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("contract tree contains a symlink: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[filepath.ToSlash(relative)] = data
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read contract tree %s: %w", root, err)
	}
	return result, nil
}
