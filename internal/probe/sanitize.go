package probe

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type SanitizeFinding struct {
	File   string `json:"file"`
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

var (
	cookieAssignmentPattern = regexp.MustCompile(`(?i)(SESSDATA|bili_jct|buvid3|buvid4|DedeUserID)\s*=`)
	emailPattern            = regexp.MustCompile(`(?i)\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b`)
)

var sensitiveFixtureKeys = map[string]struct{}{
	"sessdata":      {},
	"bilijct":       {},
	"csrf":          {},
	"csrftoken":     {},
	"buvid3":        {},
	"buvid4":        {},
	"accesskey":     {},
	"accesstoken":   {},
	"refreshtoken":  {},
	"cookie":        {},
	"setcookie":     {},
	"authorization": {},
}

func AuditSanitization(root string) ([]SanitizeFinding, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("probe: resolve sanitize root: %w", err)
	}
	var findings []SanitizeFinding
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".json") || !hasPathComponent(path, "responses") {
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
		file := filepath.ToSlash(relative)
		if cookieAssignmentPattern.Match(data) {
			findings = append(findings, SanitizeFinding{File: file, Path: "$", Reason: "contains a credential-like Cookie assignment"})
		}
		var value any
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}
		auditSanitizedValue(file, "$", value, &findings)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("probe: sanitize audit: %w", err)
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].File == findings[j].File {
			return findings[i].Path < findings[j].Path
		}
		return findings[i].File < findings[j].File
	})
	return findings, nil
}

func auditSanitizedValue(file, path string, value any, findings *[]SanitizeFinding) {
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			childPath := path + "." + key
			if _, sensitive := sensitiveFixtureKeys[normalizeKey(key)]; sensitive && !isRedactedValue(child) {
				*findings = append(*findings, SanitizeFinding{File: file, Path: childPath, Reason: "contains a non-redacted credential field"})
			}
			auditSanitizedValue(file, childPath, child, findings)
		}
	case []any:
		for index, child := range current {
			auditSanitizedValue(file, fmt.Sprintf("%s[%d]", path, index), child, findings)
		}
	case string:
		if !strings.Contains(current, "://") && emailPattern.MatchString(current) && !strings.Contains(strings.ToLower(current), "redacted") {
			*findings = append(*findings, SanitizeFinding{File: file, Path: path, Reason: "contains an email-like value"})
		}
	}
}

func normalizeKey(value string) string {
	return strings.Map(func(character rune) rune {
		if character >= 'A' && character <= 'Z' {
			return character + ('a' - 'A')
		}
		if character >= 'a' && character <= 'z' || character >= '0' && character <= '9' {
			return character
		}
		return -1
	}, value)
}

func isRedactedValue(value any) bool {
	if value == nil {
		return true
	}
	text, ok := value.(string)
	if !ok {
		return false
	}
	text = strings.ToLower(strings.TrimSpace(text))
	return text == "" || text == "<redacted>" || text == "${csrf}" || strings.HasPrefix(text, "fixture-") || strings.Contains(text, "sanitized") || strings.Contains(text, "redacted")
}

func hasPathComponent(path, want string) bool {
	for _, component := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
		if component == want {
			return true
		}
	}
	return false
}
