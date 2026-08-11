// Command bpi-probe audits the committed contract snapshot, generates the Go
// parity catalog, and runs explicitly gated read-only live probes.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yuelioi/bpi-go/internal/probe"
)

const usage = `usage: bpi-probe <command> [flags]

Offline commands:
  audit           verify source/mapping parity, fixtures, snapshot hashes, and privacy
  lock            regenerate the deterministic contract snapshot lock
  sanitize-audit  scan committed response fixtures for high-confidence secrets
  api-doc         generate or check docs/api-index.md

Network command:
  batch-run       run read-only contracts (requires BPI_PROBE=1 and --read-only)

Authenticated profiles read raw Cookie headers from BPI_COOKIE_NORMAL and
BPI_COOKIE_VIP. The SDK itself never reads credentials from the environment.

Run "bpi-probe <command> -h" for command flags.`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "bpi-probe:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		fmt.Println(usage)
		return nil
	}
	switch args[0] {
	case "audit":
		return runAudit(args[1:])
	case "lock":
		return runLock(args[1:])
	case "sanitize-audit":
		return runSanitizeAudit(args[1:])
	case "api-doc":
		return runAPIDoc(args[1:])
	case "batch-run":
		return runBatch(args[1:])
	default:
		return fmt.Errorf("unknown command %q\n%s", args[0], usage)
	}
}

func runAudit(args []string) error {
	set := flag.NewFlagSet("audit", flag.ContinueOnError)
	source := set.String("source", "parity/source.json", "locked source manifest")
	implemented := set.String("implemented", "parity/implemented.json", "Go implementation manifest")
	contracts := set.String("contracts", "testdata/contracts", "contract snapshot root")
	lock := set.String("lock", "parity/contracts.lock.json", "contract snapshot lock")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return errors.New("audit accepts flags only")
	}
	report, err := probe.Audit(*source, *implemented, *contracts, *lock)
	if err != nil {
		return err
	}
	return writeStdoutJSON(report)
}

func runLock(args []string) error {
	set := flag.NewFlagSet("lock", flag.ContinueOnError)
	sourcePath := set.String("source", "parity/source.json", "locked source manifest")
	contracts := set.String("contracts", "testdata/contracts", "contract snapshot root")
	output := set.String("output", "parity/contracts.lock.json", "lock output path")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return errors.New("lock accepts flags only")
	}
	var source probe.SourceManifest
	if err := probe.ReadJSON(*sourcePath, &source); err != nil {
		return err
	}
	if source.Source.Commit != probe.SourceCommit {
		return fmt.Errorf("source manifest commit is %s, want %s", source.Source.Commit, probe.SourceCommit)
	}
	lock, err := probe.BuildContractLock(*contracts, source.Source.Commit)
	if err != nil {
		return err
	}
	if err := probe.WriteJSON(*output, lock); err != nil {
		return err
	}
	fmt.Printf("wrote %s with %d files\n", *output, len(lock.Files))
	return nil
}

func runSanitizeAudit(args []string) error {
	set := flag.NewFlagSet("sanitize-audit", flag.ContinueOnError)
	contracts := set.String("contracts", "testdata/contracts", "contract snapshot root")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return errors.New("sanitize-audit accepts flags only")
	}
	findings, err := probe.AuditSanitization(*contracts)
	if err != nil {
		return err
	}
	if len(findings) != 0 {
		if err := writeStdoutJSON(findings); err != nil {
			return err
		}
		return fmt.Errorf("sanitize audit found %d issue(s)", len(findings))
	}
	fmt.Println("sanitize-audit ok: no high-confidence sensitive values found")
	return nil
}

func runAPIDoc(args []string) error {
	set := flag.NewFlagSet("api-doc", flag.ContinueOnError)
	source := set.String("source", "parity/source.json", "locked source manifest")
	implemented := set.String("implemented", "parity/implemented.json", "Go implementation manifest")
	output := set.String("output", "docs/api-index.md", "Markdown output path")
	check := set.Bool("check", false, "verify output is current without writing")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return errors.New("api-doc accepts flags only")
	}
	catalog, err := probe.GenerateCatalog(*source, *implemented)
	if err != nil {
		return err
	}
	if *check {
		current, err := os.ReadFile(*output)
		if err != nil {
			return fmt.Errorf("read catalog: %w", err)
		}
		if !bytes.Equal(current, []byte(catalog)) {
			return fmt.Errorf("%s is stale; run bpi-probe api-doc", *output)
		}
		fmt.Printf("api-doc check ok: %s\n", *output)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*output, []byte(catalog), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", *output)
	return nil
}

func runBatch(args []string) error {
	set := flag.NewFlagSet("batch-run", flag.ContinueOnError)
	contracts := set.String("contracts", "testdata/contracts", "contract snapshot root")
	profileList := set.String("profiles", "anonymous", "comma-separated anonymous,normal,vip profiles")
	readOnly := set.Bool("read-only", false, "required read-only risk gate")
	output := set.String("output", "", "optional summary JSON path (response bodies are never written)")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return errors.New("batch-run accepts flags only")
	}
	profiles, err := probe.ParseProfiles(*profileList)
	if err != nil {
		return err
	}
	summary, runErr := probe.RunBatch(context.Background(), probe.BatchConfig{
		ContractRoot: *contracts,
		CookieHeaders: map[string]string{
			"normal": os.Getenv("BPI_COOKIE_NORMAL"),
			"vip":    os.Getenv("BPI_COOKIE_VIP"),
		},
		Profiles: profiles,
		ReadOnly: *readOnly,
	})
	if *output != "" {
		if err := probe.WriteJSON(*output, summary); err != nil {
			return err
		}
		fmt.Printf("wrote safe probe summary %s\n", *output)
	} else if summary.SchemaVersion != 0 {
		if err := writeStdoutJSON(summary); err != nil {
			return err
		}
	}
	return runErr
}

func writeStdoutJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}
